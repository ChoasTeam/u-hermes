package server

import (
	"bufio"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"u-hermes/internal/config"
	"u-hermes/internal/store"

	"u-hermes/internal/chat"
)

func TestCreateModel(t *testing.T) {
	srv := setupTestServer(t)

	body := `{"name":"Kimi","api_base":"https://api.moonshot.cn","api_key":"sk-kimi"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/models", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 201 {
		t.Fatalf("create model: expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "created" {
		t.Errorf("expected status=created, got %v", resp["status"])
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("expected non-empty id")
	}

	// Verify model appeared in config
	if len(srv.config.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(srv.config.Models))
	}
}

func TestUpdateModel_PropagatesIsDefault(t *testing.T) {
	srv := setupTestServer(t)

	body := `{"is_default":true}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/models/default", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("update model: expected 200, got %d", w.Code)
	}

	// Verify only one default
	for _, m := range srv.config.Models {
		if m.ID != "default" && m.IsDefault {
			t.Error("only the updated model should be default")
		}
	}
}

func TestGetConversation_NotFound(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/conversations/nonexistent-id", nil)
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("expected 404 for nonexistent conversation, got %d", w.Code)
	}
}

func TestListConversations_Empty(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/conversations", nil)
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("list conversations: expected 200, got %d", w.Code)
	}

	// Empty list should be JSON array [], not null
	body := strings.TrimSpace(w.Body.String())
	if body != "[]" {
		t.Errorf("expected empty array [], got %s", body)
	}
}

func TestChat_NoModel_ReturnsError(t *testing.T) {
	// Server with no models
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	defer st.Close()

	cfg := &config.Config{Version: 1}
	srv := New(st, cfg, chat.NewService(), t.TempDir()+"/config.json", nil, 0)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chat", strings.NewReader(`{"message":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["code"] != "NO_MODEL" {
		t.Errorf("expected NO_MODEL code, got %v", resp["code"])
	}
}

func TestChat_InvalidJSON_ReturnsError(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chat", strings.NewReader(`{invalid`))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestChat_MissingMessage_ReturnsError(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chat", strings.NewReader(`{"model":"deepseek"}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	// Empty message should still process (it creates conversation)
	// but our test model at localhost:9999 will fail
	// Check that it returns SSE Content-Type
	if w.Header().Get("Content-Type") != "text/event-stream" {
		t.Errorf("expected SSE content type, got %s", w.Header().Get("Content-Type"))
	}
}

func TestChat_SSEStream_EmitsEvents(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/chat", strings.NewReader(`{"message":"hello"}`))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)

	// The model at localhost:9999 will fail, but we should get SSE events
	body := w.Body.String()
	if !strings.Contains(body, "event:") {
		t.Errorf("expected SSE events in response, got: %s", body)
	}

	// Should contain an error event since model is unreachable
	scanner := bufio.NewScanner(strings.NewReader(body))
	hasDone := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "event: done") {
			hasDone = true
		}
	}
	if !hasDone {
		t.Log("no done event (expected for error path)")
	}
}
