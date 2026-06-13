package config

import (
	"fmt"
	"os"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	argon2cmd "github.com/5-bare-bones/5bb__sphinx/commands/config/argon2"
	"github.com/5-bare-bones/5bb__sphinx/commands/config/create"
	"github.com/5-bare-bones/5bb__sphinx/commands/config/edit"
	"github.com/5-bare-bones/5bb__sphinx/config"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Read configuration file
sphinx config`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Read the configuration file",
		Aliases: []string{"cfg"},
		Example: example,
		RunE:    runConfig(),
	}

	cmd.AddCommand(argon2cmd.NewCmd(vault), create.NewCmd(), edit.NewCmd(vault))

	return cmd
}

func runConfig() command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		path := config.Filename()
		data, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "reading configuration file")
		}

		content := strings.TrimSpace(string(data))
		fmt.Printf(`
File location: %s

%s
`, path, content)

		return nil
	}
}
