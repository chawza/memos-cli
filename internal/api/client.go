package api

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	BaseURL    string
	Token      string
	httpClient *http.Client
}

type ClientOption func(*Client)

func WithTimeout(seconds int) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = time.Duration(seconds) * time.Second
	}
}

func WithTLSSkipVerify(skip bool) ClientOption {
	return func(c *Client) {
		if skip {
			transport, ok := c.httpClient.Transport.(*http.Transport)
			if !ok {
				transport = &http.Transport{}
			}
			transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			c.httpClient.Transport = transport
		}
	}
}

func NewClient(baseURL, token string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Ping() error {
	_, _, err := c.ListMemos(1, "", "", "")
	if err != nil {
		return fmt.Errorf("connectivity check failed: %w", err)
	}
	return nil
}

func (c *Client) do(method, path string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		bodyReader = strings.NewReader(string(b))
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("API error %d: %s", resp.StatusCode, strings.TrimSpace(string(bodyBytes))),
		}
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
