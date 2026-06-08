package store

import "time"

type Message struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Role           string `json:"role"`
	Content        string `json:"content"`
	TokensUsed     int    `json:"tokens_used"`
	CreatedAt      int64  `json:"created_at"`
}

func (s *Store) CreateMessage(msg *Message) error {
	if msg.CreatedAt == 0 {
		msg.CreatedAt = time.Now().UnixMilli()
	}
	_, err := s.db.Exec(
		"INSERT INTO messages (id, conversation_id, role, content, tokens_used, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		msg.ID, msg.ConversationID, msg.Role, msg.Content, msg.TokensUsed, msg.CreatedAt,
	)
	return err
}

func (s *Store) ListMessages(conversationID string) ([]Message, error) {
	rows, err := s.db.Query(
		"SELECT id, conversation_id, role, content, tokens_used, created_at FROM messages WHERE conversation_id = ? ORDER BY created_at ASC, rowid ASC",
		conversationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.Role, &m.Content, &m.TokensUsed, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (s *Store) UpdateMessageContent(id, content string, tokensUsed int) error {
	_, err := s.db.Exec(
		"UPDATE messages SET content = ?, tokens_used = ? WHERE id = ?",
		content, tokensUsed, id,
	)
	return err
}
