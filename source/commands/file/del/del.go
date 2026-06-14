package del

import (
	"fmt"
	"io"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Remove a file
sphinx file del Sample

* Remove a directory
sphinx file del SampleDir/

* Remove multiple files
sphinx file del Sample Sample2 Sample3`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:     "del <names>",
		Aliases: []string{"delete", "rm"},
		Short:   "Remove files or directories",
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.File, true),
		RunE:    runDelete(vault, r),
	}
}

func runDelete(vault *bolt.DB, r io.Reader) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		if !terminal.Confirm(r, "Are you sure you want to proceed?") {
			return nil
		}

		names := make([]string, 0, len(args))
		for _, name := range args {
			name = command_helper.NormalizeName(name, true)

			if !strings.HasSuffix(name, "/") {
				names = append(names, name)
				fmt.Println("Remove:", name)
				continue
			}

			files, err := file.ListNames(vault)
			if err != nil {
				return err
			}

			for _, f := range files {
				if strings.HasPrefix(f, name) {
					names = append(names, f)
					fmt.Println("Remove:", f)
				}
			}
		}

		return file.Remove(vault, names...)
	}
}
