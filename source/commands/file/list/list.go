package list

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/orderedmap"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/tree"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const (
	_ = 1 << (10 * iota)
	// KB - 1024 bytes
	KB
	// MB - 1048576 bytes
	MB
	// GB - 1073741824 bytes
	GB
	// TB - 1099511627776 bytes
	TB
)

const example = `
* Filter files by name
sphinx file list Sample --filter

* List all files
sphinx file list`

type listOptions struct {
	filter bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := listOptions{}
	cmd := &cobra.Command{
		Use:     "list <name>",
		Aliases: []string{"ls"},
		Short:   "List files",
		Example: example,
		Args:    command_helper.MustExistList(vault, command_helper.File),
		RunE:    runList(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = listOptions{}
		},
	}

	cmd.Flags().BoolVarP(&opts.filter, "filter", "f", false, "filter by name")

	return cmd
}

func runList(vault *bolt.DB, opts *listOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		name := strings.Join(args, " ")
		name = command_helper.NormalizeName(name)

		// List all
		if name == "" {
			files, err := file.ListNames(vault)
			if err != nil {
				return err
			}

			tree.Print(files)
			return nil
		}

		// Filter by name
		if opts.filter {
			files, err := file.ListNames(vault)
			if err != nil {
				return err
			}

			var matches []string
			for _, file := range files {
				matched, err := regexp.MatchString(name, file)
				if err != nil {
					return err
				}

				if matched {
					matches = append(matches, file)
				}
			}

			if len(matches) == 0 {
				return errors.New("no files were found")
			}

			tree.Print(matches)
			return nil
		}

		// List one
		f, err := file.GetCheap(vault, name)
		if err != nil {
			return err
		}

		printFile(f)
		return nil
	}
}

func printFile(f *protobuf.FileCheap) {
	parts := strings.Split(f.Name, "/")
	path := strings.Join(parts[:len(parts)-1], "/")
	bytes := f.Size
	size := fmt.Sprintf("%d bytes", bytes)
	createdAt := time.Unix(f.CreatedAt, 0)
	updatedAt := time.Unix(f.UpdatedAt, 0)

	switch {
	case bytes >= TB:
		size = fmt.Sprintf("%d TB", bytes/TB)
	case bytes >= GB:
		size = fmt.Sprintf("%d GB", bytes/GB)
	case bytes >= MB:
		size = fmt.Sprintf("%d MB", bytes/MB)
	case bytes >= KB:
		size = fmt.Sprintf("%d KB", bytes/KB)
	}

	mp := orderedmap.New()
	mp.Set("Path", "/"+path)
	mp.Set("Size", size)
	mp.Set("Created at", createdAt.String())
	if !updatedAt.IsZero() {
		mp.Set("Updated at", updatedAt.String())
	}

	fmt.Println(command_helper.BuildBox(f.Name, mp))
}
