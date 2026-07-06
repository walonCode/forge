package health

import (
	"context"
	"database/sql"
	"time"
)

type Service struct {
	db      *sql.DB
	version string
}

var start = time.Now()

func newService(db *sql.DB, version string) *Service {
	return &Service{
		db:      db,
		version: version,
	}
}

func (s *Service) Health() HealthResponse {
	return HealthResponse{
		Status:  "ok",
		Version: s.version,
		Uptime:  time.Since(start).Round(time.Second).String(),
	}
}

func (s *Service) Ready(ctx context.Context) (ReadinessResponse, bool) {
	checks := map[string]string{}
	ok := true

	// check database connectivity
	if err := s.db.PingContext(ctx); err != nil {
		checks["db"] = "unreachable"
		ok = false
	} else {
		checks["db"] = "ok"
	}

	status := "ready"
	if !ok {
		status = "unavailable"
	}

	return ReadinessResponse{
		Status: status,
		Checks: checks,
	}, ok
}
