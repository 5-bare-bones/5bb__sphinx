package riddle

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:  "riddle",
		Level: registry.TierTutor,
		// Riddles are self-contained files; solving one needs no vault.
		Stateless: true,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd()
		},
	})
}
