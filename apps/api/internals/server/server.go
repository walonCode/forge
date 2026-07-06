package server

import (
	"api/internals/middleware"
	"api/pkg/configs"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api/docs"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	router *chi.Mux
	db     *sql.DB
	config *configs.Config
	logger *slog.Logger
}

func New(db *sql.DB, cfg *configs.Config, logger *slog.Logger) *Server {
	s := &Server{
		router: chi.NewRouter(),
		db:     db,
		config: cfg,
		logger: logger,
	}

	s.mountMiddleware()
	s.mountRoutes()

	return s
}

func (s *Server) mountMiddleware() {
	s.router.Use(chimiddleware.RealIP)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"X-Correlation-ID"},
		// token auth uses the Authorization header, not cookies; credentialed
		// mode is incompatible with a wildcard origin, so keep it off.
		AllowCredentials: false,
		MaxAge:           300,
	}))
	s.router.Use(middleware.MaxBodyBytes)
	s.router.Use(middleware.CreateCorrelationIdMiddleware)
	s.router.Use(middleware.LoggerMiddleware(s.logger))
}

// add soon
func (s *Server) mountRoutes() {
	//metrics
	s.router.Handle("/metrics", promhttp.Handler())

	//doc
	s.router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", s.config.Port)),
	))

	//mount the modules
	s.mountModules()
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Port),
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	//channel to recieve OS signal
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	//channel to recieve server error
	serverErr := make(chan error, 1)

	go func() {
		s.logger.Info("server starting", "addr", srv.Addr, "version", s.config.AppVersion)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}

	case sig := <-shutdown:
		s.logger.Info("shutdown signal recieved", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			_ = srv.Close()
			return fmt.Errorf("server shutdown gracefully failed: %w", err)
		}

		s.logger.Info("server stopped gracefully")
	}

	return nil
}
