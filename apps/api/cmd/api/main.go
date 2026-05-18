package main

import (
	"api/internals/server"
	"database/sql"
	"log/slog"
	"os"
)

func main(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//real db and cfg comming soon
	var cfg any
	var db *sql.DB
	srv := server.New(db, cfg)

	if err := srv.Start(); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}