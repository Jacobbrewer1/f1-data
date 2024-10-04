package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	vaulty2 "github.com/Jacobbrewer1/vaulty"
	_ "github.com/go-sql-driver/mysql"
	vault2 "github.com/hashicorp/vault/api"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type DatabaseConnector interface {
	ConnectDB() (*Database, error)
}

type databaseConnector struct {
	ctx            context.Context
	client         vaulty2.Client
	vip            *viper.Viper
	currentSecrets *vault2.Secret
}

func NewDatabaseConnector(opts ...ConnectionOption) (DatabaseConnector, error) {
	c := new(databaseConnector)

	for _, opt := range opts {
		opt(c)
	}

	if c.ctx == nil {
		c.ctx = context.Background()
	}

	if c.client == nil {
		return nil, errors.New("no vault client provided")
	} else if c.vip == nil {
		return nil, errors.New("no viper configuration provided")
	} else if c.currentSecrets == nil {
		return nil, errors.New("no current secrets provided")
	}

	return c, nil
}

// ConnectDB connects to the database
func (d *databaseConnector) ConnectDB() (*Database, error) {
	if !d.vip.IsSet("vault") {
		return nil, errors.New("no vault configuration found")
	}

	d.vip.Set("database.connection_string", generateConnectionStr(d.vip, d.currentSecrets))
	sqlxDb, err := createConnection(d.vip)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db := NewDatabase(sqlxDb)

	go func() {
		err := vaulty2.RenewLease(d.ctx, d.client, d.vip.GetString("vault.database.path"), d.currentSecrets, func() (*vault2.Secret, error) {
			slog.Warn("Vault lease expired, reconnecting to database")

			vs, err := d.client.GetSecret(d.ctx, d.vip.GetString("vault.database.path"))
			if err != nil {
				return nil, fmt.Errorf("error getting secrets from vault: %w", err)
			}

			dbConnectionString := generateConnectionStr(d.vip, vs)
			d.vip.Set("database.connection_string", dbConnectionString)

			newDb, err := createConnection(d.vip)
			if err != nil {
				return nil, fmt.Errorf("error connecting to database: %w", err)
			}

			if err := db.Reconnect(d.ctx, newDb); err != nil {
				return nil, fmt.Errorf("error reconnecting to database: %w", err)
			}

			slog.Info("Database reconnected")

			return vs, nil
		})
		if err != nil {
			slog.Error("Error renewing vault lease", slog.String(loggingKeyError, err.Error()))
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

func generateConnectionStr(v *viper.Viper, vs *vault2.Secret) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=90s&multiStatements=true&parseTime=true",
		vs.Data["username"],
		vs.Data["password"],
		v.GetString("database.host"),
		v.GetString("database.schema"),
	)
}
