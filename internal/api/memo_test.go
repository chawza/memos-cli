package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func testServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("fixtures", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to load fixture %s: %v", name, err)
	}
	return data
}

func TestCreateMemo(t *testing.T) {
	fixture := loadFixture(t, "create_memo_request.json")
	var expectedReq CreateMemo
	if err := json.Unmarshal(fixture, &expectedReq); err != nil {
		t.Fatalf("parse fixture: %v", err)
	}

	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/memos" {
			t.Errorf("expected /api/v1/memos, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		var req CreateMemo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Content != expectedReq.Content {
			t.Errorf("expected content %q, got %q", expectedReq.Content, req.Content)
		}
		if req.Visibility != expectedReq.Visibility {
			t.Errorf("expected visibility %q, got %q", expectedReq.Visibility, req.Visibility)
		}
		if req.State != expectedReq.State {
			t.Errorf("expected state %q, got %q", expectedReq.State, req.State)
		}
		if req.Pinned != expectedReq.Pinned {
			t.Errorf("expected pinned %v, got %v", expectedReq.Pinned, req.Pinned)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_full.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.CreateMemo(&CreateMemo{
		Content:    expectedReq.Content,
		Visibility: expectedReq.Visibility,
		State:      expectedReq.State,
		Pinned:     expectedReq.Pinned,
	})
	if err != nil {
		t.Fatalf("CreateMemo: %v", err)
	}

	if memo.Name != "memos/1" {
		t.Errorf("expected name memos/1, got %s", memo.Name)
	}
	if memo.State != "NORMAL" {
		t.Errorf("expected state NORMAL, got %s", memo.State)
	}
	if memo.UID != "abc123def456" {
		t.Errorf("expected uid abc123def456, got %s", memo.UID)
	}
	if memo.Creator != "users/1" {
		t.Errorf("expected creator users/1, got %s", memo.Creator)
	}
	if memo.CreateTime != "2024-01-15T10:30:00Z" {
		t.Errorf("expected createTime 2024-01-15T10:30:00Z, got %s", memo.CreateTime)
	}
	if memo.Content != "# Hello\n\nThis is a **test** memo with some content." {
		t.Errorf("unexpected content: %q", memo.Content)
	}
	if memo.Visibility != "PRIVATE" {
		t.Errorf("expected visibility PRIVATE, got %s", memo.Visibility)
	}
	if memo.Pinned != false {
		t.Errorf("expected pinned false, got %v", memo.Pinned)
	}
	if len(memo.Tags) != 3 || memo.Tags[0] != "test" || memo.Tags[1] != "demo" || memo.Tags[2] != "hello" {
		t.Errorf("expected tags [test demo hello], got %v", memo.Tags)
	}
	if memo.Property == nil {
		t.Fatal("expected property to be present")
	}
	if memo.Property.HasLink != false {
		t.Errorf("expected property.hasLink false, got %v", memo.Property.HasLink)
	}
	if memo.Property.HasTaskList != true {
		t.Errorf("expected property.hasTaskList true, got %v", memo.Property.HasTaskList)
	}
	if memo.Property.Title != "Hello" {
		t.Errorf("expected property.title Hello, got %q", memo.Property.Title)
	}
	if memo.Location == nil {
		t.Fatal("expected location to be present")
	}
	if memo.Location.Placeholder != "San Francisco, CA" {
		t.Errorf("expected location placeholder %q, got %q", "San Francisco, CA", memo.Location.Placeholder)
	}
	if memo.Location.Latitude != 37.7749 {
		t.Errorf("expected location latitude 37.7749, got %v", memo.Location.Latitude)
	}
	if memo.Location.Longitude != -122.4194 {
		t.Errorf("expected location longitude -122.4194, got %v", memo.Location.Longitude)
	}
}

func TestGetMemo(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/memos/42" {
			t.Errorf("expected /api/v1/memos/42, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_full.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.GetMemo("42")
	if err != nil {
		t.Fatalf("GetMemo: %v", err)
	}
	if memo.Name != "memos/1" {
		t.Errorf("expected name memos/1, got %s", memo.Name)
	}
	if memo.Content != "# Hello\n\nThis is a **test** memo with some content." {
		t.Errorf("unexpected content: %q", memo.Content)
	}
	if memo.Snippet != "Hello This is a test memo with some content." {
		t.Errorf("expected snippet %q, got %q", "Hello This is a test memo with some content.", memo.Snippet)
	}
}

func TestGetMemoWithPrefix(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/42" {
			t.Errorf("expected /api/v1/memos/42, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_full.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.GetMemo("memos/42")
	if err != nil {
		t.Fatalf("GetMemo: %v", err)
	}
	if memo.Name != "memos/1" {
		t.Errorf("expected name memos/1, got %s", memo.Name)
	}
}

func TestListMemos(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/memos" {
			t.Errorf("expected /api/v1/memos, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("pageSize") != "10" {
			t.Errorf("expected pageSize=10, got %s", r.URL.Query().Get("pageSize"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_list.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memos, next, err := c.ListMemos(10, "", "", "")
	if err != nil {
		t.Fatalf("ListMemos: %v", err)
	}
	if len(memos) != 2 {
		t.Fatalf("expected 2 memos, got %d", len(memos))
	}
	if memos[0].Name != "memos/1" {
		t.Errorf("expected first memo memos/1, got %s", memos[0].Name)
	}
	if memos[0].Visibility != "PUBLIC" {
		t.Errorf("expected first memo visibility PUBLIC, got %s", memos[0].Visibility)
	}
	if memos[0].Pinned != true {
		t.Errorf("expected first memo pinned true, got %v", memos[0].Pinned)
	}
	if len(memos[0].Tags) != 1 || memos[0].Tags[0] != "work" {
		t.Errorf("expected first memo tags [work], got %v", memos[0].Tags)
	}
	if memos[1].Name != "memos/2" {
		t.Errorf("expected second memo memos/2, got %s", memos[1].Name)
	}
	if memos[1].Property == nil || memos[1].Property.HasLink != true {
		t.Errorf("expected second memo property.hasLink true")
	}
	if next != "page2token" {
		t.Errorf("expected next page token 'page2token', got %q", next)
	}
}

func TestListMemosEmpty(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_list_empty.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memos, next, err := c.ListMemos(10, "", "", "")
	if err != nil {
		t.Fatalf("ListMemos: %v", err)
	}
	if len(memos) != 0 {
		t.Errorf("expected 0 memos, got %d", len(memos))
	}
	if next != "" {
		t.Errorf("expected empty nextPageToken, got %q", next)
	}
}

func TestListMemosWithFilter(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("filter") != "creator == \"users/1\"" {
			t.Errorf("expected filter=%q, got %q", "creator == \"users/1\"", r.URL.Query().Get("filter"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_list_empty.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	_, _, err := c.ListMemos(10, "", "creator == \"users/1\"", "")
	if err != nil {
		t.Fatalf("ListMemos: %v", err)
	}
}

func TestListMemosWithState(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("state") != "ARCHIVED" {
			t.Errorf("expected state=ARCHIVED, got %s", r.URL.Query().Get("state"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_list_empty.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	_, _, err := c.ListMemos(10, "", "", "ARCHIVED")
	if err != nil {
		t.Fatalf("ListMemos: %v", err)
	}
}

func TestListMemosPagination(t *testing.T) {
	callCount := 0
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Query().Get("pageToken") == "" {
			w.Write(loadFixture(t, "memo_list.json"))
		} else {
			if r.URL.Query().Get("pageToken") != "page2token" {
				t.Errorf("expected pageToken=page2token, got %s", r.URL.Query().Get("pageToken"))
			}
			w.Write(loadFixture(t, "memo_list_page2.json"))
		}
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")

	memos, next, err := c.ListMemos(10, "", "", "")
	if err != nil {
		t.Fatalf("page 1: %v", err)
	}
	if len(memos) != 2 {
		t.Fatalf("expected 2 memos on page 1, got %d", len(memos))
	}
	if next != "page2token" {
		t.Fatalf("expected next pageToken page2token, got %q", next)
	}

	memos2, next2, err := c.ListMemos(10, next, "", "")
	if err != nil {
		t.Fatalf("page 2: %v", err)
	}
	if len(memos2) != 1 {
		t.Fatalf("expected 1 memo on page 2, got %d", len(memos2))
	}
	if memos2[0].Name != "memos/3" {
		t.Errorf("expected memo/3 on page 2, got %s", memos2[0].Name)
	}
	if next2 != "" {
		t.Errorf("expected empty nextPageToken on last page, got %q", next2)
	}
	if callCount != 2 {
		t.Errorf("expected 2 server calls, got %d", callCount)
	}
}

func TestUpdateMemo(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/memos/5" {
			t.Errorf("expected /api/v1/memos/5, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("updateMask") != "content,visibility" {
			t.Errorf("expected updateMask=content,visibility, got %s", r.URL.Query().Get("updateMask"))
		}

		var req UpdateMemo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if *req.Content != "Updated content here" {
			t.Errorf("expected content %q, got %q", "Updated content here", *req.Content)
		}
		if *req.Visibility != "PRIVATE" {
			t.Errorf("expected visibility %q, got %q", "PRIVATE", *req.Visibility)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_full.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	content := "Updated content here"
	visibility := "PRIVATE"
	memo, err := c.UpdateMemo("5", &UpdateMemo{
		Name:       "memos/5",
		Content:    &content,
		Visibility: &visibility,
	}, "content,visibility")
	if err != nil {
		t.Fatalf("UpdateMemo: %v", err)
	}
	if memo.Name != "memos/1" {
		t.Errorf("expected name memos/1, got %s", memo.Name)
	}
}

func TestDeleteMemo(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/memos/99" {
			t.Errorf("expected /api/v1/memos/99, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	if err := c.DeleteMemo("99"); err != nil {
		t.Fatalf("DeleteMemo: %v", err)
	}
}

func TestStructuredError(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "api_error.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	_, err := c.GetMemo("nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", apiErr.StatusCode)
	}
	if apiErr.Code != 5 {
		t.Errorf("expected code 5, got %d", apiErr.Code)
	}
	if apiErr.Message != "Memo not found" {
		t.Errorf("expected message %q, got %q", "Memo not found", apiErr.Message)
	}
	if apiErr.Details == nil {
		t.Error("expected details to be present (even if empty)")
	}
}

func TestHTTPErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{"BadRequest", http.StatusBadRequest, `{"code": 3, "message": "Invalid request", "details": []}`},
		{"Unauthorized", http.StatusUnauthorized, `{"code": 16, "message": "Unauthorized", "details": []}`},
		{"Forbidden", http.StatusForbidden, `{"code": 7, "message": "Permission denied", "details": []}`},
		{"InternalServerError", http.StatusInternalServerError, `{"code": 13, "message": "Internal server error", "details": []}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := testServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.body))
			})
			defer server.Close()

			c := NewClient(server.URL, "test-token")
			_, err := c.GetMemo("1")
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("expected *APIError, got %T", err)
			}
			if apiErr.StatusCode != tt.statusCode {
				t.Errorf("expected status %d, got %d", tt.statusCode, apiErr.StatusCode)
			}
		})
	}
}

func TestPing(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(loadFixture(t, "memo_list_empty.json"))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	if err := c.Ping(); err != nil {
		t.Fatalf("Ping: %v", err)
	}
}

func TestPingFailure(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"code": 14, "message": "Service unavailable", "details": []}`))
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	if err := c.Ping(); err == nil {
		t.Fatal("expected Ping to fail, got nil")
	}
}
