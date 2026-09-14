package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/johan253/idme/internal/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all signing keys in the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		rows, err := db.New(deps.pool).ListSigningKeys(ctx)
		if err != nil {
			return fmt.Errorf("list signing keys: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KID\tACTIVE\tCREATED")
		for _, r := range rows {
			active := ""
			if r.IsActive {
				active = "*"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n",
				r.Kid,
				active,
				r.CreatedAt.Time.Format("2026-01-02 15:04:05"),
			)
		}
		return w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
