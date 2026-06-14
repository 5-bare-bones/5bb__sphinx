package importt

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	"github.com/5-bare-bones/5bb__sphinx/vault/totp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Import
sphinx vault import keepass --path path/to/file

* Import and delete the file:
sphinx vault import 1password --erase --path path/to/file`

type importOptions struct {
	path  string
	erase bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := importOptions{}
	cmd := &cobra.Command{
		Use:   "import <manager-name>",
		Short: "Import entries",
		Long: `Import entries from other password managers. Format: CSV.

If an entry already exists it will be overwritten.

Delete the CSV used with the erase flag, the file will be deleted only if no errors were encountered.

Supported:
	• 1Password
	• Bitwarden
   	• Keepass/X/XC
	• Lastpass`,
		Example: example,
		Args:    command_helper.SupportedManagers(),
		RunE:    runImport(vault, &opts),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = importOptions{}
		},
	}

	f := cmd.Flags()
	f.StringVarP(&opts.path, "path", "p", "", "source file path")
	f.BoolVarP(&opts.erase, "erase", "e", false, "erase the file on exit (only if there are no errors)")

	return cmd
}

func runImport(vault *bolt.DB, opts *importOptions) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		manager := strings.Join(args, " ")
		manager = strings.ToLower(manager)

		if opts.path == "" {
			return command_helper.ErrInvalidPath
		}
		ext := filepath.Ext(opts.path)
		if ext == "" || ext == "." {
			opts.path += ".csv"
		}

		records, err := readCSV(opts.path)
		if err != nil {
			return err
		}

		if err := createEntries(vault, manager, records); err != nil {
			return err
		}

		if opts.erase {
			if err := command_helper.Erase(opts.path); err != nil {
				return err
			}
			fmt.Println("Erased file at", opts.path)
		}

		fmt.Println("Successfully imported the entries from", manager)
		return nil
	}
}

func createEntries(vault *bolt.DB, manager string, records [][]string) error {
	// [1:] used to skip headers
	records = records[:][1:]
	entries := make([]*protobuf.Entry, len(records))

	switch manager {
	case "keepass", "keepassx":
		for i, record := range records {
			entries[i] = &protobuf.Entry{
				Name:     command_helper.NormalizeName(record[0]),
				Username: record[1],
				Password: record[2],
				URL:      record[3],
				Notes:    record[4],
				Expires:  "Never",
			}
		}

	case "keepassxc":
		for i, record := range records {
			entries[i] = &protobuf.Entry{
				// Join folder and name
				Name:     command_helper.NormalizeName(record[0] + "/" + record[1]),
				Username: record[2],
				Password: record[3],
				URL:      record[4],
				Notes:    record[5],
				Expires:  "Never",
			}
		}

	case "1password":
		for i, record := range records {
			entries[i] = &protobuf.Entry{
				Name:     command_helper.NormalizeName(record[0]),
				Username: record[2],
				Password: record[3],
				URL:      record[1],
				Notes:    fmt.Sprintf("%s.\nMember number: %s.\nRecovery Codes: %s", record[4], record[5], record[6]),
				Expires:  "Never",
			}
		}

	case "lastpass":
		for i, record := range records {
			entries[i] = &protobuf.Entry{
				// Join folder and name
				Name:     command_helper.NormalizeName(record[5] + "/" + record[4]),
				Username: record[1],
				Password: record[2],
				URL:      record[0],
				Notes:    record[3],
				Expires:  "Never",
			}
		}

	case "bitwarden":
		for i, record := range records {
			// Join folder and name
			name := command_helper.NormalizeName(record[0] + "/" + record[3])
			entries[i] = &protobuf.Entry{
				Name:     name,
				Username: record[7],
				Password: record[8],
				URL:      record[6],
				Notes:    record[4],
				Expires:  "Never",
			}

			// Create TOTP if the entry has one
			if err := createTOTP(vault, name, record[9]); err != nil {
				return err
			}
		}
	}

	return entry.Create(vault, entries...)
}

func createTOTP(vault *bolt.DB, name, rawToken string) error {
	if rawToken == "" {
		return nil
	}

	t := &protobuf.TOTP{
		Name: name,
		Raw:  rawToken,
		// Bitwarden uses 6 digits by default
		Digits: 6,
	}

	return totp.Create(vault, t)
}

func readCSV(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "opening file")
	}
	defer f.Close()

	fInfo, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "obtaining file information")
	}

	if fInfo.Size() == 0 {
		return nil, errors.New("the CSV file is empty")
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "reading csv data")
	}

	return records, nil
}
