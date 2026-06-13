package edit

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/sig"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Edit a file
sphinx file edit Sample --editor micro

* Write a file's content to a temporary file and log its path
sphinx file edit Sample -l`

type editOptions struct {
	editor string
	log    bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := editOptions{}
	cmd := &cobra.Command{
		Use:   "edit <name>",
		Short: "Edit a file",
		Long: `Edit a file.

Caution: a temporary file is created with a random name, it will be erased right after the first save but it could still be read by a malicious actor.

Notes:
	- Some editors flush the changes to the disk when closed, sphinx won't notice any modifications until then.
	- Modifying the file with a different program will prevent sphinx from erasing the file as its being blocked by another process.`,
		Example: example,
		Args:    command_helper.MustExist(vault, command_helper.File),
		RunE:    runEdit(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = editOptions{}
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.editor, "editor", "e", "", "file editor command")
	f.BoolVarP(&opts.log, "log", "l", false, "log the temporary file path and wait for modifications")

	return cmd
}

func runEdit(vault *bolt.DB, opts *editOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		var bin string
		if !opts.log {
			if opts.editor == "" {
				opts.editor = command_helper.SelectEditor()
			}
			var err error
			bin, err = exec.LookPath(opts.editor)
			if err != nil {
				return errors.Errorf("executable %q not found", opts.editor)
			}
		}

		oldFile, err := file.Get(vault, name)
		if err != nil {
			return err
		}

		filename, err := createTempFile(filepath.Ext(oldFile.Name), oldFile.Content)
		if err != nil {
			return errors.Wrap(err, "creating temporary file")
		}
		sig.Signal.AddCleanup(func() error { return command_helper.Erase(filename) })
		defer command_helper.Erase(filename)

		if opts.log {
			logTempFilename(filename)
			if err := watchFile(filename); err != nil {
				return err
			}
		} else {
			if err := runEditor(bin, filename, opts.editor); err != nil {
				return err
			}
		}

		if err := update(vault, oldFile, filename); err != nil {
			return err
		}

		fmt.Printf("\n%q updated\n", name)
		return nil
	}
}

// createTempFile creates a temporary file and returns its name.
func createTempFile(ext string, content []byte) (string, error) {
	f, err := os.CreateTemp("", "*"+ext)
	if err != nil {
		return "", errors.Wrap(err, "creating file")
	}

	if _, err := f.Write(content); err != nil {
		return "", errors.Wrap(err, "writing file")
	}

	if err := f.Close(); err != nil {
		return "", errors.Wrap(err, "closing file")
	}

	return f.Name(), nil
}

func logTempFilename(filename string) {
	fmt.Printf(`Temporary file path: %s

Caution: if any process is accessing the file at the time of modification, sphinx won't be able to erase it
`, filepath.ToSlash(filename))
}

// runEditor opens the temporary file with the editor chosen and blocks until the user saves the changes.
func runEditor(bin, filename, editor string) error {
	edit := exec.Command(bin, filename)
	edit.Stdin = os.Stdin
	edit.Stdout = os.Stdout

	if err := edit.Start(); err != nil {
		return errors.Wrapf(err, "running %s", editor)
	}

	if err := watchFile(filename); err != nil {
		return err
	}

	return edit.Wait()
}

func watchFile(filename string) error {
	done := make(chan struct{}, 1)
	errCh := make(chan error, 1)
	go command_helper.WatchFile(filename, done, errCh)

	// Block until an event is received or an error occurs
	select {
	case <-done:

	case err := <-errCh:
		return err
	}

	return nil
}

// update reads the edited content and updates the file record.
func update(vault *bolt.DB, old *protobuf.File, filename string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "reading file")
	}
	content = bytes.TrimSpace(content)

	new := &protobuf.File{
		Name:      old.Name,
		Content:   content,
		Size:      int64(len(content)),
		CreatedAt: old.CreatedAt,
		UpdatedAt: time.Now().Unix(),
	}

	if err := file.Create(vault, new); err != nil {
		return errors.Wrap(err, "updating file")
	}

	return nil
}
