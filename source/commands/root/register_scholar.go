//go:build scholar

package root

// Tier 3 (scholar). Includes file/move and file/del, which attach to the adept
// "file" group — the per-subcommand gating case.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/card"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/vault"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/entry/rotate"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/del"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/move"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/topt"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/vault/stats"
)
