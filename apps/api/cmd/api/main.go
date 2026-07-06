// Package main is the entry point of the Forge API.
//
//	@title			Forge API
//	@version		1.0
//	@description	Task management API for Forge.
//	@host			localhost:8080
//	@BasePath		/
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by your token, e.g. "Bearer eyJhbGci..."
package main

import (
	"api/internals/server"
	"api/pkg/configs"
	"api/pkg/database"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := configs.Load()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := server.New(db.Client, cfg, logger)

	if err := srv.Start(); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}
