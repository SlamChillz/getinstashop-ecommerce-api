package server

import (
	"github.com/gin-gonic/gin"
	"github.com/slamchillz/getinstashop-ecommerce-api/config"
	db "github.com/slamchillz/getinstashop-ecommerce-api/internal/db/sqlc"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/handlers"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/routers"
	"github.com/slamchillz/getinstashop-ecommerce-api/pkg/token"
)

// Server struct
type Server struct {
	config  config.Config
	token   *token.JWT
	router  *gin.Engine
	store   db.Store
	handler *handlers.AllHandler
}

// NewServer Create a new server instance
func NewServer(config config.Config, store db.Store) (*Server, error) {
	jwt, err := token.NewJWT(config.JwtSecret)
	if err != nil {
		return nil, err
	}
	server := &Server{config: config, token: jwt, store: store}
	server.setupHandler().setupRouter()
	return server, nil
}

// Instantiate all handlers
func (server *Server) setupHandler() *Server {
	server.handler = handlers.RegisterHandlers(server.store, server.token)
	return server
}

// Register application routers
func (server *Server) setupRouter() *Server {
	server.router = routers.InitRouters(server.handler, server.token)
	return server
}

func (server *Server) Router() *gin.Engine {
	return server.router
}

func (server *Server) TokenCreator() *token.JWT {
	return server.token
}

// Start server on the given address
func (server *Server) Start() error {
	return server.router.Run(server.config.HTTPServerAddress)
}
