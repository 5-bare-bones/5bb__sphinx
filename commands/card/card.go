package card

import (
	"os"

	cadd "github.com/5-bare-bones/5bb__sphinx/commands/card/add"
	ccopy "github.com/5-bare-bones/5bb__sphinx/commands/card/copy"
	crm "github.com/5-bare-bones/5bb__sphinx/commands/card/del"
	cedit "github.com/5-bare-bones/5bb__sphinx/commands/card/edit"
	cls "github.com/5-bare-bones/5bb__sphinx/commands/card/list"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
sphinx card (add|copy|edit|list|del)`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "card",
		Short:   "Card operations",
		Example: example,
	}

	cmd.AddCommand(
		cadd.NewCmd(vault, os.Stdin),
		ccopy.NewCmd(vault),
		cedit.NewCmd(vault),
		cls.NewCmd(vault),
		crm.NewCmd(vault, os.Stdin),
	)

	return cmd
}
