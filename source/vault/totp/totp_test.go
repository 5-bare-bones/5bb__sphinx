package totp

import (
	"crypto/rand"
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

func TestTOTP(t *testing.T) {
	vault := setContext(t)

	totp := &protobuf.TOTP{
		Name:   "test",
		Raw:    "IFGEWRKSIFJUMR2R",
		Digits: 6,
	}

	// Create destroys the buffer, hence we cannot use their fields anymore
	t.Run("Create", create(vault, totp))
	t.Run("Get", get(vault, totp))
	t.Run("List", list(vault, totp))
	t.Run("List names", listNames(vault, totp))
	t.Run("Remove", remove(vault, totp.Name))
}

func create(vault *bolt.DB, totp *protobuf.TOTP) func(*testing.T) {
	return func(t *testing.T) {
		err := Create(vault, totp)
		assert.NoError(t, err)
	}
}

func get(vault *bolt.DB, expected *protobuf.TOTP) func(*testing.T) {
	return func(t *testing.T) {
		got, err := Get(vault, expected.Name)
		assert.NoError(t, err)

		equal := proto.Equal(expected, got)
		assert.True(t, equal)
	}
}

func list(vault *bolt.DB, expected *protobuf.TOTP) func(*testing.T) {
	return func(t *testing.T) {
		totps, err := List(vault)
		assert.NoError(t, err)

		assert.NotZero(t, len(totps), "Expected one or more totps")

		got := totps[0]
		equal := proto.Equal(expected, got)
		assert.True(t, equal)
	}
}

func listNames(vault *bolt.DB, expected *protobuf.TOTP) func(*testing.T) {
	return func(t *testing.T) {
		totps, err := ListNames(vault)
		assert.NoError(t, err)

		assert.NotZero(t, len(totps), "Expected one or more totps")

		got := totps[0]
		assert.Equal(t, expected.Name, got)
	}
}

func remove(vault *bolt.DB, name string) func(*testing.T) {
	return func(t *testing.T) {
		err := Remove(vault, name)
		assert.NoError(t, err)
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
			name: string('\x00'),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			err := Create(vault, &protobuf.TOTP{Name: tc.name})
			assert.Error(t, err)
		})
	}
}

func TestGetError(t *testing.T) {
	vault := setContext(t)

	_, err := Get(vault, "non-existent")
	assert.Error(t, err)
}

func TestCryptErrors(t *testing.T) {
	vault := setContext(t)

	// Create the one used by Get and List
	name := "test"
	err := Create(vault, &protobuf.TOTP{Name: name})
	assert.NoError(t, err)

	// Try to get the TOTP with another password
	config.Set("auth.password", memguard.NewEnclave([]byte("invalid")))

	_, err = Get(vault, name)
	assert.Error(t, err)
	_, err = List(vault)
	assert.Error(t, err)
}

func TestProtoErrors(t *testing.T) {
	vault := setContext(t)

	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.TOTP.GetName())
		buf := make([]byte, 64)
		rand.Read(buf)
		encBuf, _ := crypt.Encrypt(buf)
		return b.Put([]byte("unformatted"), encBuf)
	})
	assert.NoError(t, err, "Failed writing invalid type")

	_, err = Get(vault, "unformatted")
	assert.Error(t, err)

	_, err = List(vault)
	assert.Error(t, err)
}

func TestKeyError(t *testing.T) {
	vault := setContext(t)

	err := Create(vault, &protobuf.TOTP{Name: ""})
	assert.Error(t, err)
}

func setContext(t testing.TB) *bolt.DB {
	return vault_helper.SetContext(t, bucket.TOTP.GetName())
}
