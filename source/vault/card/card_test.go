package card

import (
	"testing"

	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/crypt"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/awnumar/memguard"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func TestCard(t *testing.T) {
	vault := setContext(t)

	c := &protobuf.Card{
		Name:         "test",
		Type:         "debit",
		Number:       "4403650939814064",
		SecurityCode: "1234",
		ExpireDate:   "16/2024",
	}

	t.Run("Create", create(vault, c))
	t.Run("Get", get(vault, c))
	t.Run("List", list(vault, c))
	t.Run("List names", listNames(vault, c))
	t.Run("Remove", remove(vault, c.Name))
	t.Run("Update", update(vault))
	t.Run("Update name", updateName(vault))
}

func create(vault *bolt.DB, c *protobuf.Card) func(*testing.T) {
	return func(t *testing.T) {
		err := Create(vault, c)
		assert.NoError(t, err)
	}
}

func get(vault *bolt.DB, expected *protobuf.Card) func(*testing.T) {
	return func(t *testing.T) {
		got, err := Get(vault, expected.Name)
		assert.NoError(t, err)

		if !proto.Equal(expected, got) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func list(vault *bolt.DB, expected *protobuf.Card) func(*testing.T) {
	return func(t *testing.T) {
		cards, err := List(vault)
		assert.NoError(t, err)

		assert.NotZero(t, len(cards), "Expected one or more cards")

		got := cards[0]
		if !proto.Equal(expected, got) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func listNames(vault *bolt.DB, expected *protobuf.Card) func(*testing.T) {
	return func(t *testing.T) {
		cards, err := ListNames(vault)
		assert.NoError(t, err)

		if len(cards) == 0 {
			t.Fatal("Expected one or more cards, got 0")
		}

		got := cards[0]
		if got != expected.Name {
			t.Errorf("Expected %s, got %s", expected.Name, got)
		}
	}
}

func remove(vault *bolt.DB, name string) func(*testing.T) {
	return func(t *testing.T) {
		err := Remove(vault, name)
		assert.NoError(t, err)
	}
}

func update(vault *bolt.DB) func(*testing.T) {
	return func(t *testing.T) {
		oldCard := &protobuf.Card{Name: "test"}
		err := Create(vault, oldCard)
		assert.NoError(t, err)

		newCard := &protobuf.Card{Name: "test", Type: "debit"}
		err = Update(vault, oldCard.Name, newCard)
		assert.NoError(t, err)

		_, err = Get(vault, newCard.Name)
		assert.NoError(t, err)
	}
}

func updateName(vault *bolt.DB) func(*testing.T) {
	return func(t *testing.T) {
		oldCard := &protobuf.Card{Name: "old"}
		err := Create(vault, oldCard)
		assert.NoError(t, err)

		newCard := &protobuf.Card{Name: "new"}
		err = Update(vault, oldCard.Name, newCard)
		assert.NoError(t, err)

		_, err = Get(vault, newCard.Name)
		assert.NoError(t, err)

		_, err = Get(vault, oldCard.Name)
		assert.Error(t, err)
	}
}

func TestCreateErrors(t *testing.T) {
	vault := setContext(t)

	cases := []struct {
		desc string
		name string
	}{
		{
			desc: "Invalid name",
			name: "",
		},
		{
			desc: "Null characters",
			name: string([]rune{'\x00'}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := Create(vault, &protobuf.Card{Name: tc.name})
			assert.Error(t, err)
		})
	}
}

func TestGetError(t *testing.T) {
	vault := setContext(t)

	_, err := Get(vault, "non-existent")
	assert.Error(t, err)
}

func TestUpdateError(t *testing.T) {
	vault := setContext(t)

	name := string([]rune{'\x00'})
	err := Update(vault, "old", &protobuf.Card{Name: name})
	assert.Error(t, err)
}

func TestCryptErrors(t *testing.T) {
	vault := setContext(t)

	name := "crypt-errors"
	err := Create(vault, &protobuf.Card{Name: name})
	assert.NoError(t, err)

	// Try to get the card with another password
	config.Set("auth.password", memguard.NewEnclave([]byte("invalid")))

	_, err = Get(vault, name)
	assert.Error(t, err)
	_, err = List(vault)
	assert.Error(t, err)
}

func TestProtoErrors(t *testing.T) {
	vault := setContext(t)

	name := "unformatted"
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Card.GetName())
		buf := make([]byte, 64)
		encBuf, _ := crypt.Encrypt(buf)
		return b.Put([]byte(name), encBuf)
	})
	assert.NoError(t, err, "Failed writing invalid type")

	_, err = Get(vault, name)
	assert.Error(t, err)
	_, err = List(vault)
	assert.Error(t, err)
}

func TestKeyError(t *testing.T) {
	vault := setContext(t)

	err := Create(vault, &protobuf.Card{Name: ""})
	assert.Error(t, err)
}

func setContext(t testing.TB) *bolt.DB {
	return vault_helper.SetContext(t, bucket.Card.GetName())
}
