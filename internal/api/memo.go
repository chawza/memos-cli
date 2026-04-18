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

func (c *Client) CreateMemo(memo *CreateMemo) (*Memo, error) {
	var result Memo
	err := c.do(http.MethodPost, "/api/v1/memos", memo, &result)
	return &result, err
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

func (c *Client) UpdateMemo(id string, memo *UpdateMemo, updateMask string) (*Memo, error) {
	path := "/api/v1/" + memoName(id)
	if updateMask != "" {
		path += "?updateMask=" + url.QueryEscape(updateMask)
	}
	var result Memo
	err := c.do(http.MethodPatch, path, memo, &result)
	return &result, err
}

func (c *Client) DeleteMemo(id string) error {
	return c.do(http.MethodDelete, "/api/v1/"+memoName(id), nil, nil)
}

func (c *Client) ListMemoComments(memoID string) ([]Memo, string, error) {
	path := "/api/v1/" + memoName(memoID) + "/comments"
	var resp ListMemosResponse
	err := c.do(http.MethodGet, path, nil, &resp)
	return resp.Memos, resp.NextPageToken, err
}

func (c *Client) CreateMemoComment(memoID string, comment *CreateMemo) (*Memo, error) {
	path := "/api/v1/" + memoName(memoID) + "/comments"
	var result Memo
	err := c.do(http.MethodPost, path, comment, &result)
	return &result, err
}

func (c *Client) ListMemoReactions(memoID string) ([]Reaction, string, error) {
	path := "/api/v1/" + memoName(memoID) + "/reactions"
	var resp ListReactionsResponse
	err := c.do(http.MethodGet, path, nil, &resp)
	return resp.Reactions, resp.NextPageToken, err
}

func (c *Client) UpsertMemoReaction(memoID string, reaction *UpsertReaction) (*Reaction, error) {
	path := "/api/v1/" + memoName(memoID) + "/reactions"
	req := UpsertReactionRequest{
		Name:     memoName(memoID),
		Reaction: *reaction,
	}
	var result Reaction
	err := c.do(http.MethodPost, path, req, &result)
	return &result, err
}

func (c *Client) DeleteMemoReaction(memoID, reactionID string) error {
	path := "/api/v1/" + memoName(memoID) + "/reactions/" + reactionID
	return c.do(http.MethodDelete, path, nil, nil)
}

func (c *Client) ListMemoAttachments(memoID string) ([]Attachment, string, error) {
	path := "/api/v1/" + memoName(memoID) + "/attachments"
	var resp ListAttachmentsResponse
	err := c.do(http.MethodGet, path, nil, &resp)
	return resp.Attachments, resp.NextPageToken, err
}

func (c *Client) SetMemoAttachments(memoID string, attachments []Attachment) error {
	path := "/api/v1/" + memoName(memoID) + "/attachments"
	req := SetAttachmentsRequest{
		Name:        memoName(memoID),
		Attachments: attachments,
	}
	return c.do(http.MethodPatch, path, req, nil)
}
