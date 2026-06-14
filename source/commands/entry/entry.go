package entry

import (
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
sphinx entry (add|copy|edit|list|del|rotate)`

// NewCmd returns the bare "entry" group command.
//
// Its subcommands register themselves into the registry with Parent "entry" and
// are attached by commands/root, so each verb (add, copy, edit, list, del,
// rotate) can be gated to its own tier independently of the group.
func NewCmd(vault *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:     "entry",
		Short:   "Entry operations",
		Example: example,
	}
}
