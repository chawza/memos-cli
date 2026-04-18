package api

type Memo struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	UID        string `json:"uid"`
	Creator    string `json:"creator"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
	DisplayTime string `json:"displayTime"`
	Content    string `json:"content"`
	Visibility string `json:"visibility"`
	Pinned     bool   `json:"pinned"`
	Snippet    string `json:"snippet"`
}

type CreateMemoRequest struct {
	Memo   *CreateMemo `json:"memo"`
	MemoID string      `json:"memoId,omitempty"`
}

type CreateMemo struct {
	Content    string `json:"content"`
	Visibility string `json:"visibility,omitempty"`
	State      string `json:"state,omitempty"`
	Pinned     bool   `json:"pinned,omitempty"`
}

type UpdateMemoRequest struct {
	Memo       *UpdateMemo `json:"memo"`
	UpdateMask string      `json:"updateMask"`
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
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}
