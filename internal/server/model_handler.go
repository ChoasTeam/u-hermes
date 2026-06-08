package server

import (
	"net/http"

	"u-hermes/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) handleCreateModel(c *gin.Context) {
	var newModel config.ModelConfig
	if err := c.ShouldBindJSON(&newModel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if newModel.ID == "" {
		newModel.ID = uuid.New().String()
	}
	s.config.Models = append(s.config.Models, newModel)

	if err := config.Save(s.configPath, s.config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save config"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created", "id": newModel.ID})
}

func (s *Server) handleListModels(c *gin.Context) {
	models := s.config.Models
	if models == nil {
		models = []config.ModelConfig{}
	}
	safe := make([]gin.H, len(models))
	for i, m := range models {
		safe[i] = gin.H{
			"id":         m.ID,
			"name":       m.Name,
			"api_base":   m.APIBase,
			"is_default": m.IsDefault,
			"has_key":    m.APIKey != "",
		}
	}
	c.JSON(http.StatusOK, safe)
}

func (s *Server) handleUpdateModel(c *gin.Context) {
	id := c.Param("id")
	var update config.ModelConfig
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	for i, m := range s.config.Models {
		if m.ID == id {
			if update.Name != "" {
				s.config.Models[i].Name = update.Name
			}
			if update.APIBase != "" {
				s.config.Models[i].APIBase = update.APIBase
			}
			if update.APIKey != "" {
				s.config.Models[i].APIKey = update.APIKey
			}
			s.config.Models[i].IsDefault = update.IsDefault

			if update.IsDefault {
				for j := range s.config.Models {
					if j != i {
						s.config.Models[j].IsDefault = false
					}
				}
			}

			if err := config.Save(s.configPath, s.config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save config"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"status": "updated"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
}

func (s *Server) handleTestModel(c *gin.Context) {
	id := c.Param("id")
	var model *config.ModelConfig
	for _, m := range s.config.Models {
		if m.ID == id {
			model = &m
			break
		}
	}
	if model == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "model not found"})
		return
	}

	if err := s.chatSvc.TestConnection(model.APIBase, model.APIKey, model.Name); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
