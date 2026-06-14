package auth

import (
	"fmt"
	"testing"

	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestParameters(t *testing.T) {
	vault := setContext(t)

	key := []byte("test")
	expected := Params{
		AuthKey: key,
		Argon2: Argon2{
			Iterations: 1,
			Memory:     1500000,
			Threads:    4,
		},
		UseKeyfile: true,
	}

	err := Register(vault, key, expected)
	assert.NoError(t, err, "Registration failed")

	got, err := GetParams(vault)
	assert.NoError(t, err)

	// Force auth key to be the same as it's encrypted and it won't match
	got.AuthKey = key
	assert.Equal(t, expected, got)
}

func TestEmptyParameters(t *testing.T) {
	vault := setContext(t)
	tx, _ := vault.Begin(true)
	tx.DeleteBucket(bucket.Auth.GetName())
	tx.Commit()

	expected := Params{}
	got, err := GetParams(vault)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestSetParametersInvalidKeys(t *testing.T) {
	vault := setContext(t)

	key := []byte("invalid")
	params := Params{
		AuthKey: key,
		Argon2: Argon2{
			Iterations: 1,
			Memory:     1,
			Threads:    1,
		},
		UseKeyfile: true,
	}

	cases := []struct {
		key  *[]byte
		desc string
	}{
		{
			desc: "iterations",
			key:  &iterKey,
		},
		{
			desc: "memory",
			key:  &memKey,
		},
		{
			desc: "threads",
			key:  &thKey,
		},
		{
			desc: "keyfile",
			key:  &keyfileKey,
		},
		{
			desc: "auth",
			key:  &authKey,
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Invalid %s key", tc.desc), func(t *testing.T) {
			*tc.key = nil

			tx, err := vault.Begin(true)
			assert.NoError(t, err, "Failed opening transaction")

			err = storeParams(tx, key, params)
			assert.Error(t, err)
			tx.Commit()

			// Fill the variable so we can test the others
			*tc.key = []byte("1")
		})
	}
}

func setContext(t testing.TB) *bolt.DB {
	return vault_helper.SetContext(t, bucket.Auth.GetName())
}
