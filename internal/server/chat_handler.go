package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"u-hermes/internal/chat"
	"u-hermes/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type chatRequest struct {
	ConversationID string `json:"conversation_id"`
	Message        string `json:"message"`
	Model          string `json:"model"`
}

func (s *Server) handleChat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	model := s.config.GetDefaultModel()
	if model == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no model configured", "code": "NO_MODEL"})
		return
	}
	if req.Model != "" {
		for _, m := range s.config.Models {
			if m.ID == req.Model {
				model = &m
				break
			}
		}
	}

	convID := req.ConversationID
	if convID == "" {
		convID = uuid.New().String()
		s.store.CreateConversation(&store.Conversation{ID: convID, Title: ""})
	}

	history, _ := s.store.ListMessages(convID)
	messages := make([]chat.Message, 0, len(history)+1)
	for _, m := range history {
		messages = append(messages, chat.Message{Role: m.Role, Content: m.Content})
	}
	messages = append(messages, chat.Message{Role: "user", Content: req.Message})

	if s.config.Chat.SystemPrompt != "" {
		messages = append([]chat.Message{{Role: "system", Content: s.config.Chat.SystemPrompt}}, messages...)
	}

	userMsgID := uuid.New().String()
	s.store.CreateMessage(&store.Message{
		ID: userMsgID, ConversationID: convID, Role: "user", Content: req.Message,
	})

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.WriteHeader(http.StatusOK)

	assistantMsgID := uuid.New().String()
	var fullContent string
	totalTokens := 0

	flusher, _ := c.Writer.(http.Flusher)

	tokens, err := s.chatSvc.Chat(model.APIBase, model.APIKey, model.Name, messages, func(token string) error {
		fullContent += token
		data, _ := json.Marshal(gin.H{"token": token})
		fmt.Fprintf(c.Writer, "event: token\ndata: %s\n\n", string(data))
		if flusher != nil {
			flusher.Flush()
		}
		return nil
	})
	totalTokens = tokens

	if err != nil {
		data, _ := json.Marshal(gin.H{"error": err.Error()})
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", string(data))
		if flusher != nil {
			flusher.Flush()
		}
		if fullContent != "" {
			s.store.CreateMessage(&store.Message{
				ID: assistantMsgID, ConversationID: convID,
				Role: "assistant", Content: fullContent, TokensUsed: totalTokens,
			})
		}
		return
	}

	s.store.CreateMessage(&store.Message{
		ID: assistantMsgID, ConversationID: convID,
		Role: "assistant", Content: fullContent, TokensUsed: totalTokens,
	})

	conv, _ := s.store.GetConversation(convID)
	if conv != nil && conv.Title == "" {
		title := req.Message
		if len([]rune(title)) > 30 {
			title = string([]rune(title)[:30]) + "..."
		}
		s.store.UpdateConversationTitle(convID, title)
	}

	data, _ := json.Marshal(gin.H{
		"total_tokens":    totalTokens,
		"conversation_id": convID,
		"message_id":      assistantMsgID,
	})
	fmt.Fprintf(c.Writer, "event: done\ndata: %s\n\n", string(data))
	if flusher != nil {
		flusher.Flush()
	}
}
