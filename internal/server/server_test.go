package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"u-hermes/internal/chat"
	"u-hermes/internal/config"
	"u-hermes/internal/store"
)

func setupTestServer(t *testing.T) *Server {
	t.Helper()
	st, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	cfg := &config.Config{
		Version: 1,
		Models: []config.ModelConfig{
			{ID: "default", Name: "TestModel", APIBase: "http://localhost:9999", APIKey: "sk-test", IsDefault: true},
		},
		Chat: config.ChatConfig{SystemPrompt: "You are helpful."},
	}

	configDir := t.TempDir()
	configPath := configDir + "/config.json"

	chatSvc := chat.NewService()
	srv := New(st, cfg, chatSvc, configPath, nil, 0)
	return srv
}

func TestHealthEndpoint(t *testing.T) {
	srv := setupTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/health", nil)
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "ok" {
		t.Errorf("expected ok, got %v", resp["status"])
	}
	if resp["configured"] != true {
		t.Errorf("expected configured=true")
	}
}

func TestListModels_StripsAPIKey(t *testing.T) {
	srv := setupTestServer(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/models", nil)
	srv.Engine().ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var models []map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &models)
	if len(models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(models))
	}
	if _, exists := models[0]["api_key"]; exists {
		t.Error("api_key should not be exposed")
	}
	if models[0]["has_key"] != true {
		t.Error("has_key should be true")
	}
}

func TestConversationsCRUD(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/conversations", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 201 {
		t.Fatalf("create: expected 201, got %d", w.Code)
	}
	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	convID := created["id"].(string)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/conversations", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("list: expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/conversations/"+convID, nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("get: expected 200, got %d", w.Code)
	}

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/conversations/"+convID, nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("delete: expected 200, got %d", w.Code)
	}
}

func TestSettingsGetUpdate(t *testing.T) {
	srv := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/settings", nil)
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("get settings: %d", w.Code)
	}

	body := `{"chat":{"system_prompt":"new prompt"}}`
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/settings", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	srv.Engine().ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("update settings: %d: %s", w.Code, w.Body.String())
	}

	if srv.config.Chat.SystemPrompt != "new prompt" {
		t.Errorf("system prompt not updated: %s", srv.config.Chat.SystemPrompt)
	}
}
