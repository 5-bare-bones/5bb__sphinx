package add

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Add a new card
sphinx card add Sample`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:     "add <name>",
		Short:   "Add a card",
		Aliases: []string{"create", "new"},
		Example: example,
		Args:    command_helper.MustNotExist(vault, command_helper.Card),
		RunE:    runAdd(vault, r),
	}
}

func runAdd(vault *bolt.DB, r io.Reader) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		c, err := input(vault, name, r)
		if err != nil {
			return err
		}

		if err := card.Create(vault, c); err != nil {
			return err
		}

		fmt.Printf("\n%q added\n", name)
		return nil
	}
}

func input(vault *bolt.DB, name string, r io.Reader) (*protobuf.Card, error) {
	reader := bufio.NewReader(r)
	c := &protobuf.Card{
		Name:         name,
		Type:         terminal.ScanOneLine(reader, "Type"),
		Number:       terminal.ScanOneLine(reader, "Number"),
		SecurityCode: terminal.ScanOneLine(reader, "Security code"),
		ExpireDate:   terminal.ScanOneLine(reader, "Expire date"),
		Notes:        terminal.ScanMultipleLines(reader, "Notes"),
	}

	return c, nil
}
