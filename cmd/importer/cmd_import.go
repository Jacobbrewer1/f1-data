package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/Jacobbrewer1/f1-data/pkg/repositories"
	"github.com/Jacobbrewer1/f1-data/pkg/repositories/importer"
	importerSvc "github.com/Jacobbrewer1/f1-data/pkg/services/importer"
	"github.com/Jacobbrewer1/f1-data/pkg/vault"
	"github.com/google/subcommands"
	"github.com/spf13/viper"
)

type importCmd struct {
	// configLocation is the location of the config file
	configLocation string
}

func (i *importCmd) Name() string {
	return "import"
}

func (i *importCmd) Synopsis() string {
	return "Import data from a file into the database"
}

func (i *importCmd) Usage() string {
	return `import:
  Import data from a file into the database.
`
}

func (i *importCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&i.configLocation, "config", "config.json", "The location of the config file")
}

func (i *importCmd) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	v := viper.New()
	v.SetConfigFile(i.configLocation)
	db, err := i.setup(ctx, v)
	if err != nil {
		slog.Error("Error setting up database", slog.String("error", err.Error()))
		return subcommands.ExitFailure
	}

	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("Error closing database", slog.String("error", err.Error()))
		}
	}()

	slog.Info("Database connection established")

	repo := importer.NewRepository(db)
	svc := importerSvc.NewService(repo, v.GetString("importer.f1_base_url"))

	toYear := v.GetInt("importer.to_year")
	// If toYear is -1, set it to the current year
	if toYear == -1 {
		toYear = time.Now().Year()
	}

	err = svc.Import(v.GetInt("importer.from_year"), toYear)
	if err != nil {
		slog.Error("Error importing data", slog.String("error", err.Error()))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (i *importCmd) setup(ctx context.Context, v *viper.Viper) (db *repositories.Database, err error) {
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if !v.IsSet("vault") {
		return nil, errors.New("vault configuration not found")
	}

	vaultDb := &repositories.VaultDB{
		Client:         nil,
		Vip:            v,
		Enabled:        false,
		CurrentSecrets: nil,
	}

	slog.Info("Vault configuration found, attempting to connect")
	vaultDb.Enabled = true

	vc, err := vault.NewClientUserPass(v)
	if err != nil {
		return nil, fmt.Errorf("error creating vault client: %w", err)
	}

	vaultDb.Client = vc

	slog.Debug("Vault client created")

	vs, err := vc.GetSecret(ctx, v.GetString("vault.database.path"))
	if errors.Is(err, vault.ErrSecretNotFound) {
		return nil, fmt.Errorf("secrets not found in vault: %s", v.GetString("vault.database.path"))
	} else if err != nil {
		return nil, fmt.Errorf("error getting secrets from vault: %w", err)
	}

	slog.Debug("Vault secrets retrieved")
	vaultDb.CurrentSecrets = vs
	db, err = repositories.ConnectDB(ctx, vaultDb)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	slog.Info("Database connection generate from vault secrets")

	return db, nil
}
