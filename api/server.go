package api

import (
	"bank/jsonlog"
	"bank/storage"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  storage.Store
	router *gin.Engine
	l      *jsonlog.Logger
}

func NewServer(store storage.Store, l *jsonlog.Logger) *Server {
	server := &Server{store: store, l: l}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	account := router.Group("/v1")
	{
		account.POST("/accounts", server.CreateAccount)
		account.GET("/accounts/:id", server.getAccountByID)
		account.GET("/accounts", server.getAllAccounts)
		account.POST("/transfers", server.createTransfer)
	}
	server.router = router
	return server
}

func (srv *Server) Start(addres string) error {
	return srv.router.Run(addres)
}
