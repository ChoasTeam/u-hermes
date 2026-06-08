package chat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChat_StreamsTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer sk-test" {
			t.Errorf("missing auth header")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		chunks := []string{
			`data: {"choices":[{"delta":{"content":"Hello"}}]}`,
			`data: {"choices":[{"delta":{"content":" world"}}]}`,
			`data: {"choices":[{"delta":{"content":"!"}}],"usage":{"total_tokens":15}}`,
			`data: [DONE]`,
		}
		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n\n"))
			w.(http.Flusher).Flush()
		}
	}))
	defer server.Close()

	svc := NewService()
	var collected string
	tokens, err := svc.Chat(server.URL, "sk-test", "test-model", []Message{
		{Role: "user", Content: "hi"},
	}, func(token string) error {
		collected += token
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if collected != "Hello world!" {
		t.Errorf("expected 'Hello world!', got '%s'", collected)
	}
	if tokens != 15 {
		t.Errorf("expected 15 tokens, got %d", tokens)
	}
}

func TestChat_APIError_ReturnsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`{"error":"invalid api key"}`))
	}))
	defer server.Close()

	svc := NewService()
	_, err := svc.Chat(server.URL, "bad-key", "model", []Message{
		{Role: "user", Content: "hi"},
	}, func(token string) error { return nil })

	if err == nil || !strings.Contains(err.Error(), "401") {
		t.Errorf("expected 401 error, got: %v", err)
	}
}

func TestTestConnection_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.Stream {
			t.Error("test connection should use stream=false")
		}
		w.WriteHeader(200)
		w.Write([]byte(`{"choices":[{"message":{"content":"hi"}}]}`))
	}))
	defer server.Close()

	svc := NewService()
	if err := svc.TestConnection(server.URL, "sk-test", "model"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTestConnection_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer server.Close()

	svc := NewService()
	if err := svc.TestConnection(server.URL, "sk-test", "model"); err == nil {
		t.Error("expected error for 500 response")
	}
}
