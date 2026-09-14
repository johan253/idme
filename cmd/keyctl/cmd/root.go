// Package cmd contains all commands for the `keyctl` cli
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/johan253/idme/internal/config"
	"github.com/johan253/idme/internal/keys"
	"github.com/spf13/cobra"
)

var (
	deps *cliDeps
)

type cliDeps struct {
	cfg *config.Config
	pool *pgxpool.Pool
	cipher keys.Cipher
}

var rootCmd = &cobra.Command{
	Use:   "keyctl",
	Short: "Manage signing keys for the idme auth server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, pool, cipher, err := bootstrap(cmd.Context())
		deps = &cliDeps{
			cfg: cfg,
			pool: pool,
			cipher: cipher,
		}
		return err
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if deps.pool != nil { deps.pool.Close() }
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func bootstrap(ctx context.Context) (*config.Config, *pgxpool.Pool, keys.Cipher, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load config: %w", err)
	}
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connect db: %w", err)
	}
	cipher, err := keys.NewAESGCMCipher(cfg.Kek)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("build cipher: %w", err)
	}

	return cfg, pool, cipher, nil
}
