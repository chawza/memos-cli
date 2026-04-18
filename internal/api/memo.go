package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func memoName(id string) string {
	if strings.HasPrefix(id, "memos/") {
		return id
	}
	return "memos/" + id
}

func (c *Client) CreateMemo(req *CreateMemoRequest) (*Memo, error) {
	var memo Memo
	err := c.do(http.MethodPost, "/api/v1/memos", req, &memo)
	return &memo, err
}

func (c *Client) ListMemos(pageSize int, pageToken, filter, state string) ([]Memo, string, error) {
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
	if state != "" {
		q.Set("state", state)
	}

	path := "/api/v1/memos"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	var resp ListMemosResponse
	err := c.do(http.MethodGet, path, nil, &resp)
	return resp.Memos, resp.NextPageToken, err
}

func (c *Client) GetMemo(id string) (*Memo, error) {
	var memo Memo
	err := c.do(http.MethodGet, "/api/v1/"+memoName(id), nil, &memo)
	return &memo, err
}

func (c *Client) UpdateMemo(id string, req *UpdateMemoRequest) (*Memo, error) {
	var memo Memo
	err := c.do(http.MethodPatch, "/api/v1/"+memoName(id), req, &memo)
	return &memo, err
}

func (c *Client) DeleteMemo(id string) error {
	return c.do(http.MethodDelete, "/api/v1/"+memoName(id), nil, nil)
}
