package store

import (
	"testing"

	"github.com/google/uuid"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir())
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestCreateAndListConversations(t *testing.T) {
	s := newTestStore(t)

	c1 := &Conversation{ID: uuid.New().String(), Title: "First chat"}
	if err := s.CreateConversation(c1); err != nil {
		t.Fatalf("create: %v", err)
	}

	c2 := &Conversation{ID: uuid.New().String(), Title: "Second chat"}
	if err := s.CreateConversation(c2); err != nil {
		t.Fatalf("create: %v", err)
	}

	list, err := s.ListConversations()
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("expected 2, got %d", len(list))
	}
	if list[0].Title != "Second chat" {
		t.Errorf("expected 'Second chat' first, got '%s'", list[0].Title)
	}
}

func TestDeleteConversation_CascadesMessages(t *testing.T) {
	s := newTestStore(t)

	convID := uuid.New().String()
	s.CreateConversation(&Conversation{ID: convID, Title: "Test"})
	s.CreateMessage(&Message{ID: uuid.New().String(), ConversationID: convID, Role: "user", Content: "hello"})

	s.DeleteConversation(convID)

	msgs, _ := s.ListMessages(convID)
	if len(msgs) != 0 {
		t.Errorf("expected 0 messages after cascade delete, got %d", len(msgs))
	}
}
