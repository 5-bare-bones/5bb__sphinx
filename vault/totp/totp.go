package totp

import (
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	bolt "go.etcd.io/bbolt"
)

// Create a new TOTP.
func Create(vault *bolt.DB, totp *protobuf.TOTP) error {
	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.TOTP.GetName())
		return vault_helper.Put(b, totp)
	})
}

// Get retrieves the TOTP with the specified name.
func Get(vault *bolt.DB, name string) (*protobuf.TOTP, error) {
	totp := &protobuf.TOTP{}
	if err := vault_helper.Get(vault, name, totp); err != nil {
		return nil, err
	}

	return totp, nil
}

// List returns a list with all the TOTPs.
func List(vault *bolt.DB) ([]*protobuf.TOTP, error) {
	return vault_helper.List(vault, &protobuf.TOTP{})
}

// ListNames returns a slice with all the totps names.
func ListNames(vault *bolt.DB) ([]string, error) {
	return vault_helper.ListNames(vault, bucket.TOTP.GetName())
}

// Remove removes one or more totps from the vault.
func Remove(vault *bolt.DB, names ...string) error {
	return vault_helper.Remove(vault, bucket.TOTP.GetName(), names...)
}
