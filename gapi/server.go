package gapi

import (
	"bank/jsonlog"
	"bank/pb"
	"bank/storage"
	"bank/token"
	"bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      storage.Store
	l          *jsonlog.Logger
	tokenMaker token.Maker
}

func NewServer(conf util.Config, store storage.Store, l *jsonlog.Logger) *Server {
	tokenMaker, err := token.NewJwtMaker(conf.TokenKey)
	if err != nil {
		return nil
	}
	server := &Server{
		config:     conf,
		store:      store,
		l:          l,
		tokenMaker: tokenMaker,
	}
	return server
}
