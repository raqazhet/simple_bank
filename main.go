package main

import (
	"database/sql"
	"log"
	"os"

	"bank/api"
	"bank/jsonlog"
	"bank/storage"

	_ "github.com/lib/pq"
)

const (
	dbdriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/test?sslmode=disable"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("cannot connect to db:", err)
	}
}

func run() error {
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := sql.Open(dbdriver, dbSource)
	if err != nil {
		logger.PrintError(err, map[string]string{
			"connect to db": err.Error(),
		})
		return err
	}
	store := storage.NewStorage(db)
	server := api.NewServer(store, logger)
	if err := server.Start(":4000"); err != nil {
		logger.PrintError(err, map[string]string{
			"start server": err.Error(),
		})
		return err
	}
	return nil
}
