// Package debug is a hero/dev-tier introspection tool: it shows the raw bbolt
// bucket layout and decrypt-dumps the contents of your OWN vault. It is gated
// behind the dev build tag and is never part of a student build.
package debug

import (
	"fmt"
	"io"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/db/entry"
	"github.com/5-bare-bones/5bb__sphinx/db/file"
	"github.com/5-bare-bones/5bb__sphinx/db/totp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// NewCmd returns the debug command.
func NewCmd(db *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:   "debug",
		Short: "Inspect raw buckets and decrypt-dump your own vault (dev only)",
		Example: `
sphinx debug`,
		RunE: runDebug(db),
	}
}

func runDebug(db *bolt.DB) cmdutil.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		if err := dumpBuckets(out, db); err != nil {
			return err
		}
		return dumpRecords(out, db)
	}
}

// dumpBuckets prints each bucket and its key count, including kure_auth.
func dumpBuckets(out io.Writer, db *bolt.DB) error {
	fmt.Fprintln(out, "buckets:")
	return db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Fprintf(out, "  %-12s %d keys\n", string(name), b.Stats().KeyN)
			return nil
		})
	})
}

// dumpRecords decrypts and prints every record. This deliberately exposes
// secrets in plaintext — that is the whole point of a hero-tier debug dump.
func dumpRecords(out io.Writer, db *bolt.DB) error {
	entries, err := entry.List(db)
	if err != nil {
		return errors.Wrap(err, "listing entries")
	}
	fmt.Fprintf(out, "\nentries (%d):\n", len(entries))
	for _, e := range entries {
		fmt.Fprintf(out, "  %s  user=%q pass=%q url=%q\n", e.Name, e.Username, e.Password, e.URL)
	}

	cards, err := card.List(db)
	if err != nil {
		return errors.Wrap(err, "listing cards")
	}
	fmt.Fprintf(out, "\ncards (%d):\n", len(cards))
	for _, c := range cards {
		fmt.Fprintf(out, "  %s  number=%q cvc=%q exp=%q\n", c.Name, c.Number, c.SecurityCode, c.ExpireDate)
	}

	totps, err := totp.List(db)
	if err != nil {
		return errors.Wrap(err, "listing totp")
	}
	fmt.Fprintf(out, "\ntotp (%d):\n", len(totps))
	for _, t := range totps {
		fmt.Fprintf(out, "  %s  raw=%q digits=%d\n", t.Name, t.Raw, t.Digits)
	}

	names, err := file.ListNames(db)
	if err != nil {
		return errors.Wrap(err, "listing files")
	}
	fmt.Fprintf(out, "\nfiles (%d):\n", len(names))
	for _, n := range names {
		fmt.Fprintf(out, "  %s\n", n)
	}

	return nil
}
