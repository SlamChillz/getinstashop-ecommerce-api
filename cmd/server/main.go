package server

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/routers"
)

// Server struct
type Server struct {
	config config.Config
	router *gin.Engine
}

// NewServer Create a new server instance
func NewServer(config config.Config) (*Server, error) {
	server := &Server{config: config}
	server.setupRouter()
	return server, nil
}

// Register application routers
func (server *Server) setupRouter() {
	server.router = routers.InitRouters()
}

// Start server on the given address
func (server *Server) Start() error {
	return server.router.Run(server.config.HTTPServerAddress)
}
