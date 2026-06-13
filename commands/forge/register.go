package forge

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:  "forge",
		Level: registry.TierTutor,
		// forge mints fresh vault files; it never touches the default vault, so
		// it runs without a login.
		Stateless: true,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd()
		},
	})
}
