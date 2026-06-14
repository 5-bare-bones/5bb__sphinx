package touch

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:   "touch",
		Parent: "file",
		Level:  registry.TierAdept,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.Vault)
		},
	})
}
