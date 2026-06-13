package file

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

// The "file" group shell. Its subcommands register themselves with Parent
// "file"; root attaches them. add/show/edit/list/touch arrive at adept, move/del at
// scholar.
func init() {
	registry.Register(registry.Entry{
		Verb:  "file",
		Level: registry.TierAdept,
		New: func(c registry.BuildContext) *cobra.Command {
			return NewCmd(c.DB)
		},
	})
}
