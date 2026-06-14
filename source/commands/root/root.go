package root

import (
	"fmt"
	"os"
	"runtime/debug"
	"slices"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/internal/buildinfo"
	"github.com/5-bare-bones/5bb__sphinx/internal/registry"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// builtinStateless are the cobra/help built-ins that never touch the vault.
// Real commands declare statelessness via registry.Entry.Stateless instead, so
// the set stays correct for whatever tier was compiled in.
var builtinStateless = map[string]struct{}{
	"help":      {},
	"-v":        {},
	"--version": {},
}

type rootOptions struct {
	version bool
}

// NewCmd assembles the root command from whatever registered itself into the
// registry in this build. Which commands that is depends on the tier-tagged
// import aggregators (see register_*.go); lower tiers literally compile in
// fewer command packages.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := rootOptions{}
	cmd := &cobra.Command{
		Use:           "sphinx",
		Short:         "sphinx ~ CLI password manager with sessions",
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		RunE: runRoot(&opts),
	}
	applyHelpTheme(cmd)

	cmd.Flags().BoolVarP(&opts.version, "version", "v", false, "display sphinx version")

	ctx := registry.BuildContext{Vault: vault, In: os.Stdin, Out: os.Stdout}
	entries := registry.Entries()

	// First pass: build every top-level command and index it by verb so that
	// subcommands can attach to their parent group.
	groups := make(map[string]*cobra.Command, len(entries))
	for _, e := range entries {
		if e.Parent != "" {
			continue
		}
		c := e.New(ctx)
		applyIdentity(c, e)
		groups[e.Verb] = c
		cmd.AddCommand(c)
	}

	// Second pass: attach subcommands to their parent group. A child whose
	// parent was not compiled into this build is skipped defensively.
	for _, e := range entries {
		if e.Parent == "" {
			continue
		}
		if parent, ok := groups[e.Parent]; ok {
			child := e.New(ctx)
			applyIdentity(child, e)
			parent.AddCommand(child)
		}
	}

	// Hide the auto-generated `help` command so it doesn't duplicate the
	// `-h/--help` flag in the command listing, mirroring how the `completion`
	// command is hidden above. The flag still provides per-command help.
	cmd.InitDefaultHelpCmd()
	if hc, _, err := cmd.Find([]string{"help"}); err == nil {
		hc.Hidden = true
	}

	return cmd
}

// applyIdentity makes the registry verb the command's primary name and adds its
// aliases. If the constructor already customized the command (it sets its own
// Use and aliases, e.g. the rank-renamed `it`), the name is left untouched.
func applyIdentity(c *cobra.Command, e registry.Entry) {
	if e.Verb != "" && len(c.Aliases) == 0 {
		if _, rest, found := strings.Cut(c.Use, " "); found {
			c.Use = e.Verb + " " + rest
		} else {
			c.Use = e.Verb
		}
	}
	for _, a := range e.Aliases {
		if a == c.Name() || slices.Contains(c.Aliases, a) {
			continue
		}
		c.Aliases = append(c.Aliases, a)
	}
}

func runRoot(opts *rootOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		if opts.version {
			printVersion()
			return nil
		}

		_ = cmd.Usage()
		return nil
	}
}

func printVersion() {
	var rev string
	if bi, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range bi.Settings {
			if setting.Key == "vcs.revision" {
				rev = setting.Value
				break
			}
		}
	}

	fmt.Printf("sphinx %s [tier: %s] %s\n", buildinfo.Version, buildinfo.Tier, rev)
}

// IsStatelessCommand reports whether the named command can run without opening
// the vault. It consults the registry so the answer reflects the compiled-in
// tier, plus the cobra built-ins.
func IsStatelessCommand(command string) bool {
	if _, ok := builtinStateless[command]; ok {
		return true
	}
	for _, e := range registry.Entries() {
		if e.Parent == "" && e.Verb == command && e.Stateless {
			return true
		}
	}
	return false
}
