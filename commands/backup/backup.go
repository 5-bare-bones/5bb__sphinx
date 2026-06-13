package backup

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/sig"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

const example = `
* Create a file backup
sphinx backup --path path/to/file

* Serve the vault on a local server, port 7777
sphinx backup --http --port 7777

* Download vault
curl localhost:7777 > vault_name`

type backupOptions struct {
	path  string
	port  uint16
	httpB bool
}

// NewCmd returns a new command.
func NewCmd(vault *bolt.DB) *cobra.Command {
	opts := backupOptions{}
	cmd := &cobra.Command{
		Use:     "backup",
		Short:   "Create vault backup",
		Example: example,
		RunE:    opts.runBackup(vault),
		PostRun: func(cmd *cobra.Command, args []string) {
			// Reset variables (session)
			opts = backupOptions{
				port: 8080,
			}
		},
	}

	f := cmd.Flags()
	f.BoolVar(&opts.httpB, "http", false, "serve vault file on a local server")
	f.StringVar(&opts.path, "path", "", "destination file path")
	f.Uint16Var(&opts.port, "port", 8080, "server port")

	cmd.MarkFlagsMutuallyExclusive("http", "path")

	return cmd
}

func (opts *backupOptions) runBackup(vault *bolt.DB) command_helper.RunErrorFunction {
	return func(cmd *cobra.Command, args []string) error {
		if opts.httpB {
			return serveFile(vault, opts.port)
		}

		return fileBackup(vault, opts.path)
	}
}

// serveFile serves the file on localhost.
func serveFile(vault *bolt.DB, port uint16) error {
	if port == 0 {
		return errors.New("invalid port")
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	sig.Signal.AddCleanup(func() error {
		// Do not exit after a signal as we are handling the shutdown
		sig.Signal.KeepAlive()
		fmt.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			return errors.Wrap(err, "graceful shutdown")
		}

		if err := server.Close(); err != nil {
			return errors.Wrap(err, "closing server")
		}
		return nil
	})

	// Register route only once, otherwise it will panic if
	// called multiple times inside a session
	var once sync.Once
	once.Do(func() {
		http.HandleFunc("/", httpBackup(vault))
	})
	fmt.Printf("Serving vault on http://localhost:%d (Press Ctrl+C to quit)\n", port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return errors.Wrap(err, "starting server")
	}

	return nil
}

// fileBackup writes the vault to a new file.
func fileBackup(vault *bolt.DB, path string) error {
	if path == "" {
		return command_helper.ErrInvalidPath
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return errors.Wrap(err, "making directory")
	}

	if err := os.Chdir(dir); err != nil {
		return errors.Wrap(err, "changing working directory")
	}

	newDB, err := bolt.Open(filepath.Base(path), 0o600, nil)
	if err != nil {
		return errors.Wrap(err, "opening vault backup")
	}

	if err := bolt.Compact(newDB, vault, 0); err != nil {
		return errors.Wrap(err, "copying vault backup")
	}

	if err := newDB.Close(); err != nil {
		return errors.Wrap(err, "closing vault backup")
	}

	abs, _ := filepath.Abs(path)
	fmt.Println("Backup created at", abs)
	return nil
}

// httpBackup writes a consistent view of the vault to a http endpoint.
func httpBackup(vault *bolt.DB) http.HandlerFunc {
	name := filepath.Base(config.GetString("vault.path"))
	disposition := fmt.Sprintf(`attachment; filename=%q`, name)

	return func(w http.ResponseWriter, r *http.Request) {
		err := vault.View(func(tx *bolt.Tx) error {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", disposition)
			w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
			if _, err := tx.WriteTo(w); err != nil {
				return errors.Wrap(err, "writing the vault")
			}

			return nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// writeTo writes the entire vault to a writer.
func writeTo(vault *bolt.DB, w io.Writer) error {
	return vault.View(func(tx *bolt.Tx) error {
		if _, err := tx.WriteTo(w); err != nil {
			return errors.Wrap(err, "writing the vault")
		}
		return nil
	})
}
