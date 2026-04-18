package api

import "fmt"

type Memo struct {
	Name        string         `json:"name"`
	State       string         `json:"state"`
	UID         string         `json:"uid"`
	Creator     string         `json:"creator"`
	CreateTime  string         `json:"createTime"`
	UpdateTime  string         `json:"updateTime"`
	DisplayTime string         `json:"displayTime"`
	Content     string         `json:"content"`
	Visibility  string         `json:"visibility"`
	Pinned      bool           `json:"pinned"`
	Snippet     string         `json:"snippet"`
	Tags        []string       `json:"tags"`
	Property    *MemoProperty  `json:"property"`
	Location    *MemoLocation  `json:"location"`
}

type MemoProperty struct {
	HasLink          bool   `json:"hasLink"`
	HasTaskList      bool   `json:"hasTaskList"`
	HasCode          bool   `json:"hasCode"`
	HasIncompleteTasks bool `json:"hasIncompleteTasks"`
	Title            string `json:"title"`
}

type MemoLocation struct {
	Placeholder string  `json:"placeholder"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type CreateMemo struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility,omitempty"`
	State      string `json:"state,omitempty"`
	Pinned     bool   `json:"pinned,omitempty"`
}

type UpdateMemo struct {
	Name       string  `json:"name"`
	Content    *string `json:"content,omitempty"`
	Visibility *string `json:"visibility,omitempty"`
	State      *string `json:"state,omitempty"`
	Pinned     *bool   `json:"pinned,omitempty"`
}

type ListMemosResponse struct {
	Memos         []Memo `json:"memos"`
	NextPageToken string `json:"nextPageToken"`
}

type APIError struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Details    []any  `json:"details"`
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("API error %d (code %d): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("API error %d", e.StatusCode)
}

type Reaction struct {
	Name         string `json:"name"`
	Creator      string `json:"creator"`
	ContentId    string `json:"contentId"`
	ReactionType string `json:"reactionType"`
	CreateTime   string `json:"createTime"`
}

type Attachment struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
}

type UpsertReactionRequest struct {
	Name     string   `json:"name"`
	Reaction Reaction `json:"reaction"`
}

type UpsertReaction struct {
	ReactionType string `json:"reactionType"`
}

type ListReactionsResponse struct {
	Reactions     []Reaction `json:"reactions"`
	NextPageToken string     `json:"nextPageToken"`
	TotalSize     int        `json:"totalSize"`
}

type ListAttachmentsResponse struct {
	Attachments   []Attachment `json:"attachments"`
	NextPageToken string       `json:"nextPageToken"`
}

type SetAttachmentsRequest struct {
	Name        string       `json:"name"`
	Attachments []Attachment `json:"attachments"`
}
