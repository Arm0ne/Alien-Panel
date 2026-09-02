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
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string, timeout time.Duration) (*Client, error) {
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
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), token: token, http: &http.Client{Timeout: timeout}}, nil
}

func (client *Client) PostJSON(ctx context.Context, path string, payload any, result any) error {
	encoded, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode central request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+"/"+strings.TrimLeft(path, "/"), bytes.NewReader(encoded))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	if client.token != "" {
		request.Header.Set("Authorization", "Bearer "+client.token)
	}
	request.Header.Set("X-Request-ID", requestID())
	request.Header.Set("X-Request-Timestamp", strconv.FormatInt(time.Now().Unix(), 10))
	response, err := client.http.Do(request)
	if err != nil {
		return fmt.Errorf("central POST %s: %w", path, err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("central POST %s returned HTTP %d", path, response.StatusCode)
	}
	var envelope struct {
		Code string          `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 8<<20)).Decode(&envelope); err != nil {
		return fmt.Errorf("decode central response: %w", err)
	}
	if envelope.Code != "0000" {
		return fmt.Errorf("central POST %s failed: %s", path, envelope.Msg)
	}
	if result != nil && len(envelope.Data) > 0 && string(envelope.Data) != "null" {
		if err := json.Unmarshal(envelope.Data, result); err != nil {
			return fmt.Errorf("decode central payload: %w", err)
		}
	}
	return nil
}

func requestID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return fmt.Sprintf("agent-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
