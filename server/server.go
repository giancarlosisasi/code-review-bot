package server

import (
	"fmt"
	"net/http"

	"github.com/giancarlosisasi/code-review-bot/config"
	"github.com/giancarlosisasi/code-review-bot/database"
	"github.com/gin-gonic/gin"
)

type Server struct {
	inMemoryDatabase *database.InMemoryDatabase
	config           *config.Config
}

func NewServer(inMemoryDatabase *database.InMemoryDatabase, config *config.Config) *Server {
	return &Server{
		inMemoryDatabase: inMemoryDatabase,
		config:           config,
	}
}

func (s *Server) Run() error {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	if s.config.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	err := router.Run(fmt.Sprintf(":%d", s.config.Port))

	return err
}
