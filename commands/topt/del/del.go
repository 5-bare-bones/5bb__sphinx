package del

import (
	"fmt"
	"io"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/totp"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Remove a TOTP
sphinx topt del Sample

* Remove a directory
sphinx topt del SampleDir/

* Remove multiple totp
sphinx topt del Sample Sample2 Sample3`

// NewCmd returns the a new command.
func NewCmd(vault *bolt.DB, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:     "del <names>",
		Short:   "Remove two-factor authentication codes or directories",
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.TOTP, true),
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

			totps, err := totp.ListNames(vault)
			if err != nil {
				return err
			}

			for _, c := range totps {
				if strings.HasPrefix(c, name) {
					names = append(names, c)
					fmt.Println("Remove:", c)
				}
			}
		}

		return totp.Remove(vault, names...)
	}
}
