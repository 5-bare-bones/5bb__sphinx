package file

import (
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
sphinx file (add|del|show|edit|list|move|touch)`

// NewCmd returns the bare "file" group command.
//
// Its subcommands no longer wire themselves up here; each registers itself into
// the registry with Parent "file" and is attached by commands/root. This lets
// higher-tier subcommands (move, del — scholar) be gated independently of the
// lower-tier ones (add, show, list, touch — adept).
func NewCmd(vault *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:     "file",
		Short:   "File operations",
		Example: example,
	}
}
