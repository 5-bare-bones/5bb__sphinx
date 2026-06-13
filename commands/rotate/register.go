package rotate

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:  "rotate",
		Level: registry.TierScholar,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.DB)
		},
	})
}
