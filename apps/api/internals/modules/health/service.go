package health

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	db *sql.DB
}

var start = time.Now()

func newService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) Health() HealthResponse {
	return HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
		Uptime:  fmt.Sprintf("%s", time.Since(start).Round(time.Second)),
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
