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
)

type Client struct {
	baseURL       string
	username      string
	password      string
	http          *http.Client
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
		http:     &http.Client{Timeout: timeout, Jar: jar},
	}, nil
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
	form := url.Values{"username": {client.username}, "password": {client.password}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint("/login"), strings.NewReader(form.Encode()))
	if err != nil {
		client.recordLoginFailure()
		return err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.http.Do(request)
	if err != nil {
		client.recordLoginFailure()
		return fmt.Errorf("xpanel login request: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		client.recordLoginFailure()
		return fmt.Errorf("xpanel login returned HTTP %d", response.StatusCode)
	}
	var payload Response
	if err := json.NewDecoder(io.LimitReader(response.Body, 1<<20)).Decode(&payload); err != nil {
		client.recordLoginFailure()
		return fmt.Errorf("decode xpanel login response: %w", err)
	}
	if !payload.Success {
		client.recordLoginFailure()
		return fmt.Errorf("xpanel login failed: %s", safeMessage(payload.Msg))
	}
	client.loginFailures = 0
	client.nextLoginAt = time.Time{}
	return nil
}

func (client *Client) Get(ctx context.Context, path string) (Response, error) {
	return client.doGet(ctx, path, true)
}

func (client *Client) doGet(ctx context.Context, path string, retry bool) (Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, client.endpoint("/panel/api/"+strings.TrimLeft(path, "/")), nil)
	if err != nil {
		return Response{}, err
	}
	response, err := client.http.Do(request)
	if err != nil {
		return Response{}, fmt.Errorf("xpanel GET %s: %w", path, err)
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized && retry {
		client.mu.Lock()
		loginErr := client.loginLocked(ctx)
		client.mu.Unlock()
		if loginErr != nil {
			return Response{}, loginErr
		}
		return client.doGet(ctx, path, false)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return Response{}, fmt.Errorf("xpanel GET %s returned HTTP %d", path, response.StatusCode)
	}
	var payload Response
	if err := json.NewDecoder(io.LimitReader(response.Body, 8<<20)).Decode(&payload); err != nil {
		return Response{}, fmt.Errorf("decode xpanel response: %w", err)
	}
	if !payload.Success {
		return Response{}, fmt.Errorf("xpanel GET %s failed: %s", path, safeMessage(payload.Msg))
	}
	return payload, nil
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
