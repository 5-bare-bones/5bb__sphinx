package show

import (
	"bytes"
	"io"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/file"

	"github.com/atotto/clipboard"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Write one file
sphinx show Sample

* Write one file and copy content to the clipboard
sphinx show Sample -c

* Write multiple files
sphinx show sample1 sample2 sample3`

type catOptions struct {
	copy bool
}

// NewCmd returns a new command.
func NewCmd(db *bolt.DB, w io.Writer) *cobra.Command {
	opts := catOptions{}
	cmd := &cobra.Command{
		Use:     "show <name>",
		Short:   "Read file and write to standard output",
		Aliases: []string{"cat"},
		Example: example,
		Args:    cmdutil.MustExist(db, cmdutil.File),
		RunE:    runCat(db, w, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = catOptions{}
		},
	}

	cmd.Flags().BoolVarP(&opts.copy, "copy", "c", false, "copy file content to the clipboard")

	return cmd
}

func runCat(db *bolt.DB, w io.Writer, opts *catOptions) cmdutil.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		for i, name := range args {
			if name == "" {
				return errors.Errorf("name [%d] is invalid", i)
			}

			name = cmdutil.NormalizeName(name)
			f, err := file.Get(db, name)
			if err != nil {
				return err
			}

			buf := bytes.NewBuffer(f.Content)
			buf.WriteString("\n")

			if opts.copy {
				if err := clipboard.WriteAll(buf.String()); err != nil {
					return errors.Wrap(err, "writing to clipboard")
				}
			}

			if _, err := io.Copy(w, buf); err != nil {
				return errors.Wrap(err, "copying content")
			}
		}

		return nil
	}
}
