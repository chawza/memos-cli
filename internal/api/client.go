package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client wraps the Memos API v1.
type Client struct {
	BaseURL string
	Token   string
	httpClient *http.Client
}

// NewClient returns a Memos API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		Token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// do sends an HTTP request with auth headers and decodes the response.
func (c *Client) do(method, path string, body interface{}, resp interface{}) error {
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

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(res.Body)
		return fmt.Errorf("API error %d: %s", res.StatusCode, string(bodyBytes))
	}

	if resp != nil {
		if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}

// CreateMemoRequest is the request body for creating a memo.
type CreateMemoRequest struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility,omitempty"`
	State      string `json:"state,omitempty"`
	Pinned     bool   `json:"pinned,omitempty"`
}

// Memo represents a Memos memo object.
type Memo struct {
	ID         string    `json:"id"`
	Name       string    `json:"name,omitempty"`
	Content    string    `json:"content"`
	Visibility string    `json:"visibility"`
	State      string    `json:"state"`
	Pinned     bool      `json:"pinned"`
	CreatorID  string    `json:"creatorId,omitempty"`
	CreateTime time.Time `json:"createTime,omitempty"`
	UpdateTime time.Time `json:"updateTime,omitempty"`
}

// CreateMemo creates a new memo.
func (c *Client) CreateMemo(req *CreateMemoRequest) (*Memo, error) {
	var memo Memo
	err := c.do(http.MethodPost, "/api/v1/memos", req, &memo)
	return &memo, err
}

// ListMemos lists memos. Pass filter as an AIP-160 filter string, e.g. "visibility == \"PRIVATE\"".
func (c *Client) ListMemos(pageSize int, pageToken, filter string) ([]Memo, string, error) {
	q := url.Values{}
	if pageSize > 0 {
		q.Set("pageSize", fmt.Sprintf("%d", pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	if filter != "" {
		q.Set("filter", filter)
	}

	path := "/api/v1/memos"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var struct {
		Memos     []Memo `json:"memos"`
		NextPageToken string `json:"nextPageToken"`
	}
	err := c.do(http.MethodGet, path, nil, &struct{
		Memos     []Memo `json:"memos"`
		NextPageToken string `json:"nextPageToken"`
	}{
		&Memos,
		"",
	})
	return struct{
		Memos     []Memo `json:"memos"`
		NextPageToken string `json:"nextPageToken"`
	}.Memos, struct{
		Memos     []Memo `json:"memos"`
		NextPageToken string `json:"nextPageToken"`
	}.NextPageToken, err
}

// GetMemo fetches a single memo by ID.
func (c *Client) GetMemo(id string) (*Memo, error) {
	var memo Memo
	err := c.do(http.MethodGet, fmt.Sprintf("/api/v1/memos/%s", id), nil, &memo)
	return &memo, err
}

// UpdateMemoRequest is the request body for updating a memo.
type UpdateMemoRequest struct {
	Content    *string `json:"content,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	State      *string `json:"state,omitempty"`
	Pinned     *bool   `json:"pinned,omitempty"`
}

// UpdateMemo updates an existing memo. updateMask is a comma-separated list of fields to update.
func (c *Client) UpdateMemo(id string, req *UpdateMemoRequest, updateMask string) (*Memo, error) {
	path := fmt.Sprintf("/api/v1/memos/%s?updateMask=%s", id, updateMask)
	var memo Memo
	err := c.do(http.MethodPatch, path, req, &memo)
	return &memo, err
}

// DeleteMemo deletes a memo by ID.
func (c *Client) DeleteMemo(id string) error {
	return c.do(http.MethodDelete, fmt.Sprintf("/api/v1/memos/%s", id), nil, nil)
}
