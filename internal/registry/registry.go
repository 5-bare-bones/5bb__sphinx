// Package registry is the single source of truth for which commands exist in a
// given build of sphinx.
//
// Each command package contributes itself by calling Register from an init()
// function that lives in a small, build-tagged register.go file. Because the
// build tag gates whether that file compiles in at all, a lower-tier build
// literally does not register (and therefore never assembles) the commands of a
// higher tier — the command surface IS the access control.
//
// commands/root consumes Entries() to build the cobra tree, so --help, shell
// completion, the interactive `it` walker and the stateless dispatcher all
// reflect the compiled-in tier with no additional filtering.
package registry

import (
	"io"
	"sort"
	"sync"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// Tier is a student rank. It is recorded on every Entry as metadata (used for
// help grouping, `ascend` and docs). It is NOT the access-control mechanism:
// gating happens at compile time via build tags, not by comparing tiers at
// runtime.
type Tier int

const (
	TierApprentice Tier = iota + 1 // L1, default build (no build tag)
	TierAdept                      // L2, build tag: adept
	TierScholar                    // L3, build tag: scholar
	TierKeeper                     // L4, build tag: keeper
	TierMaster                     // L5, build tag: master
	TierHero                       // L6, build tag: dev
	TierTutor                      // L7, build tag: tutor (orthogonal to the ladder)
)

// String returns the lowercase rank name.
func (t Tier) String() string {
	switch t {
	case TierApprentice:
		return "apprentice"
	case TierAdept:
		return "adept"
	case TierScholar:
		return "scholar"
	case TierKeeper:
		return "keeper"
	case TierMaster:
		return "master"
	case TierHero:
		return "hero"
	case TierTutor:
		return "tutor"
	default:
		return "unknown"
	}
}

// BuildContext carries everything any command constructor might need. It exists to
// collapse the several historical NewCmd signatures ((vault), (vault, io.Reader),
// (vault, io.Writer), ()) behind one uniform Constructor type so that every
// command registers identically.
type BuildContext struct {
	Vault *bolt.DB
	In    io.Reader
	Out   io.Writer
}

// Constructor builds a single cobra command from a BuildCtx.
type Constructor func(BuildContext) *cobra.Command

// Entry is the (noun, verb, constructor) tuple the design calls for.
type Entry struct {
	// Verb is the command's primary name, e.g. "add", "edit", "move". root makes
	// this the command's name even if the constructor's Use differs, so a
	// command can be renamed (e.g. "2fa" -> "topt") purely from its register.go.
	Verb string
	// Aliases are alternative names the command also answers to (e.g. keeping
	// "2fa" working after renaming the verb to "topt").
	Aliases []string
	// Parent is the key of the top-level group a subcommand attaches to
	// ("file", "card", "topt", "config"). Empty for top-level commands and for
	// the group shells themselves.
	Parent string
	// Noun is reserved for the future dual-grammar tree (e.g. "entry edit" vs
	// "edit entry"). Populated but unused today.
	Noun string
	// Level is the tier at which the command first appears. Metadata only.
	Level Tier
	// Stateless marks commands that must run without opening the vault. It
	// replaces the previously hard-coded statelessCommands map and stays
	// tier-correct automatically: a command that is not compiled in cannot be
	// marked stateless.
	Stateless bool
	// New constructs the cobra command.
	New Constructor
}

var (
	mu      sync.Mutex
	entries []Entry
)

// Register adds an entry to the registry. It is meant to be called from a
// command package's init(); registration order across packages is unspecified,
// so Entries returns a deterministically sorted snapshot.
func Register(e Entry) {
	mu.Lock()
	defer mu.Unlock()
	entries = append(entries, e)
}

// Entries returns a stable-sorted (Parent, Verb) snapshot of everything that
// registered itself in this build.
func Entries() []Entry {
	mu.Lock()
	defer mu.Unlock()

	out := make([]Entry, len(entries))
	copy(out, entries)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Parent != out[j].Parent {
			return out[i].Parent < out[j].Parent
		}
		return out[i].Verb < out[j].Verb
	})
	return out
}
