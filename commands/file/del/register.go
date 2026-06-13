package del

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:   "del",
		Parent: "file",
		Level:  registry.TierScholar,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.Vault, c.In)
		},
	})
}
