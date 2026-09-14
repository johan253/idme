package cmd

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/johan253/idme/internal/db"
	"github.com/spf13/cobra"
)

var (
	pruneOlderThan time.Duration
	pruneDryRun    bool
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Short: "Delete inactive signing keys older than a cutoff",
	Long: `Remove inactive signing keys from the database that are older than
the given duration. The cutoff should exceed your JWT TTL plus a grace period,
so no in-flight token could still reference a deleted key.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		cutoff := pgtype.Timestamptz{
			Time:  time.Now().Add(-pruneOlderThan),
			Valid: true,
		}
		q := db.New(deps.pool)
		rows, err := q.ListPrunableSigningKeys(ctx, cutoff)
		if err != nil {
			return fmt.Errorf("list prunable keys: %w", err)
		}
		if len(rows) == 0 {
			fmt.Fprintln(cmd.ErrOrStderr(), "nothing to prune")
			return nil
		}
		for _, r := range rows {
			age := time.Since(r.CreatedAt.Time).Round(time.Second)
			fmt.Printf("%s\t%s ago\n", r.Kid, age)
		}
		if pruneDryRun {
			fmt.Fprintf(cmd.ErrOrStderr(), "dry-run: %d key(s) would be deleted\n", len(rows))
			return nil
		}
		n, err := q.PruneSigningKeys(ctx, cutoff)
		if err != nil {
			return fmt.Errorf("prune signing keys: %w", err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "deleted %d key(s)\n", n)
		return nil
	},
}

func init() {
	pruneCmd.Flags().DurationVar(&pruneOlderThan, "older-than", 24*time.Hour,
		"delete inactive keys older than this duration")
	pruneCmd.Flags().BoolVar(&pruneDryRun, "dry-run", false,
		"print what would be deleted, don't actually delete")
	rootCmd.AddCommand(pruneCmd)
}
