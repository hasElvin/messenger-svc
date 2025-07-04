package rest

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/hasElvin/messenger-svc/docs"
	"github.com/hasElvin/messenger-svc/internal/adapters/rest/handlers"
	"github.com/hasElvin/messenger-svc/internal/core/ports"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	messageHandler *handlers.MessageHandler
	utilityHandler *handlers.UtilityHandler
	router         *gin.Engine
}

func NewServer(messageService ports.MessageService, utilityService ports.UtilityService) *Server {
	messageHandler := handlers.NewMessageHandler(messageService)
	utilityHandler := handlers.NewUtilityHandler(utilityService)

	router := gin.Default()
	router.Use(cors.Default())

	server := &Server{
		messageHandler: messageHandler,
		utilityHandler: utilityHandler,
		router:         router,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.router.POST("/start", s.messageHandler.StartAutoSender)
	s.router.POST("/stop", s.messageHandler.StopAutoSender)
	s.router.GET("/sent", s.messageHandler.GetSentMessages)

	s.router.GET("/ping", s.utilityHandler.Ping)
	s.router.POST("/seed", s.utilityHandler.SeedSampleMessages)
	s.router.DELETE("/clear", s.utilityHandler.ClearDatabase)

	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
