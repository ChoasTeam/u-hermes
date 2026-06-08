package store

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func (s *Store) CreateConversation(conv *Conversation) error {
	now := time.Now().Unix()
	if conv.CreatedAt == 0 {
		conv.CreatedAt = now
	}
	if conv.UpdatedAt == 0 {
		conv.UpdatedAt = now
	}
	_, err := s.db.Exec(
		"INSERT INTO conversations (id, title, created_at, updated_at) VALUES (?, ?, ?, ?)",
		conv.ID, conv.Title, conv.CreatedAt, conv.UpdatedAt,
	)
	return err
}

func (s *Store) ListConversations() ([]Conversation, error) {
	rows, err := s.db.Query(
		"SELECT id, title, created_at, updated_at FROM conversations ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convs []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		convs = append(convs, c)
	}
	return convs, rows.Err()
}

func (s *Store) GetConversation(id string) (*Conversation, error) {
	var c Conversation
	err := s.db.QueryRow(
		"SELECT id, title, created_at, updated_at FROM conversations WHERE id = ?", id,
	).Scan(&c.ID, &c.Title, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) UpdateConversationTitle(id, title string) error {
	_, err := s.db.Exec(
		"UPDATE conversations SET title = ?, updated_at = ? WHERE id = ?",
		title, time.Now().Unix(), id,
	)
	return err
}

func (s *Store) DeleteConversation(id string) error {
	_, err := s.db.Exec("DELETE FROM conversations WHERE id = ?", id)
	return err
}
