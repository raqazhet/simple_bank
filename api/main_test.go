package api

import (
	"testing"
	"time"

	"bank/jsonlog"
	"bank/storage"
	"bank/util"
)

func newTestServer(t *testing.T, store storage.Store) *Server {
	config := util.Config{
		AccessTokenDuration: time.Minute,
	}
	logger := jsonlog.Logger{}
	server := NewServer(config, store, &logger)

	return server
}
