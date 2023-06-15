package api

import (
	"bank/jsonlog"
	"bank/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *storage.Storage
	router *gin.Engine
	l      *jsonlog.Logger
}

func NewServer(store *storage.Storage, l *jsonlog.Logger) *Server {
	server := &Server{store: store, l: l}
	router := gin.Default()
	account := router.Group("/v1")
	{
		account.POST("/accounts", server.CreateAccount)
		account.GET("/accounts/:id", server.getAccountByID)
	}
	server.router = router
	return server
}

func (srv *Server) Start(addres string) error {
	return srv.router.Run(addres)
}
