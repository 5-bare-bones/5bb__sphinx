package phrase

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/GGP1/atoll"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Add an entry generating a random passphrase
sphinx add phrase Sample --length 6 --separator $ --include atoll --exclude admin,login --list nolist`

type phraseOptions struct {
	list, separator string
	incl, excl      []string
	length          uint64
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB, r io.Reader) *cobra.Command {
	opts := phraseOptions{}
	cmd := &cobra.Command{
		Use:     "phrase <name>",
		Short:   "Create an entry using a passphrase",
		Aliases: []string{"passphrase"},
		Example: example,
		Args:    command_helper.MustNotExist(vault, command_helper.Entry),
		RunE:    runPhrase(vault, r, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = phraseOptions{
				separator: " ",
			}
		},
	}

	f := cmd.Flags()
	f.Uint64VarP(&opts.length, "length", "l", 0, "number of words")
	f.StringVarP(&opts.separator, "separator", "s", " ", "character that separates each word")
	f.StringSliceVarP(&opts.incl, "include", "i", nil, "words to include in the passphrase")
	f.StringSliceVarP(&opts.excl, "exclude", "e", nil, "words to exclude from the passphrase")
	f.StringVarP(&opts.list, "list", "L", "WordList", "passphrase list used {NoList|WordList|SyllableList}")

	return cmd
}

func runPhrase(vault *bolt.DB, r io.Reader, opts *phraseOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		if opts.length < 1 {
			return command_helper.ErrInvalidLength
		}

		e, err := entryInput(r, name)
		if err != nil {
			return err
		}

		e.Password, err = genPassphrase(opts)
		if err != nil {
			return err
		}

		if err := entry.Create(vault, e); err != nil {
			return err
		}

		fmt.Printf("\n%q added\n", name)
		return nil
	}
}

func entryInput(r io.Reader, name string) (*protobuf.Entry, error) {
	reader := bufio.NewReader(r)

	username := terminal.ScanOneLine(reader, "Username")
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
		URL:      url,
		Expires:  exp,
		Notes:    notes,
	}
	return entry, nil
}

// genPassphrase returns a customized random passphrase.
func genPassphrase(opts *phraseOptions) (string, error) {
	l := atoll.WordList

	if opts.list != "" {
		opts.list = strings.ReplaceAll(opts.list, " ", "")

		switch strings.ToLower(opts.list) {
		case "nolist", "no":
			l = atoll.NoList

		case "wordlist", "word":
			// Do nothing as it's the default

		case "syllablelist", "syllable":
			l = atoll.SyllableList

		default:
			return "", errors.Errorf("invalid list: %q", opts.list)
		}
	}

	p := &atoll.Passphrase{
		Length:    opts.length,
		Separator: opts.separator,
		Include:   opts.incl,
		Exclude:   opts.excl,
		List:      l,
	}

	passphrase, err := atoll.NewSecret(p)
	if err != nil {
		return "", err
	}

	return string(passphrase), nil
}
