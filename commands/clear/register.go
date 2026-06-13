package clear

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:      "clear",
		Level:     registry.TierApprentice,
		Stateless: true,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd()
		},
	})
}
