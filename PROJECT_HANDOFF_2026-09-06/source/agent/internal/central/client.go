package central

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"xpanel-central/agent/internal/retry"
)

type Client struct {
	baseURL     string
	token       string
	http        *http.Client
	retryPolicy retry.Policy
}

func (client *Client) Heartbeat(ctx context.Context, payload any) error {
	return client.PostJSON(ctx, "/agent/v1/heartbeat", payload, nil)
}

func (client *Client) Sync(ctx context.Context, payload any, result any) error {
	return client.PostJSON(ctx, "/agent/v1/sync", payload, result)
}

func NewClient(baseURL, token string, timeout time.Duration) (*Client, error) {
	return NewClientWithRetryPolicy(baseURL, token, timeout, retry.DefaultPolicy)
}

// NewClientWithRetryPolicy is useful for deployments that need a different
// failure budget (for example, a staging environment with short test delays).
// Production callers should normally use NewClient so the documented
// 30s/2m/10m backoff is retained.
func NewClientWithRetryPolicy(baseURL, token string, timeout time.Duration, policy retry.Policy) (*Client, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("central URL is required")
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("central URL must include scheme and host")
	}
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), token: token, http: &http.Client{Timeout: timeout}, retryPolicy: policy.Normalized()}, nil
}

func (client *Client) PostJSON(ctx context.Context, path string, payload any, result any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode central request: %w", err)
	}
	logicalRequestID := requestID()
	policy := client.retryPolicy.Normalized()
	var lastErr error
	for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
		var retryable bool
		lastErr, retryable = client.postJSONAttempt(ctx, path, encoded, result, logicalRequestID)
		if lastErr == nil {
			return nil
		}
		if !retryable || attempt == policy.MaxAttempts {
			break
		}
		if err := retry.Wait(ctx, policy.Delay(attempt)); err != nil {
			return fmt.Errorf("central POST %s retry canceled: %w", path, err)
		}
	}
	return lastErr
}

func (client *Client) postJSONAttempt(ctx context.Context, path string, encoded []byte, result any, logicalRequestID string) (error, bool) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+"/"+strings.TrimLeft(path, "/"), bytes.NewReader(encoded))
	if err != nil {
		return err, false
	}
	request.Header.Set("Content-Type", "application/json")
	if client.token != "" {
		request.Header.Set("Authorization", "Bearer "+client.token)
	}
	request.Header.Set("X-Request-ID", logicalRequestID)
	request.Header.Set("X-Request-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	response, err := client.http.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("central POST %s: %w", path, ctx.Err()), false
		}
		return fmt.Errorf("central POST %s: %w", path, err), true
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		retryable := retryableHTTPStatus(response.StatusCode)
		return fmt.Errorf("central POST %s returned HTTP %d", path, response.StatusCode), retryable
	}
	var envelope struct {
		Code string          `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 8<<20)).Decode(&envelope); err != nil {
		return fmt.Errorf("decode central response: %w", err), true
	}
	if envelope.Code != "0000" {
		// Application errors are generally validation/auth failures and must not
		// be replayed. The sync endpoint itself is idempotent by sync_id, but
		// retrying arbitrary business errors would hide actionable configuration
		// problems.
		return fmt.Errorf("central POST %s failed: %s", path, envelope.Msg), false
	}
	if result != nil && len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, result); err != nil {
			return fmt.Errorf("decode central payload: %w", err), false
		}
	}
	return nil, false
}

func retryableHTTPStatus(status int) bool {
	return status == http.StatusRequestTimeout || status == http.StatusTooEarly || status == http.StatusTooManyRequests || status >= 500
}

func requestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("agent-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
