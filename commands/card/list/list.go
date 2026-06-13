package list

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/orderedmap"
	"github.com/5-bare-bones/5bb__sphinx/pb"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/tree"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* List one, show sensitive information and QR code
sphinx card list Sample -s -q

* Filter by name
sphinx card  Sample -f

* List all
sphinx card list`

type listOptions struct {
	filter, qr, show bool
}

// NewCmd returns a new command.
func NewCmd(db *bolt.DB) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:     "list <name>",
		Short:   "List cards",
		Aliases: []string{"ls"},
		Example: example,
		Args:    cmdutil.MustExistList(db, cmdutil.Card),
		RunE:    runList(db, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = listOptions{}
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.filter, "filter", "f", false, "filter by name")
	f.BoolVarP(&opts.qr, "qr", "q", false, "display the QR code of the number on the terminal")
	f.BoolVarP(&opts.show, "show", "s", false, "show card number and security code")

	return cmd
}

func runList(db *bolt.DB, opts *listOptions) cmdutil.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = cmdutil.NormalizeName(name)

		// List all
		if name == "" {
			cards, err := card.ListNames(db)
			if err != nil {
				return err
			}

			tree.Print(cards)
			return nil
		}

		// Filter by name
		if opts.filter {
			cards, err := card.ListNames(db)
			if err != nil {
				return err
			}

			var matches []string
			for _, card := range cards {
				matched, err := regexp.MatchString(name, card)
				if err != nil {
					return err
				}

				if matched {
					matches = append(matches, card)
				}
			}

			if len(matches) == 0 {
				return errors.New("no cards were found")
			}

			tree.Print(matches)
			return nil
		}

		// List one
		c, err := card.Get(db, name)
		if err != nil {
			return err
		}

		if opts.qr {
			return terminal.DisplayQRCode(c.Number)
		}

		printCard(name, c, opts.show)
		return nil
	}
}

func printCard(name string, c *pb.Card, show bool) {
	if !show {
		c.Number = "••••••••••••••••"
		c.SecurityCode = "•••"
	}

	mp := orderedmap.New()
	mp.Set("Type", c.Type)
	mp.Set("Number", c.Number)
	mp.Set("Security code", c.SecurityCode)
	mp.Set("Expire date", c.ExpireDate)
	mp.Set("Notes", c.Notes)

	fmt.Println(cmdutil.BuildBox(name, mp))
}
