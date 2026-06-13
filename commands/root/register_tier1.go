package root

// Tier 1 (apprentice) — the default build, no build tag. These command packages
// are compiled into every binary.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/add"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/clear"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/generate"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/it"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/list"
)
