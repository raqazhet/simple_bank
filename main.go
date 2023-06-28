package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"bank/api"
	"bank/gapi"
	"bank/jsonlog"
	"bank/pb"
	"bank/storage"
	"bank/util"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	// if err := runGinServer(config, store); err != nil {
	// 	logger.PrintError(err, map[string]string{
	// 		"runGinServerErr:": err.Error(),
	// 	})
	// 	return err
	// }
	if err := runGrpcServer(config, store); err != nil {
		logger.PrintError(err, map[string]string{
			"runGrpcServerErr:": err.Error(),
		})
	}
	return nil
}

func runGinServer(config util.Config, store storage.Store) error {
	l := jsonlog.Logger{}
	server := api.NewServer(config, store, &l)
	err := server.Start(config.HTTPServerAddress)
	if err != nil {
		return err
	}
	return nil
}

func runGrpcServer(config util.Config, store storage.Store) error {
	server := gapi.NewServer(config, store)
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Print("cannot create listener grpc")
		return err
	}
	log.Printf("start gRPC server at %s", listener.Addr())
	go runGateWayServer(config, store)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Printf("cannot start grpc server %v", err)
		return err
	}
	return nil
}

func runGateWayServer(conf util.Config, store storage.Store) {
	server := gapi.NewServer(conf, store)
	/*protoc генерирует JSON поля в camel case.
	Чтобы использовать тот же регистр, что и в proto файле, мы должны задать для параметра UseProtoNames значение true.
	И это можно сделать передавая MarshalerOption при создании gRPC мультиплексора.*/
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatalf("cannot register handler server %v", err)
	}
	// Этот мультиплексор фактически будет получать HTTP-запросы от клиентов.
	// Поэтому, чтобы преобразовать их в формат gRPC,
	// нам придется перенаправить их в gRPC мультиплексор, который мы создали ранее.
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)
	listener, err := net.Listen("tcp", conf.HTTPServerAddress)
	if err != nil {
		log.Fatalf("cannot create listener %s", err)
	}
	log.Printf("start HTTP gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("cannot start HTTP gateway server %s", err)
	}
}
