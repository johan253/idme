package cmd

import (
	"fmt"

	"github.com/johan253/idme/internal/keys"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use: "generate",
	Short: "Generate a new signing key pair and insert into the database",
	Long: `Generate a new Ed25519 signing keypair, encrypt the private key with
the configured JWK_KEK, and insert it as an inactive row. Promote it with
'keyctl promote <kid>' once every pod has picked it up.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		rotator := keys.NewRotator(deps.pool, deps.cipher)
		kid, err := rotator.Generate(ctx)
		if err != nil {
			return fmt.Errorf("generate signing key: %w", err)
		}
		fmt.Println(kid)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
