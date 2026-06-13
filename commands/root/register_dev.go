//go:build dev

package root

// Tier 6 (hero) — the contributor's dev surface. Layered onto a master base:
// a dev build passes -tags 'adept scholar keeper master dev'.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/debug"
)
