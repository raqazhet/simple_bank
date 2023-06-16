package main

import (
	"database/sql"
	"log"
	"os"

	"bank/api"
	"bank/jsonlog"
	"bank/storage"
	"bank/util"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal("cannot connect to db:", err)
	}
}

func run() error {
	config, err := util.LoadConfig("./")
	if err != nil {
		return err
	}
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.PrintError(err, map[string]string{
			"connect to db": err.Error(),
		})
		return err
	}
	store := storage.NewStorage(db)
	server := api.NewServer(store, logger)
	if err := server.Start(config.ServerAddress); err != nil {
		logger.PrintError(err, map[string]string{
			"start server": err.Error(),
		})
		return err
	}

	return nil
}
