package it

import (
	"github.com/5-bare-bones/5bb__sphinx/internal/buildinfo"
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
)

func init() {
	registry.Register(registry.Entry{
		Verb:  "it",
		Level: registry.TierApprentice,
		New: func(c registry.BuildContext) *cobra.Command {
			cmd := NewCmd(c.DB)
			// The interactive walker is renamed after the binary's rank
			// (sphinx apprentice .. sphinx master), while "it" stays as an
			// alias so muscle memory and docs keep working.
			rank := buildinfo.CurrentTier().String()
			cmd.Use = rank + " <command|flags|name>"
			cmd.Aliases = []string{"it"}
			return cmd
		},
	})
}
