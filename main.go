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
	server := api.NewServer(config, store, logger)
	if err := server.Start(config.ServerAddress); err != nil {
		logger.PrintError(err, map[string]string{
			"start server": err.Error(),
		})
		return err
	}

	return nil
}
