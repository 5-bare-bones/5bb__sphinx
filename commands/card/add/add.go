package add

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/pb"
	"github.com/5-bare-bones/5bb__sphinx/terminal"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Add a new card
sphinx card add Sample`

// NewCmd returns a new command.
func NewCmd(db *bolt.DB, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:     "add <name>",
		Short:   "Add a card",
		Aliases: []string{"create", "new"},
		Example: example,
		Args:    cmdutil.MustNotExist(db, cmdutil.Card),
		RunE:    runAdd(db, r),
	}
}

func runAdd(db *bolt.DB, r io.Reader) cmdutil.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = cmdutil.NormalizeName(name)

		c, err := input(db, name, r)
		if err != nil {
			return err
		}

		if err := card.Create(db, c); err != nil {
			return err
		}

		fmt.Printf("\n%q added\n", name)
		return nil
	}
}

func input(db *bolt.DB, name string, r io.Reader) (*pb.Card, error) {
	reader := bufio.NewReader(r)
	c := &pb.Card{
		Name:         name,
		Type:         terminal.ScanOneLine(reader, "Type"),
		Number:       terminal.ScanOneLine(reader, "Number"),
		SecurityCode: terminal.ScanOneLine(reader, "Security code"),
		ExpireDate:   terminal.ScanOneLine(reader, "Expire date"),
		Notes:        terminal.ScanMultipleLines(reader, "Notes"),
	}

	return c, nil
}
