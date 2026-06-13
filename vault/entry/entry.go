package entry

import (
	"strings"

	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// Create new entries.
func Create(vault *bolt.DB, entries ...*protobuf.Entry) error {
	if len(entries) == 0 {
		return nil
	}

	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Entry.GetName())
		for _, entry := range entries {
			if err := vault_helper.Put(b, entry); err != nil {
				return err
			}
		}

		return nil
	})
}

// Get retrieves the entry with the specified name.
func Get(vault *bolt.DB, name string) (*protobuf.Entry, error) {
	entry := &protobuf.Entry{}
	if err := vault_helper.Get(vault, name, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// List returns a list with all the entries.
func List(vault *bolt.DB) ([]*protobuf.Entry, error) {
	return vault_helper.List(vault, &protobuf.Entry{})
}

// ListNames returns a list with all the entries names.
func ListNames(vault *bolt.DB) ([]string, error) {
	return vault_helper.ListNames(vault, bucket.Entry.GetName())
}

// Remove removes one or more entries from the vault.
func Remove(vault *bolt.DB, names ...string) error {
	return vault_helper.Remove(vault, bucket.Entry.GetName(), names...)
}

// Update updates an entry, it removes the old one if the name differs.
func Update(vault *bolt.DB, oldName string, entry *protobuf.Entry) error {
	if strings.ContainsRune(entry.Name, '\x00') {
		return errors.New("entry name contains null characters")
	}

	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Entry.GetName())
		if oldName != entry.Name {
			xorName := vault_helper.XorName([]byte(oldName))
			if err := b.Delete(xorName); err != nil {
				return errors.Wrap(err, "remove old entry")
			}
		}
		return vault_helper.Put(b, entry)
	})
}
