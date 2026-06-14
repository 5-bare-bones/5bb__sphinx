package argon2

import (
	"fmt"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/commands/config/argon2/test"
	authVault "github.com/5-bare-bones/5bb__sphinx/vault/auth"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
sphinx config argon2`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "argon2",
		Short:   "Display currently used argon2 parameters",
		Aliases: []string{"argon"},
		Example: example,
		RunE:    runArgon2(vault),
	}

	cmd.AddCommand(test.NewCmd())

	return cmd
}

func runArgon2(vault *bolt.DB) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		params, err := authVault.GetParams(vault)
		if err != nil {
			return err
		}

		fmt.Printf("Iterations: %d\nMemory: %d\nThreads: %d\n",
			params.Argon2.Iterations, params.Argon2.Memory, params.Argon2.Threads)
		return nil
	}
}
