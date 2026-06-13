package entry

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

func TestEntry(t *testing.T) {
	vault := setContext(t)

	e := &protobuf.Entry{
		Name:     "test",
		Username: "testing",
		URL:      "golang.org",
		Expires:  "Never",
		Notes:    "",
	}
	e2 := &protobuf.Entry{Name: "test2"}
	names := map[string]struct{}{
		e.Name:  {},
		e2.Name: {},
	}

	t.Run("Create", create(vault, e, e2))
	t.Run("Get", get(vault, e))
	t.Run("List", list(vault, e, e2))
	t.Run("List names", listNames(vault, names))
	t.Run("Remove", remove(vault, e.Name, e2.Name))
	t.Run("Update", update(vault))
	t.Run("Update name", updateName(vault))
}

func create(vault *bolt.DB, entries ...*protobuf.Entry) func(*testing.T) {
	return func(t *testing.T) {
		err := Create(vault, entries...)
		assert.NoError(t, err)
	}
}

func get(vault *bolt.DB, expected *protobuf.Entry) func(*testing.T) {
	return func(t *testing.T) {
		got, err := Get(vault, expected.Name)
		assert.NoError(t, err)

		if !proto.Equal(expected, got) {
			t.Errorf("Expected %v, got %v", expected, got)
		}
	}
}

func list(vault *bolt.DB, expected ...*protobuf.Entry) func(*testing.T) {
	return func(t *testing.T) {
		entries, err := List(vault)
		assert.NoError(t, err)

		for _, got := range entries {
			for _, exp := range expected {
				if got.Name == exp.Name {
					equal := proto.Equal(exp, got)
					assert.True(t, equal)
				}
			}
		}
	}
}

func listNames(vault *bolt.DB, names map[string]struct{}) func(*testing.T) {
	return func(t *testing.T) {
		entryNames, err := ListNames(vault)
		assert.NoError(t, err)

		for _, name := range entryNames {
			_, ok := names[name]
			assert.Truef(t, ok, "Expected %q to be in the list but it isn't", name)
		}
	}
}

func remove(vault *bolt.DB, names ...string) func(*testing.T) {
	return func(t *testing.T) {
		err := Remove(vault, names...)
		assert.NoError(t, err)
	}
}

func update(vault *bolt.DB) func(*testing.T) {
	return func(t *testing.T) {
		oldEntry := &protobuf.Entry{Name: "test"}
		err := Create(vault, oldEntry)
		assert.NoError(t, err)

		newEntry := &protobuf.Entry{Name: "test", Username: "username"}
		err = Update(vault, oldEntry.Name, newEntry)
		assert.NoError(t, err)

		_, err = Get(vault, newEntry.Name)
		assert.NoError(t, err)
	}
}

func updateName(vault *bolt.DB) func(*testing.T) {
	return func(t *testing.T) {
		oldEntry := &protobuf.Entry{Name: "old"}
		err := Create(vault, oldEntry)
		assert.NoError(t, err)

		newEntry := &protobuf.Entry{Name: "new"}
		err = Update(vault, oldEntry.Name, newEntry)
		assert.NoError(t, err)

		_, err = Get(vault, newEntry.Name)
		assert.NoError(t, err)

		_, err = Get(vault, oldEntry.Name)
		assert.Error(t, err)
	}
}

func TestCreateNone(t *testing.T) {
	vault := setContext(t)
	err := Create(vault)
	assert.NoError(t, err)

	names, err := ListNames(vault)
	assert.NoError(t, err)

	assert.Zero(t, len(names), "Expected no entries")
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
			name: string('\x00'),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := Create(vault, &protobuf.Entry{Name: tc.name})
			assert.Error(t, err)
		})
	}
}

func TestGetErrors(t *testing.T) {
	vault := setContext(t)

	_, err := Get(vault, "non-existent")
	assert.Error(t, err)
}

func TestUpdateError(t *testing.T) {
	vault := setContext(t)

	name := string([]rune{'\x00'})
	err := Update(vault, "old", &protobuf.Entry{Name: name})
	assert.Error(t, err)
}

func TestCryptErrors(t *testing.T) {
	vault := setContext(t)

	name := "test decrypt error"

	e := &protobuf.Entry{Name: name, Expires: "Never"}
	err := Create(vault, e)
	assert.NoError(t, err)

	// Try to get the entry with other password
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
		b := tx.Bucket(bucket.Entry.GetName())
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

	err := Create(vault, &protobuf.Entry{Name: ""})
	assert.Error(t, err)
}

func setContext(t testing.TB) *bolt.DB {
	return vault_helper.SetContext(t, bucket.Entry.GetName())
}
