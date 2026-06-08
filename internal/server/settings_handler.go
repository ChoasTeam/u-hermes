package server

import (
	"net/http"
	"os"

	"u-hermes/internal/config"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleGetSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"chat": gin.H{
			"system_prompt": s.config.Chat.SystemPrompt,
		},
		"version": s.config.Version,
	})
}

func (s *Server) handleUpdateSettings(c *gin.Context) {
	var update struct {
		Chat *struct {
			SystemPrompt string `json:"system_prompt"`
		} `json:"chat"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if update.Chat != nil {
		s.config.Chat.SystemPrompt = update.Chat.SystemPrompt
	}
	if err := config.Save(s.configPath, s.config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (s *Server) handleResetSettings(c *gin.Context) {
	os.Remove(s.configPath)
	os.RemoveAll(s.DataDir())

	c.JSON(http.StatusOK, gin.H{"status": "reset"})

	go func() {
		os.Exit(0)
	}()
}
