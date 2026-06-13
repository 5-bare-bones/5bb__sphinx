package export

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:  "export",
		Level: registry.TierMaster,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.DB)
		},
	})
}
