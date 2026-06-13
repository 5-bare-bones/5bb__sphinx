//go:build master

package root

// Tier 5 (master) — vault administration: rewrites, dumps, re-keys.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/config"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/export"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/import"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/restore"
)
