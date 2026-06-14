// Package buildinfo exposes identity that is stamped into the binary at build
// time via -ldflags -X. It tells a running sphinx which tier it is, what
// version it reports, and which public key to trust when verifying an `ascend`
// update manifest.
//
// Example:
//
//	go build -ldflags "\
//	  -X github.com/5-bare-bones/5bb__sphinx/internal/buildinfo.Version=1.2.3 \
//	  -X github.com/5-bare-bones/5bb__sphinx/internal/buildinfo.Tier=scholar \
//	  -X github.com/5-bare-bones/5bb__sphinx/internal/buildinfo.PubKeyB64=..."
package buildinfo

import "github.com/5-bare-bones/5bb__sphinx/internal/registry"

// These are intentionally plain vars so the linker can overwrite them with
// -ldflags -X. The defaults describe an untagged developer build.
var (
	// Version is the released version string, e.g. "1.2.3".
	Version = "dev"
	// Tier is the lowercase rank name this binary was built for, matching the
	// cumulative build tags it was compiled with ("apprentice".."master",
	// "hero", "tutor").
	Tier = "apprentice"
	// PubKeyB64 is the base64-encoded Ed25519 public key used to verify signed
	// ascend manifests. Empty in unsigned dev builds.
	PubKeyB64 = ""
)

// CurrentTier maps the linker-stamped Tier string onto the registry.Tier enum.
// Unknown values fall back to the apprentice tier.
func CurrentTier() registry.Tier {
	switch Tier {
	case "adept":
		return registry.TierAdept
	case "scholar":
		return registry.TierScholar
	case "keeper":
		return registry.TierKeeper
	case "master":
		return registry.TierMaster
	case "hero":
		return registry.TierHero
	case "tutor":
		return registry.TierTutor
	default:
		return registry.TierApprentice
	}
}
