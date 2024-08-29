package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/Jacobbrewer1/f1-data/pkg/logging"
	"github.com/Jacobbrewer1/f1-data/pkg/vault"
	_ "github.com/go-sql-driver/mysql"
	vault2 "github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	EnvDbConnStr = "DB_CONN_STR"
	EnvRedisHost = "REDIS_HOST"
)

type VaultDB struct {
	Client         vault.Client
	Vip            *viper.Viper
	Enabled        bool
	CurrentSecrets *vault.Secrets
}

// ConnectDB connects to the database
func ConnectDB(ctx context.Context, v *VaultDB) (*Database, error) {
	if !v.Vip.IsSet("vault") {
		return nil, errors.New("no vault configuration found")
	}

	v.Vip.Set("database.connection_string", generateConnectionStr(v.Vip, v.CurrentSecrets))
	sqlxDb, err := createConnection(v.Vip)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db := NewDatabase(sqlxDb)

	go func() {
		err := vault.RenewLease(ctx, v.Client, v.Vip.GetString("vault.database.path"), v.CurrentSecrets.Secret, func() (*vault2.Secret, error) {
			slog.Warn("Vault lease expired, reconnecting to database")

			vs, err := v.Client.GetSecret(ctx, v.Vip.GetString("vault.database.path"))
			if err != nil {
				return nil, fmt.Errorf("error getting secrets from vault: %w", err)
			}

			dbConnectionString := generateConnectionStr(v.Vip, vs)
			v.Vip.Set("database.connection_string", dbConnectionString)

			newDb, err := createConnection(v.Vip)
			if err != nil {
				return nil, fmt.Errorf("error connecting to database: %w", err)
			}

			if err := db.Reconnect(ctx, newDb); err != nil {
				return nil, fmt.Errorf("error reconnecting to database: %w", err)
			}

			slog.Info("Database reconnected")

			return vs.Secret, nil
		})
		if err != nil {
			slog.Error("Error renewing vault lease", slog.String(logging.KeyError, err.Error()))
			os.Exit(1) // Forces new credentials to be fetched
		}
	}()

	slog.Info("Database connection established with vault")
	return db, nil
}

func createConnection(v *viper.Viper) (*sqlx.DB, error) {
	connectionString := v.GetString("database.connection_string")
	if connectionString == "" {
		return nil, errors.New("no database connection string provided")
	}

	db, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Test the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	slog.Info("Connected to database")

	return db, nil
}

func generateConnectionStr(v *viper.Viper, vs *vault.Secrets) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=90s&multiStatements=true&parseTime=true",
		vs.Data["username"],
		vs.Data["password"],
		v.GetString("database.host"),
		v.GetString("database.schema"),
	)
}
