package main

import (
	"api/internals/server"
	"api/pkg/configs"
	"api/pkg/database"
	"log/slog"
	"os"
)

func main(){
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//real db and cfg comming soon
	cfg := configs.Load()
	db, err := database.Connect(cfg.DATABSE_URL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	
	srv := server.New(db.Client, cfg)

	if err := srv.Start(); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}