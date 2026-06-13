package add

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/commands/add/phrase"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/GGP1/atoll"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Add an entry using a custom password
sphinx add Sample -c

* Add an entry generating a random password
sphinx add Sample --length 27 --levels 1,2,3,4,5 --include & --exclude / --repeat`

type addOptions struct {
	include, exclude string
	levels           []int
	length           uint64
	custom, repeat   bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB, r io.Reader) *cobra.Command {
	opts := addOptions{}
	cmd := &cobra.Command{
		Use:     "add <name>",
		Short:   "Add an entry",
		Aliases: []string{"create", "new"},
		Example: example,
		Args:    command_helper.MustNotExist(vault, command_helper.Entry),
		RunE:    runAdd(vault, r, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = addOptions{}
		},
	}

	cmd.AddCommand(phrase.NewCmd(vault, r))

	f := cmd.Flags()
	f.BoolVarP(&opts.custom, "custom", "c", false, "use a custom password")
	f.Uint64VarP(&opts.length, "length", "l", 0, "password length")
	f.IntSliceVarP(&opts.levels, "levels", "L", []int{1, 2, 3, 4, 5}, "password levels")
	f.StringVarP(&opts.include, "include", "i", "", "characters to include in the password")
	f.StringVarP(&opts.exclude, "exclude", "e", "", "characters to exclude from the password")
	f.BoolVarP(&opts.repeat, "repeat", "r", true, "allow character repetition")

	return cmd
}

func runAdd(vault *bolt.DB, r io.Reader, opts *addOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		if !opts.custom {
			if opts.length < 1 {
				return command_helper.ErrInvalidLength
			}
			if len(opts.levels) == 0 {
				return errors.New("please specify levels")
			}
		}

		e, err := entryInput(r, name, opts.custom)
		if err != nil {
			return err
		}

		if !opts.custom {
			// Generate random password
			e.Password, err = genPassword(opts)
			if err != nil {
				return err
			}
		}

		if err := entry.Create(vault, e); err != nil {
			return err
		}

		fmt.Printf("\n%q added\n", name)
		return nil
	}
}

// genPassword returns a customized random password or an error.
func genPassword(opts *addOptions) (string, error) {
	levels := make([]atoll.Level, len(opts.levels))
	for i, lvl := range opts.levels {
		switch lvl {
		case 1:
			levels[i] = atoll.Lower
		case 2:
			levels[i] = atoll.Upper
		case 3:
			levels[i] = atoll.Digit
		case 4:
			levels[i] = atoll.Space
		case 5:
			levels[i] = atoll.Special

		default:
			return "", errors.Errorf("invalid level [%d]", lvl)
		}
	}

	p := &atoll.Password{
		Length:  opts.length,
		Levels:  levels,
		Include: opts.include,
		Exclude: opts.exclude,
		Repeat:  opts.repeat,
	}

	password, err := atoll.NewSecret(p)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func entryInput(r io.Reader, name string, custom bool) (*protobuf.Entry, error) {
	var password string
	reader := bufio.NewReader(r)

	username := terminal.ScanOneLine(reader, "Username")
	if custom {
		enclave, err := terminal.ScanPassword("Password", true)
		if err != nil {
			return nil, err
		}

		pwd, err := enclave.Open()
		if err != nil {
			return nil, errors.Wrap(err, "opening enclave")
		}

		password = pwd.String()
	}
	url := terminal.ScanOneLine(reader, "URL")
	expires := terminal.ScanOneLine(reader, "Expires [yyyy-mm-dd]")
	notes := terminal.ScanMultipleLines(reader, "Notes")

	exp, err := command_helper.FormatExpires(expires)
	if err != nil {
		return nil, err
	}

	entry := &protobuf.Entry{
		Name:     name,
		Username: username,
		Password: password,
		URL:      url,
		Expires:  exp,
		Notes:    notes,
	}

	return entry, nil
}
