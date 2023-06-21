package api

import (
	"bank/jsonlog"
	"bank/storage"
	"bank/token"
	"bank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      storage.Store
	router     *gin.Engine
	l          *jsonlog.Logger
	tokenMaker token.Maker
}

func NewServer(conf util.Config, store storage.Store, l *jsonlog.Logger) *Server {
	tokenMaker, err := token.NewJwtMaker(util.RandomString(32))
	if err != nil {
		return nil
	}
	server := &Server{
		config:     conf,
		store:      store,
		l:          l,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	return server
}

func (srv *Server) Start(addres string) error {
	return srv.router.Run(addres)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(gin.Recovery())
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	account := router.Group("/v1/accounts").Use(authMidddleware(server.tokenMaker))
	{
		account.POST("/", server.CreateAccount)
		account.GET("/:id", server.getAccountByID)
		account.GET("/", server.getAllAccounts)
		account.POST("/transfers", server.createTransfer)
	}
	user := router.Group("/v1/user")
	{
		user.POST("/", server.CreateUser)
		user.POST("/login", server.LoginUser)
	}
	server.router = router
}
