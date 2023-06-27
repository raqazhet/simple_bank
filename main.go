package main

import (
	"database/sql"
	"log"
	"os"

	"bank/api"
	"bank/gapi"
	"bank/jsonlog"
	"bank/pb"
	"bank/storage"
	"bank/util"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("cannot connect to :", err)
	}
}

func run() error {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	config, err := util.LoadConfig("./")
	if err != nil {
		logger.PrintError(err, map[string]string{
			"load config": err.Error(),
		})
		return err
	}

	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.PrintError(err, map[string]string{
			"connect to db": err.Error(),
		})
		return err
	}
	store := storage.NewStorage(db)
	if err := runGinServer(config, store); err != nil {
		logger.PrintError(err, map[string]string{
			"runGinServerErr:": err.Error(),
		})
		return err
	}
	return nil
}

func runGinServer(config util.Config, store storage.Store) error {
	l := jsonlog.Logger{}
	server := api.NewServer(config, store, &l)
	err := server.Start(config.ServerAddress)
	if err != nil {
		return err
	}
	return nil
}

func runGrpcServer(config util.Config, store storage.Store) {
	server := gapi.NewServer(config, store)
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
}
