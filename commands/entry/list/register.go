package list

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:   "list",
		Parent: "entry",
		Level:  registry.TierApprentice,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.Vault)
		},
	})
}
