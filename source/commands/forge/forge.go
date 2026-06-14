// Package forge mints per-student vaults from a YAML world-manifest. It is
// tutor-only tooling — never compiled into a student build — and doubles as a
// generator of test fixtures and reference data.
package forge

import (
	"fmt"
	"os"
	"path/filepath"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/internal/vault"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const example = `
* Mint every vault described in world.yaml into ./vaults
sphinx forge world.yaml --out ./vaults`

// manifest is the on-disk world description.
type manifest struct {
	Argon2   *argon2Spec `yaml:"argon2"`
	OutDir   string      `yaml:"output_dir"`
	Students []student   `yaml:"students"`
}

type argon2Spec struct {
	Iterations uint32 `yaml:"iterations"`
	Memory     uint32 `yaml:"memory"`
	Threads    uint32 `yaml:"threads"`
}

type student struct {
	Name       string      `yaml:"name"`
	Passphrase string      `yaml:"passphrase"`
	Entries    []entrySpec `yaml:"entries"`
	Cards      []cardSpec  `yaml:"cards"`
	TOTPs      []totpSpec  `yaml:"totps"`
}

type entrySpec struct {
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	URL      string `yaml:"url"`
	Notes    string `yaml:"notes"`
	Expires  string `yaml:"expires"`
}

type cardSpec struct {
	Name         string `yaml:"name"`
	Type         string `yaml:"type"`
	Number       string `yaml:"number"`
	SecurityCode string `yaml:"security_code"`
	ExpireDate   string `yaml:"expire_date"`
	Notes        string `yaml:"notes"`
}

type totpSpec struct {
	Name   string `yaml:"name"`
	Raw    string `yaml:"raw"`
	Digits int32  `yaml:"digits"`
}

// NewCmd returns the forge command.
func NewCmd() *cobra.Command {
	var outDir string
	cmd := &cobra.Command{
		Use:     "forge <manifest.yaml>",
		Short:   "Mint per-student vaults from a world manifest (tutor only)",
		Example: example,
		Args:    cobra.ExactArgs(1),
		RunE:    runForge(&outDir),
	}
	cmd.Flags().StringVarP(&outDir, "out", "o", "", "output directory (overrides manifest's output_dir)")
	return cmd
}

func runForge(outDir *string) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		raw, err := os.ReadFile(args[0])
		if err != nil {
			return errors.Wrap(err, "reading manifest")
		}

		var m manifest
		if err := yaml.Unmarshal(raw, &m); err != nil {
			return errors.Wrap(err, "parsing manifest")
		}
		if len(m.Students) == 0 {
			return errors.New("manifest defines no students")
		}

		dir := m.OutDir
		if *outDir != "" {
			dir = *outDir
		}
		if dir == "" {
			dir = "."
		}
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return errors.Wrap(err, "creating output directory")
		}

		params := resolveArgon2(m.Argon2)

		// Validate the whole manifest before minting anything, so a bad entry
		// never leaves a half-forged set of vaults on disk.
		seen := make(map[string]struct{}, len(m.Students))
		for _, s := range m.Students {
			if s.Name == "" {
				return errors.New("student with empty name")
			}
			if _, dup := seen[s.Name]; dup {
				return errors.Errorf("duplicate student %q", s.Name)
			}
			seen[s.Name] = struct{}{}
			if s.Passphrase == "" {
				return errors.Errorf("student %q has no passphrase", s.Name)
			}
			if _, err := os.Stat(filepath.Join(dir, s.Name+".vault")); err == nil {
				return errors.Errorf("vault already exists: %s (refusing to overwrite)",
					filepath.Join(dir, s.Name+".vault"))
			}
		}

		for _, s := range m.Students {
			path := filepath.Join(dir, s.Name+".vault")
			if err := vault.Create(path, s.Passphrase, params, seedFor(s)); err != nil {
				return errors.Wrapf(err, "forging vault for %q", s.Name)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "forged %s (%d entries, %d cards, %d totp)\n",
				path, len(s.Entries), len(s.Cards), len(s.TOTPs))
		}

		fmt.Fprintf(cmd.OutOrStdout(), "minted %d vault(s) into %s\n", len(m.Students), dir)
		return nil
	}
}

// resolveArgon2 fills any unset (zero) field from the fast fixture defaults.
func resolveArgon2(spec *argon2Spec) vault.Argon2 {
	p := vault.DefaultArgon2()
	if spec == nil {
		return p
	}
	if spec.Iterations != 0 {
		p.Iterations = spec.Iterations
	}
	if spec.Memory != 0 {
		p.Memory = spec.Memory
	}
	if spec.Threads != 0 {
		p.Threads = spec.Threads
	}
	return p
}

func seedFor(s student) vault.Seed {
	seed := vault.Seed{}
	for _, e := range s.Entries {
		seed.Entries = append(seed.Entries, &protobuf.Entry{
			Name:     e.Name,
			Username: e.Username,
			Password: e.Password,
			URL:      e.URL,
			Notes:    e.Notes,
			Expires:  e.Expires,
		})
	}
	for _, c := range s.Cards {
		seed.Cards = append(seed.Cards, &protobuf.Card{
			Name:         c.Name,
			Type:         c.Type,
			Number:       c.Number,
			SecurityCode: c.SecurityCode,
			ExpireDate:   c.ExpireDate,
			Notes:        c.Notes,
		})
	}
	for _, t := range s.TOTPs {
		digits := t.Digits
		if digits == 0 {
			digits = 6
		}
		seed.TOTPs = append(seed.TOTPs, &protobuf.TOTP{
			Name:   t.Name,
			Raw:    t.Raw,
			Digits: digits,
		})
	}
	return seed
}
