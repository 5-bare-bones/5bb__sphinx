package del

import (
	"fmt"
	"io"
	"strings"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/totp"
	"github.com/5-bare-bones/5bb__sphinx/terminal"

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
func NewCmd(db *bolt.DB, r io.Reader) *cobra.Command {
	return &cobra.Command{
		Use:     "del <names>",
		Short:   "Remove two-factor authentication codes or directories",
		Example: example,
		Args:    cmdutil.MustExist(db, cmdutil.TOTP, true),
		RunE:    runDelete(db, r),
	}
}

func runDelete(db *bolt.DB, r io.Reader) cmdutil.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		if !terminal.Confirm(r, "Are you sure you want to proceed?") {
			return nil
		}

		names := make([]string, 0, len(args))
		for _, name := range args {
			name = cmdutil.NormalizeName(name, true)

			if !strings.HasSuffix(name, "/") {
				names = append(names, name)
				fmt.Println("Remove:", name)
				continue
			}

			totps, err := totp.ListNames(db)
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

		return totp.Remove(db, names...)
	}
}
