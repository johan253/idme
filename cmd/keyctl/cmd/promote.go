
package cmd

import (
	"fmt"

	"github.com/johan253/idme/internal/keys"
	"github.com/spf13/cobra"
)

var promoteCmd = &cobra.Command{
	Use: "promote <kid>",
	Short: "Promote a signing key to active",
	Long: `Mark the given key as the sole active signing key. The database's
partial unique index guarantees at most one key is active at a time - 
concurrent promotes will fail rather than corrupt state.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		kid := args[0]
		rotator := keys.NewRotator(deps.pool, deps.cipher)
		if err := rotator.Promote(ctx, kid); err != nil {
			return fmt.Errorf("promote %s: %w", kid, err)
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "promoted %s\n", kid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(promoteCmd)
}
