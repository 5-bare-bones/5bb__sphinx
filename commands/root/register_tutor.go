//go:build tutor

package root

// Tier 7 (sphinx/tutor) — vault-minting and riddle tooling. Orthogonal to the
// student ladder; built with -tags '... tutor' on top of a master base.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/forge"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/riddle"
)
