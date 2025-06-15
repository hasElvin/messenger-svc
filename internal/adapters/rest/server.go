package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/hasElvin/messenger-svc/internal/adapters/rest/handlers"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
)

type Server struct {
	messageHandler *handlers.MessageHandler
	router         *gin.Engine
}

func NewServer(messageService ports.MessageService) *Server {
	messageHandler := handlers.NewMessageHandler(messageService)

	router := gin.Default()

	server := &Server{
		messageHandler: messageHandler,
		router:         router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.GET("/ping", s.messageHandler.Ping)
	s.router.POST("/start", s.messageHandler.StartAutoSender)
	s.router.POST("/stop", s.messageHandler.StopAutoSender)
	s.router.GET("/sent", s.messageHandler.GetSentMessages)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
