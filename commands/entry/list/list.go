package list

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/orderedmap"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/tree"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* List one and show sensitive information
sphinx entry list Sample --show

* List one and show the password QR code
sphinx entry list Sample --qr

* Filter by name
sphinx entry list Sample --filter

* List all
sphinx entry list`

type listOptions struct {
	filter, qr, show bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:     "list <name>",
		Aliases: []string{"ls", "entries"},
		Short:   "List entries",
		Long: `List entries.

Listing all the entries does not check for expired entries, this decision was taken to prevent high loads when the number of entries is elevated. Listing a single entry does notifies if it is expired.`,
		Example: example,
		Args:    command_helper.MustExistList(vault, command_helper.Entry),
		RunE:    runList(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = listOptions{}
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.filter, "filter", "f", false, "filter by name")
	f.BoolVarP(&opts.qr, "qr", "q", false, "display the password QR code on the terminal")
	f.BoolVarP(&opts.show, "show", "s", false, "show entry password")

	return cmd
}

func runList(vault *bolt.DB, opts *listOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		// List all
		if name == "" {
			entries, err := entry.ListNames(vault)
			if err != nil {
				return err
			}

			tree.Print(entries)
			return nil
		}

		// Filter by name
		if opts.filter {
			entries, err := entry.ListNames(vault)
			if err != nil {
				return err
			}

			var matches []string
			for _, entry := range entries {
				matched, err := regexp.MatchString(name, entry)
				if err != nil {
					return err
				}

				if matched {
					matches = append(matches, entry)
				}
			}

			if len(matches) == 0 {
				return errors.New("no entries were found")
			}

			tree.Print(matches)
			return nil
		}

		// List one
		e, err := entry.Get(vault, name)
		if err != nil {
			return err
		}

		if opts.qr {
			return terminal.DisplayQRCode(e.Password)
		}

		printEntry(name, e, opts.show)
		return nil
	}
}

func printEntry(name string, e *protobuf.Entry, show bool) {
	if !show {
		e.Password = "•••••••••••••••"
	}

	if expired(e.Expires) {
		e.Expires = "EXPIRED"
	}

	mp := orderedmap.New()
	mp.Set("Username", e.Username)
	mp.Set("Password", e.Password)
	mp.Set("URL", e.URL)
	mp.Set("Expires", e.Expires)
	mp.Set("Notes", e.Notes)

	fmt.Println(command_helper.BuildBox(name, mp))
}

// expired reports whether the entry's ISO expiration date is in the past.
func expired(expires string) bool {
	if expires == "Never" {
		return false
	}

	// Stored as an ISO date (YYYY-MM-DD). An unparseable value (e.g. legacy
	// data) is treated as not expired rather than silently marking it EXPIRED.
	expiration, err := time.Parse("2006-01-02", expires)
	if err != nil {
		return false
	}
	return time.Now().After(expiration)
}
