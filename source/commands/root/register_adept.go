//go:build adept

package root

// Tier 2 (adept). Cumulative: an adept build passes -tags 'adept'; a scholar
// build passes -tags 'adept scholar', so this file is included there too.
import (
	_ "github.com/5-bare-bones/5bb__sphinx/commands/entry"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/entry/copy"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/entry/del"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/entry/edit"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/add"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/edit"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/list"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/show"
	_ "github.com/5-bare-bones/5bb__sphinx/commands/file/touch"
)
