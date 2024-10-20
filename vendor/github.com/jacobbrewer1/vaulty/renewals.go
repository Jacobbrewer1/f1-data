package vaulty

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	hashiVault "github.com/hashicorp/vault/api"
)

// renewResult is a bitmask which could contain one or more of the values below
type renewResult uint8

const (
	renewError renewResult = 1 << iota
	exitRequested
	expiring // will be revoked soon
)

// RenewLease Once you've set the token for your Vault client, you will need to
// periodically renew it. Likewise, the database credentials lease will expire
// at some point and also needs to be renewed periodically.
//
// A function like this one should be run as a goroutine to avoid blocking.
// Production applications may also need to be more tolerant of failures and
// retry on errors rather than exiting.
//
// Additionally, enterprise Vault users should be aware that due to eventual
// consistency, the API may return unexpected errors when running Vault with
// performance standbys or performance replication, despite the client having
// a freshly renewed token. See the link below for several ways to mitigate
// this which are outside the scope of this code sample.
//
// ref: https://www.vaultproject.io/docs/enterprise/consistency#vault-1-7-mitigations
func RenewLease(ctx context.Context, client ClientHandler, name string, credentials *hashiVault.Secret, renewFunc RenewalFunc) error {
	slog.Debug("renewing lease", slog.String("secret", name))

	currentCreds := credentials

	for {
		res, err := leaseRenew(ctx, client, name, currentCreds)
		if err != nil {
			return fmt.Errorf("unable to renew lease: %w", err)
		} else if res&exitRequested != 0 {
			// Context was cancelled. Program is exiting.
			slog.Debug("exit requested", slog.String("secret", name))
			return nil
		}

		err = handleWatcherResult(res, func() {
			newCreds, err := renewFunc()
			if err != nil {
				slog.Error("unable to renew credentials", slog.String(loggingKeyError, err.Error()))
				os.Exit(1) // Forces new credentials to be fetched
			}

			currentCreds = newCreds
		})
		if err != nil {
			return fmt.Errorf("unable to handle watcher result: %w", err)
		}

		slog.Info("lease renewed", slog.String("secret", name))
	}
}

func leaseRenew(ctx context.Context, client ClientHandler, name string, credentials *hashiVault.Secret) (renewResult, error) {
	credentialsWatcher, err := client.Client().NewLifetimeWatcher(&hashiVault.LifetimeWatcherInput{
		Secret:    credentials,
		Increment: 3600,
	})
	if err != nil {
		return renewError, fmt.Errorf("unable to initialize credentials lifetime watcher: %w", err)
	}

	go credentialsWatcher.Start()
	defer credentialsWatcher.Stop()

	res, err := monitorWatcher(ctx, name, credentialsWatcher)
	if err != nil {
		return renewError, fmt.Errorf("unable to monitor watcher: %w", err)
	}

	return res, nil
}
