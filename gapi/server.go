package gapi

import (
	"bank/pb"
	"bank/storage"
	"bank/token"
	"bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      storage.Store
	tokenMaker token.Maker
}

func NewServer(conf util.Config, store storage.Store) *Server {
	tokenMaker, err := token.NewJwtMaker(conf.TokenKey)
	if err != nil {
		return nil
	}
	server := &Server{
		config:     conf,
		store:      store,
		tokenMaker: tokenMaker,
	}
	return server
}
