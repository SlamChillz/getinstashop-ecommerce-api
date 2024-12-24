package server

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/routers"
)

// Server struct
type Server struct {
	config  config.Config
	router  *gin.Engine
	store   db.Store
	handler *handlers.AllHandler
}

// NewServer Create a new server instance
func NewServer(config config.Config, store db.Store) (*Server, error) {
	server := &Server{config: config, store: store}
	server.setupHandler().setupRouter()
	return server, nil
}

// Instantiate all handlers
func (server *Server) setupHandler() *Server {
	server.handler = handlers.RegisterHandlers(server.store)
	return server
}

// Register application routers
func (server *Server) setupRouter() *Server {
	server.router = routers.InitRouters(server.handler)
	return server
}

// Start server on the given address
func (server *Server) Start() error {
	return server.router.Run(server.config.HTTPServerAddress)
}
