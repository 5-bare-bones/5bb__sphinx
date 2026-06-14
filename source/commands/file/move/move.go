package move

import (
	"fmt"
	"path/filepath"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Move a file
sphinx file move oldFile newFile

* Move a directory
sphinx file move oldDir/ newDir/

* Move a file into a directory
sphinx file move oldDir/test.txt newDir/`

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	return &cobra.Command{
		Use:     "move <src> <dst>",
		Aliases: []string{"mv"},
		Short:   "Move a file or directory",
		Long: `Move a file or directory.

In case any of the paths contains spaces within it, it must be enclosed by double quotes.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.Errorf("accepts 2 arg(s), received %d", len(args))
			}

			oldName := command_helper.NormalizeName(args[0], true)
			if err := command_helper.Exists(vault, oldName, command_helper.File); err == nil {
				return errors.Errorf("there's no file nor directory named %q", strings.TrimSuffix(oldName, "/"))
			}
			return nil
		},
		Example: example,
		RunE:    runMove(vault),
	}
}

func runMove(vault *bolt.DB) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		oldName := args[0]
		newName := args[1]
		if oldName == "" || newName == "" {
			return errors.New("invalid format, use: sphinx file move <oldName> <newName>")
		}

		oldName = command_helper.NormalizeName(oldName, true)
		newName = command_helper.NormalizeName(newName, true)
		oldNameIsDir := strings.HasSuffix(oldName, "/")
		newNameIsDir := strings.HasSuffix(newName, "/")

		if oldNameIsDir {
			if !newNameIsDir {
				return errors.New("cannot move a directory into a file")
			}
			return mvDir(vault, oldName, newName)
		}

		// Move file into directory
		if !oldNameIsDir && newNameIsDir {
			newName += filepath.Base(oldName)
		}

		if filepath.Ext(newName) == "" {
			newName += filepath.Ext(oldName)
		}

		if err := command_helper.Exists(vault, newName, command_helper.File); err != nil {
			return err
		}

		if err := file.Rename(vault, oldName, newName); err != nil {
			return err
		}

		fmt.Printf("\n%q moved to %q\n", oldName, newName)
		return nil
	}
}

func mvDir(vault *bolt.DB, oldName, newName string) error {
	names, err := file.ListNames(vault)
	if err != nil {
		return err
	}

	fmt.Printf("Moving %q directory into %q...\n", strings.TrimSuffix(oldName, "/"), strings.TrimSuffix(newName, "/"))

	for _, name := range names {
		if strings.HasPrefix(name, oldName) {
			if err := file.Rename(vault, name, newName+strings.TrimPrefix(name, oldName)); err != nil {
				return err
			}
		}
	}

	return nil
}
