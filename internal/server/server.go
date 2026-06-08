package server

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"u-hermes/internal/chat"
	"u-hermes/internal/config"
	"u-hermes/internal/store"

	"github.com/gin-gonic/gin"
)

type Server struct {
	engine     *gin.Engine
	store      *store.Store
	config     *config.Config
	chatSvc    *chat.Service
	configPath string
	port       int
}

func New(st *store.Store, cfg *config.Config, chatSvc *chat.Service, configPath string, webAssets fs.FS, port int) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	s := &Server{
		engine:     engine,
		store:      st,
		config:     cfg,
		chatSvc:    chatSvc,
		configPath: configPath,
		port:       port,
	}

	api := engine.Group("/api")
	{
		api.GET("/health", s.handleHealth)
		api.POST("/chat", s.handleChat)
		api.GET("/conversations", s.handleListConversations)
		api.POST("/conversations", s.handleCreateConversation)
		api.GET("/conversations/:id", s.handleGetConversation)
		api.DELETE("/conversations/:id", s.handleDeleteConversation)
		api.POST("/models", s.handleCreateModel)
		api.GET("/models", s.handleListModels)
		api.PUT("/models/:id", s.handleUpdateModel)
		api.POST("/models/:id/test", s.handleTestModel)
		api.GET("/settings", s.handleGetSettings)
		api.PUT("/settings", s.handleUpdateSettings)
		api.POST("/settings/reset", s.handleResetSettings)
	}

	// Serve embedded SPA or fallback to disk
	distFS, err := fs.Sub(webAssets, "web/dist")
	if err != nil {
		engine.Static("/assets", "./web/dist/assets")
		engine.StaticFile("/", "./web/dist/index.html")
		engine.NoRoute(func(c *gin.Context) {
			c.File("./web/dist/index.html")
		})
	} else {
		engine.StaticFS("/assets", mustSub(distFS, "assets"))
		engine.GET("/", func(c *gin.Context) {
			c.FileFromFS("index.html", http.FS(distFS))
		})
		engine.NoRoute(func(c *gin.Context) {
			c.FileFromFS("index.html", http.FS(distFS))
		})
	}

	return s
}

func mustSub(fsys fs.FS, dir string) http.FileSystem {
	sub, err := fs.Sub(fsys, dir)
	if err != nil {
		panic(err)
	}
	return http.FS(sub)
}

func (s *Server) Engine() *gin.Engine {
	return s.engine
}

func (s *Server) ConfigPath() string {
	return s.configPath
}

func (s *Server) DataDir() string {
	return filepath.Dir(s.configPath) + "/data"
}
