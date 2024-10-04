package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Jacobbrewer1/vaulty/repositories"
	"github.com/alexliesenfeld/health"
)

func healthHandler(db *repositories.Database) http.Handler {
	checker := health.NewChecker(
		// Disable caching of the results of the checks.
		health.WithCacheDuration(0),
		health.WithDisabledCache(),

		// Set a timeout of 10 seconds for the entire health check.
		health.WithTimeout(10*time.Second),

		// Monitor the health of the database.
		health.WithCheck(health.Check{
			Name: "database",
			Check: func(ctx context.Context) error {
				if err := db.PingContext(ctx); err != nil {
					return fmt.Errorf("failed to ping database: %w", err)
				}
				return nil
			},
			Timeout:            3 * time.Second,
			MaxTimeInError:     0,
			MaxContiguousFails: 0,
			StatusListener: func(ctx context.Context, name string, state health.CheckState) {
				slog.Info("database health check status changed",
					slog.String("name", name),
					slog.String("state", string(state.Status)),
				)
			},
			Interceptors:         nil,
			DisablePanicRecovery: false,
		}),
	)

	return health.NewHandler(checker)
}
