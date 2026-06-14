package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/5-bare-bones/5bb__sphinx/auth"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/terminal"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

type record struct {
	key   []byte
	value []byte
}

const confMessage = "This script modifies all records names by performing a XOR operation " +
	"against the authentication key. Are you sure you want to proceed?"

func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("couldn't initialize the configuration: %v", err)
	}

	vaultPath := filepath.Clean(config.GetString("vault.path"))
	vault, err := bolt.Open(vaultPath, 0o600, &bolt.Options{Timeout: 200 * time.Millisecond})
	if err != nil {
		log.Fatalf("couldn't open the vault: %v", err)
	}

	if err := auth.Login(vault); err != nil {
		log.Fatalf("couldn't log in: %v", err)
	}

	if err := xorNames(vault, os.Stdin); err != nil {
		log.Fatalf("couldn't xor names: %v", err)
	}
}

func xorNames(vault *bolt.DB, r io.Reader) error {
	if !terminal.Confirm(r, confMessage) {
		return nil
	}

	tx, err := vault.Begin(true)
	if err != nil {
		return errors.Wrap(err, "starting transaction")
	}
	defer tx.Rollback()

	buckets := bucket.GetNames()
	for _, bucket := range buckets {
		b := tx.Bucket(bucket)
		cursor := b.Cursor()
		mp := make(map[string]record, b.Stats().KeyN)

		// The bucket mustn't be modified inside the loop; this will result in undefined behavior
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			mp[string(k)] = record{
				key:   vault_helper.XorName(k),
				value: v,
			}
		}

		for oldName, newRecord := range mp {
			if err := b.Put(newRecord.key, newRecord.value); err != nil {
				return errors.Wrap(err, "saving new record")
			}

			if err := b.Delete([]byte(oldName)); err != nil {
				return errors.Wrap(err, "deleting old record")
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "committing transaction")
	}

	return nil
}
