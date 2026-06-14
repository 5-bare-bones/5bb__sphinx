package card

import (
	"strings"

	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// Create a new bank card.
func Create(vault *bolt.DB, card *protobuf.Card) error {
	return vault.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Card.GetName())
		return vault_helper.Put(b, card)
	})
}

// Get retrieves the card with the specified name.
func Get(vault *bolt.DB, name string) (*protobuf.Card, error) {
	card := &protobuf.Card{}
	if err := vault_helper.Get(vault, name, card); err != nil {
		return nil, err
	}

	return card, nil
}

// List returns a list with all the cards.
func List(vault *bolt.DB) ([]*protobuf.Card, error) {
	return vault_helper.List(vault, &protobuf.Card{})
}

// ListNames returns a list with all the cards names.
func ListNames(vault *bolt.DB) ([]string, error) {
	return vault_helper.ListNames(vault, bucket.Card.GetName())
}

// Remove removes one or more cards from the vault.
func Remove(vault *bolt.DB, names ...string) error {
	return vault_helper.Remove(vault, bucket.Card.GetName(), names...)
}

// Update updates a card, it removes the old one if the name differs.
func Update(vault *bolt.DB, oldName string, card *protobuf.Card) error {
	if strings.ContainsRune(card.Name, '\x00') {
		return errors.New("card name contains null characters")
	}

	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Card.GetName())
		if oldName != card.Name {
			xorName := vault_helper.XorName([]byte(oldName))
			if err := b.Delete(xorName); err != nil {
				return errors.Wrap(err, "remove old card")
			}
		}
		return vault_helper.Put(b, card)
	})
}
