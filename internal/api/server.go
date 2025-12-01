package api

import (
	"github.com/example/ledger/internal/service"
	"github.com/labstack/echo/v4"
)

type Server struct {
	router  *echo.Echo
	service *service.Service
}

func NewServer(service *service.Service) *Server {
	server := &Server{
		service: service,
		router:  echo.New(),
	}

	server.setupRoutes()
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Start(address)
}

func (server *Server) setupRoutes() {
	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.POST("/transactions", server.createTransfer)
}
