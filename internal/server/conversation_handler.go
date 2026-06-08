package server

import (
	"net/http"

	"u-hermes/internal/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) handleListConversations(c *gin.Context) {
	convs, err := s.store.ListConversations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if convs == nil {
		convs = []store.Conversation{}
	}
	c.JSON(http.StatusOK, convs)
}

func (s *Server) handleCreateConversation(c *gin.Context) {
	conv := &store.Conversation{ID: uuid.New().String()}
	if err := s.store.CreateConversation(conv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, conv)
}

func (s *Server) handleGetConversation(c *gin.Context) {
	id := c.Param("id")
	conv, err := s.store.GetConversation(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if conv == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	msgs, _ := s.store.ListMessages(id)
	if msgs == nil {
		msgs = []store.Message{}
	}
	c.JSON(http.StatusOK, gin.H{"conversation": conv, "messages": msgs})
}

func (s *Server) handleDeleteConversation(c *gin.Context) {
	id := c.Param("id")
	if err := s.store.DeleteConversation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
