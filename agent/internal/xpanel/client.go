package xpanel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"xpanel-central/agent/internal/retry"
)

type Client struct {
	baseURL       string
	username      string
	password      string
	http          *http.Client
	retryPolicy   retry.Policy
	mu            sync.Mutex
	loginFailures int
	nextLoginAt   time.Time
}

type Response struct {
	Success bool            `json:"success"`
	Msg     string          `json:"msg"`
	Obj     json.RawMessage `json:"obj"`
	Data    json.RawMessage `json:"data"`
}

func NewClient(baseURL, basePath, username, password string, timeout time.Duration) (*Client, error) {
	return NewClientWithRetryPolicy(baseURL, basePath, username, password, timeout, retry.DefaultPolicy)
}

// NewClientWithRetryPolicy is intended for staging and tests that need a
// shorter retry window. NewClient retains the production 30s/2m/10m policy.
func NewClientWithRetryPolicy(baseURL, basePath, username, password string, timeout time.Duration, policy retry.Policy) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, errors.New("xpanel URL is required")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, errors.New("xpanel URL must include scheme and host")
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("create cookie jar: %w", err)
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/") + "/" + strings.Trim(strings.TrimSpace(basePath), "/"),
		username: username,
		password: password,
		// X-Panel commonly redirects unauthenticated API requests to its HTML
		// login page. Do not follow that redirect here: doing so hides the
		// authentication challenge and makes the JSON decoder fail on '<'.
		// The GET path below handles the redirect by logging in and retrying.
		http: &http.Client{
			Timeout:       timeout,
			Jar:           jar,
			CheckRedirect: noRedirect,
		},
		retryPolicy: policy.Normalized(),
	}, nil
}

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

func (client *Client) Login(ctx context.Context) error {
	client.mu.Lock()
	defer client.mu.Unlock()
	return client.loginLocked(ctx)
}

func (client *Client) loginLocked(ctx context.Context) error {
	if !client.nextLoginAt.IsZero() && time.Now().Before(client.nextLoginAt) {
		return fmt.Errorf("xpanel login backoff active until %s", client.nextLoginAt.UTC().Format(time.RFC3339))
	}
	policy := client.retryPolicy.Normalized()
	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		var transient bool
		lastErr, transient = client.loginAttempt(ctx)
		if lastErr == nil {
			return nil
		}
		if !transient || attempt == policy.MaxAttempts {
			client.recordLoginFailure()
			return lastErr
		}
		if err := retry.Wait(ctx, policy.Delay(attempt)); err != nil {
			return fmt.Errorf("xpanel login retry canceled: %w", err)
		}
	}
	return lastErr
}

func (client *Client) loginAttempt(ctx context.Context) (error, bool) {
	form := url.Values{"username": {client.username}, "password": {client.password}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint("/login"), strings.NewReader(form.Encode()))
	if err != nil {
		return err, false
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.http.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("xpanel login request: %w", ctx.Err()), false
		}
		return fmt.Errorf("xpanel login request: %w", err), true
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("xpanel login returned HTTP %d", response.StatusCode), retryableHTTPStatus(response.StatusCode)
	}
	var payload Response
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		return fmt.Errorf("decode xpanel login response: %w", err), true
	}
	if !payload.Success {
		return fmt.Errorf("xpanel login failed: %s", safeMessage(payload.Msg)), false
	}
	client.loginFailures = 0
	client.nextLoginAt = time.Time{}
	return nil, false
}

func (client *Client) Get(ctx context.Context, path string) (Response, error) {
	return client.doGet(ctx, path, true)
}

func (client *Client) doGet(ctx context.Context, path string, retryAuth bool) (Response, error) {
	policy := client.retryPolicy.Normalized()
	authRetry := retryAuth
	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		response, status, err, transient := client.getAttempt(ctx, path)
		if err != nil {
			lastErr = err
		} else if isAuthChallenge(status) && authRetry {
			client.mu.Lock()
			loginErr := client.loginLocked(ctx)
			client.mu.Unlock()
			if loginErr != nil {
				return Response{}, loginErr
			}
			authRetry = false
			// A 401 only consumes the authentication retry, not the transient
			// request budget. This keeps a session refresh from reducing the
			// number of network retries available to the actual GET.
			attempt--
			continue
		} else if err != nil {
			lastErr = err
		} else if status < 200 || status >= 300 {
			lastErr = fmt.Errorf("xpanel GET %s returned HTTP %d", path, status)
		} else if response.Success {
			return response, nil
		} else {
			// X-Panel application failures are actionable configuration or
			// permission errors; replaying them would only amplify load.
			return Response{}, fmt.Errorf("xpanel GET %s failed: %s", path, safeMessage(response.Msg))
		}
		if !transient || attempt == policy.MaxAttempts {
			break
		}
		if err := retry.Wait(ctx, policy.Delay(attempt)); err != nil {
			return Response{}, fmt.Errorf("xpanel GET %s retry canceled: %w", path, err)
		}
	}
	return Response{}, lastErr
}

func (client *Client) getAttempt(ctx context.Context, path string) (Response, int, error, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.endpoint("/panel/api/"+strings.TrimLeft(path, "/")), nil)
	if err != nil {
		return Response{}, 0, err, false
	}
	response, err := client.http.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return Response{}, 0, fmt.Errorf("xpanel GET %s: %w", path, ctx.Err()), false
		}
		return Response{}, 0, fmt.Errorf("xpanel GET %s: %w", path, err), true
	}
	defer response.Body.Close()
	status := response.StatusCode
	if status == http.StatusUnauthorized {
		return Response{}, status, nil, false
	}
	if status < 200 || status >= 300 {
		return Response{}, status, nil, retryableHTTPStatus(status)
	}
	var payload Response
	if err := json.NewDecoder(io.LimitReader(response.Body, 8<<20)).Decode(&payload); err != nil {
		return Response{}, status, fmt.Errorf("decode xpanel response: %w", err), true
	}
	return payload, status, nil, false
}

func retryableHTTPStatus(status int) bool {
	return status == http.StatusRequestTimeout || status == http.StatusTooEarly || status == http.StatusTooManyRequests || status >= 500
}

func isAuthChallenge(status int) bool {
	return status == http.StatusUnauthorized || (status >= 300 && status < 400)
}

func (client *Client) recordLoginFailure() {
	client.loginFailures++
	backoff := time.Second << min(client.loginFailures-1, 8)
	if backoff > 5*time.Minute {
		backoff = 5 * time.Minute
	}
	client.nextLoginAt = time.Now().Add(backoff)
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func (client *Client) endpoint(path string) string {
	return strings.TrimRight(client.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func safeMessage(message string) string {
	message = strings.TrimSpace(message)
	if len(message) > 160 {
		return message[:160]
	}
	return message
}
