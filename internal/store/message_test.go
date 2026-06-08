package store

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateAndListMessages(t *testing.T) {
	s := newTestStore(t)

	convID := uuid.New().String()
	s.CreateConversation(&Conversation{ID: convID, Title: "Test"})

	m1 := &Message{ID: uuid.New().String(), ConversationID: convID, Role: "user", Content: "hello"}
	m2 := &Message{ID: uuid.New().String(), ConversationID: convID, Role: "assistant", Content: "hi there"}

	if err := s.CreateMessage(m1); err != nil {
		t.Fatalf("create m1: %v", err)
	}
	if err := s.CreateMessage(m2); err != nil {
		t.Fatalf("create m2: %v", err)
	}

	msgs, err := s.ListMessages(convID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(msgs) != 2 {
		t.Errorf("expected 2, got %d", len(msgs))
	}
	if msgs[0].Role != "user" {
		t.Errorf("expected user first, got %s", msgs[0].Role)
	}
}

func TestUpdateMessageContent(t *testing.T) {
	s := newTestStore(t)

	convID := uuid.New().String()
	s.CreateConversation(&Conversation{ID: convID, Title: "Test"})

	msgID := uuid.New().String()
	s.CreateMessage(&Message{ID: msgID, ConversationID: convID, Role: "assistant", Content: "partial"})

	s.UpdateMessageContent(msgID, "complete response", 42)

	msgs, _ := s.ListMessages(convID)
	if msgs[0].Content != "complete response" {
		t.Errorf("expected 'complete response', got '%s'", msgs[0].Content)
	}
	if msgs[0].TokensUsed != 42 {
		t.Errorf("expected 42 tokens, got %d", msgs[0].TokensUsed)
	}
}
