package vault

import (
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
sphinx vault (backup|restore|import|export|stats)`

// NewCmd returns the bare "vault" group command.
//
// Vault administration verbs (stats, backup, import, export, restore) register
// themselves into the registry with Parent "vault" and are attached by
// commands/root, each gated to its own tier.
func NewCmd(vault *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:     "vault",
		Short:   "Vault administration",
		Example: example,
	}
}
