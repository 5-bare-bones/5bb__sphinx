package copy

import (
	"strings"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Copy password and clean after 15m
sphinx entry copy Sample --timeout 15m

* Copy username
sphinx entry copy Sample --username

* Copy both username and password consecutively
sphinx entry copy Sample --all`

type copyOptions struct {
	timeout  time.Duration
	username bool
	all      bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := copyOptions{}
	cmd := &cobra.Command{
		Use:     "copy <name>",
		Short:   "Copy entry credentials to the clipboard",
		Aliases: []string{"cp"},
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.Entry),
		RunE:    runCopy(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = copyOptions{}
		},
	}

	f := cmd.Flags()
	f.DurationVarP(&opts.timeout, "timeout", "t", 0, "clipboard clearing timeout")
	f.BoolVarP(&opts.username, "username", "u", false, "copy entry username")
	f.BoolVarP(&opts.all, "all", "a", false, "copy entry username and password consecutively")

	return cmd
}

func runCopy(vault *bolt.DB, opts *copyOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		e, err := entry.Get(vault, name)
		if err != nil {
			return err
		}

		if opts.all {
			if err := command_helper.WriteClipboard(cmd, opts.timeout, "Username", e.Username); err != nil {
				return err
			}

			return command_helper.WriteClipboard(cmd, opts.timeout, "Password", e.Password)
		}

		field := "Password"
		value := e.Password
		if opts.username {
			field = "Username"
			value = e.Username
		}

		return command_helper.WriteClipboard(cmd, opts.timeout, field, value)
	}
}
