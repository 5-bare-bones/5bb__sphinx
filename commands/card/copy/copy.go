package copy

import (
	"strings"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Copy the number
sphinx card copy Sample

* Copy the security code
sphinx card copy Sample -c

* Copy and clean after 30s
sphinx card copy Sample -t 30s`

type copyOptions struct {
	cvc     bool
	timeout time.Duration
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := copyOptions{}
	cmd := &cobra.Command{
		Use:     "copy <name>",
		Short:   "Copy card number or security code",
		Aliases: []string{"cp"},
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.Card),
		RunE:    runCard(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = copyOptions{}
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.cvc, "cvc", "c", false, "copy card security code")
	f.DurationVarP(&opts.timeout, "timeout", "t", 0, "clipboard clearing timeout")

	return cmd
}

func runCard(vault *bolt.DB, opts *copyOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		c, err := card.Get(vault, name)
		if err != nil {
			return err
		}

		field := "Number"
		copy := c.Number
		if opts.cvc {
			field = "Security code"
			copy = c.SecurityCode
		}

		return command_helper.WriteClipboard(cmd, opts.timeout, field, copy)
	}
}
