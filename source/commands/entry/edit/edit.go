package edit

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/sig"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Edit using the standard input
sphinx entry edit Sample

* Edit using the text editor
sphinx entry edit Sample --it`

type editOptions struct {
	interactive bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := editOptions{}
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit an entry",
		Long: `Edit an entry.

If the name is edited, sphinx will remove the entry with the old name and create one with the new name.`,
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.Entry),
		RunE:    runEdit(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = editOptions{}
		},
	}

	cmd.Flags().BoolVarP(&opts.interactive, "it", "i", false, "use the text editor")

	return cmd
}

func runEdit(vault *bolt.DB, opts *editOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		oldEntry, err := entry.Get(vault, name)
		if err != nil {
			return err
		}

		// oldEntry.Expires is already an ISO date (or "Never"), so it is shown
		// to the user as-is.

		if opts.interactive {
			return useTextEditor(vault, oldEntry)
		}

		return useStdin(vault, os.Stdin, oldEntry)
	}
}

func createTempFile(e *protobuf.Entry) (string, error) {
	f, err := os.CreateTemp("", "*.json")
	if err != nil {
		return "", errors.Wrap(err, "creating temporary file")
	}

	content, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", errors.Wrap(err, "encoding entry")
	}

	if _, err := f.Write(content); err != nil {
		return "", errors.Wrap(err, "writing temporary file")
	}

	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "closing temporary file")
	}

	return f.Name(), nil
}

// readTmpFile reads the modified file and formats the card.
func readTmpFile(filename string) (*protobuf.Entry, error) {
	var e protobuf.Entry

	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&e); err != nil {
		return nil, errors.Wrap(err, "decoding file")
	}

	return &e, nil
}

// updateEntry takes the name of the entry that's being edited to check if the name was
// changed. If it was, it will remove the old one.
func updateEntry(vault *bolt.DB, name string, e *protobuf.Entry) error {
	if e.Name == "" {
		return command_helper.ErrInvalidName
	}

	// Verify that the "expires" field has a valid format
	expires, err := command_helper.FormatExpires(e.Expires)
	if err != nil {
		return err
	}

	name = command_helper.NormalizeName(name)
	e.Name = command_helper.NormalizeName(e.Name)
	e.Expires = expires

	if err := entry.Update(vault, name, e); err != nil {
		return err
	}

	fmt.Printf("\n%q added\n", e.Name)
	return nil
}

func useStdin(vault *bolt.DB, r io.Reader, oldEntry *protobuf.Entry) error {
	fmt.Println("Type '-' to clear the field (except Name and Password) or leave blank to use the current value")
	reader := bufio.NewReader(r)

	scanln := func(field, value string) string {
		input := terminal.ScanOneLine(reader, fmt.Sprintf("%s [%s]", field, value))
		if input == "-" {
			return ""
		} else if input != "" {
			return input
		}
		return value
	}

	newEntry := &protobuf.Entry{}
	newEntry.Name = scanln("Name", oldEntry.Name)
	newEntry.Username = scanln("Username", oldEntry.Username)

	enclave, err := terminal.ScanPassword("Password", true)
	if err != nil {
		if err == terminal.ErrInvalidPassword {
			// Assume the user typed an empty string to not modify the password
			newEntry.Password = oldEntry.Password
		} else {
			return err
		}
	} else {
		pwd, err := enclave.Open()
		if err != nil {
			return errors.Wrap(err, "opening enclave")
		}

		newEntry.Password = pwd.String()
	}

	newEntry.URL = scanln("URL", oldEntry.URL)
	newEntry.Expires = scanln("Expires", oldEntry.Expires)

	notes := terminal.ScanMultipleLines(reader, fmt.Sprintf("Notes [%s]", oldEntry.Notes))
	if notes == "" {
		notes = oldEntry.Notes
	} else if notes == "-" {
		notes = ""
	}
	newEntry.Notes = notes

	return updateEntry(vault, oldEntry.Name, newEntry)
}

func useTextEditor(vault *bolt.DB, oldEntry *protobuf.Entry) error {
	editor := command_helper.SelectEditor()
	bin, err := exec.LookPath(editor)
	if err != nil {
		return errors.Errorf("executable %q not found", editor)
	}

	filename, err := createTempFile(oldEntry)
	if err != nil {
		return err
	}

	sig.Signal.AddCleanup(func() error { return command_helper.Erase(filename) })
	defer command_helper.Erase(filename)

	// Open the temporary file with the selected text editor
	edit := exec.Command(bin, filename)
	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout

	if err := edit.Start(); err != nil {
		return errors.Wrapf(err, "running %s", editor)
	}

	done := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go command_helper.WatchFile(filename, done, errCh)

	// Block until an event is received or an error occurs
	select {
	case <-done:
	case err := <-errCh:
		return err
	}

	if err := edit.Wait(); err != nil {
		return err
	}

	// Read the file and update the entry
	newEntry, err := readTmpFile(filename)
	if err != nil {
		return err
	}

	rmTabs := func(old string) string {
		return strings.ReplaceAll(old, "\t", "")
	}
	newEntry.Name = rmTabs(newEntry.Name)
	newEntry.Username = rmTabs(newEntry.Username)
	newEntry.URL = rmTabs(newEntry.URL)
	newEntry.Notes = rmTabs(newEntry.Notes)

	return updateEntry(vault, oldEntry.Name, newEntry)
}
