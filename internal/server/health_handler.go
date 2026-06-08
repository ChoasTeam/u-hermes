package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func (s *Server) handleHealth(c *gin.Context) {
	modelConnected := false
	modelName := ""
	if m := s.config.GetDefaultModel(); m != nil {
		modelName = m.Name
		if err := s.chatSvc.TestConnection(m.APIBase, m.APIKey, m.Name); err == nil {
			modelConnected = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":          "ok",
		"version":         "0.1.0",
		"model":           modelName,
		"model_connected": modelConnected,
		"configured":      len(s.config.Models) > 0,
		"uptime":          int(time.Since(startTime).Seconds()),
	})
}
