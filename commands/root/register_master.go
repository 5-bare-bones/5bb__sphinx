//go:build master

package root

// Tier 5 (master) — vault administration: rewrites, dumps, re-keys.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/config"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/vault/export"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/vault/import"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/vault/restore"
)
