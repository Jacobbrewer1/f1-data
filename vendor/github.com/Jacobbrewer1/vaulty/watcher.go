package vaulty

import (
	"context"
	"fmt"
	"log/slog"

	hashiVault "github.com/hashicorp/vault/api"
)

func monitorWatcher(ctx context.Context, name string, watcher *hashiVault.LifetimeWatcher) (renewResult, error) {
	for {
		select {
		case <-ctx.Done():
			return exitRequested, nil

			// DoneCh will return if renewal fails, or if the remaining lease
			// duration is under a built-in threshold and either renewing is not
			// extending it or renewing is disabled.  In both cases, the caller
			// should attempt a re-read of the secret. Clients should check the
			// return value of the channel to see if renewal was successful.
		case err := <-watcher.DoneCh():
			// Leases created by a token get revoked when the token is revoked.
			return expiring, err

			// RenewCh is a channel that receives a message when a successful
			// renewal takes place and includes metadata about the renewal.
		case info := <-watcher.RenewCh():
			slog.Info("renewal successful", slog.String("renewed_at", info.RenewedAt.String()),
				slog.String("secret", name), slog.String("lease_duration", fmt.Sprintf("%ds", info.Secret.LeaseDuration)))
		}
	}
}

func handleWatcherResult(result renewResult, onExpire ...func()) error {
	switch {
	case result&exitRequested != 0:
		slog.Debug("result is exitRequested", slog.Int("result", int(result)))
		return nil
	case result&expiring != 0:
		slog.Debug("result is expiring", slog.Int("result", int(result)))
		if len(onExpire) == 0 {
			return fmt.Errorf("no onExpire functions provided")
		}
		for _, f := range onExpire {
			f()
		}
		return nil
	default:
		slog.Debug("no action required", slog.Int("result", int(result)))
		return nil
	}
}
