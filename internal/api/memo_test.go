package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func TestCreateMemo(t *testing.T) {
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
		if req.Content != "hello world" {
			t.Errorf("expected content 'hello world', got %q", req.Content)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Memo{
			Name:       "memos/1",
			Content:    "hello world",
			Visibility: "PRIVATE",
		})
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.CreateMemo(&CreateMemo{
		Content:    "hello world",
		Visibility: "PRIVATE",
	})
	if err != nil {
		t.Fatalf("CreateMemo: %v", err)
	}
	if memo.Name != "memos/1" {
		t.Errorf("expected name memos/1, got %s", memo.Name)
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
		json.NewEncoder(w).Encode(Memo{
			Name:    "memos/42",
			Content: "test content",
		})
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.GetMemo("42")
	if err != nil {
		t.Fatalf("GetMemo: %v", err)
	}
	if memo.Content != "test content" {
		t.Errorf("expected 'test content', got %q", memo.Content)
	}
}

func TestGetMemoWithPrefix(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/memos/42" {
			t.Errorf("expected /api/v1/memos/42, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Memo{Name: "memos/42"})
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	memo, err := c.GetMemo("memos/42")
	if err != nil {
		t.Fatalf("GetMemo: %v", err)
	}
	if memo.Name != "memos/42" {
		t.Errorf("expected name memos/42, got %s", memo.Name)
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
		json.NewEncoder(w).Encode(ListMemosResponse{
			Memos: []Memo{
				{Name: "memos/1", Content: "first"},
				{Name: "memos/2", Content: "second"},
			},
			NextPageToken: "next",
		})
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
	if next != "next" {
		t.Errorf("expected next page token 'next', got %q", next)
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

		var req UpdateMemo
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if *req.Content != "updated" {
			t.Errorf("expected content 'updated', got %q", *req.Content)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Memo{
			Name:    "memos/5",
			Content: "updated",
		})
	})
	defer server.Close()

	c := NewClient(server.URL, "test-token")
	content := "updated"
	memo, err := c.UpdateMemo("5", &UpdateMemo{Name: "memos/5", Content: &content}, "content")
	if err != nil {
		t.Fatalf("UpdateMemo: %v", err)
	}
	if memo.Content != "updated" {
		t.Errorf("expected 'updated', got %q", memo.Content)
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

func TestAPIError(t *testing.T) {
	server := testServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "memo not found"}`))
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
}
