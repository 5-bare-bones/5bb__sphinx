package vault_helper_test

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

var (
	record = &protobuf.Card{
		Name:         "test",
		Number:       "12313121",
		SecurityCode: "007",
	}
	bucketName = vault_helper.GetBucketName(record)
)

func TestGet(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return vault_helper.Put(b, record)
	})
	assert.NoError(t, err)

	got := &protobuf.Card{}
	err = vault_helper.Get(vault, record.Name, got)
	assert.NoError(t, err)

	equal := proto.Equal(record, got)
	assert.True(t, equal)
}

func TestGetBucketName(t *testing.T) {
	cases := []struct {
		desc     string
		record   vault_helper.Record
		expected []byte
	}{
		{
			desc:     "Entry",
			record:   &protobuf.Entry{},
			expected: bucket.Entry.GetName(),
		},
		{
			desc:     "Card",
			record:   &protobuf.Card{},
			expected: bucket.Card.GetName(),
		},
		{
			desc:     "File",
			record:   &protobuf.File{},
			expected: bucket.File.GetName(),
		},
		{
			desc:     "TOTP",
			record:   &protobuf.TOTP{},
			expected: bucket.TOTP.GetName(),
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			got := vault_helper.GetBucketName(tc.record)
			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestList(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

	createRecord(t, vault, record)
	record2 := &protobuf.Card{
		Name: "west",
	}
	createRecord(t, vault, record2)
	expected := []*protobuf.Card{record, record2}

	got, err := vault_helper.List(vault, &protobuf.Card{})
	assert.NoError(t, err)

	for _, e := range expected {
		for _, g := range got {
			if e.Name == g.Name {
				equal := proto.Equal(e, g)
				assert.True(t, equal)
			}
		}
	}
}

func TestListNames(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

	recordA := "a"
	recordB := "b"
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if err := b.Put(vault_helper.XorName([]byte(recordA)), nil); err != nil {
			return err
		}
		return b.Put(vault_helper.XorName([]byte(recordB)), nil)
	})
	assert.NoError(t, err)

	got, err := vault_helper.ListNames(vault, bucketName)
	assert.NoError(t, err)

	// We expect the names xored with the auth key. They should be ordered
	expected := []string{recordA, recordB}
	assert.Equal(t, expected, got)
}

func TestListNamesNil(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

	err := vault.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
	assert.NoError(t, err, "Failed deleting the file bucket")

	list, err := vault_helper.ListNames(vault, bucketName)
	assert.NoError(t, err)
	assert.Nil(t, list)
}

func TestPut(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)
	createRecord(t, vault, record)
}

func TestRemove(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

	recordA := "a"
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return b.Put(vault_helper.XorName([]byte(recordA)), nil)
	})
	assert.NoError(t, err)

	err = vault_helper.Remove(vault, bucketName, recordA)
	assert.NoError(t, err)

	expected := make([]string, 0)
	got, err := vault_helper.ListNames(vault, bucketName)
	assert.NoError(t, err)

	assert.Equal(t, expected, got)
}

func TestRemoveNone(t *testing.T) {
	err := vault_helper.Remove(nil, nil)
	assert.NoError(t, err)
}

func TestCryptErrors(t *testing.T) {
	vault := vault_helper.SetContext(t, bucket.Entry.GetName())

	name := "test decrypt error"

	e := &protobuf.Entry{Name: name, Expires: "Never"}
	createRecord(t, vault, e)

	// Try to get the entry with other password
	config.Set("auth.password", memguard.NewEnclave([]byte("invalid")))

	err := vault_helper.Get(vault, name, &protobuf.Entry{})
	assert.Error(t, err)
	// For some reason it does not fail if a card struct is used
	_, err = vault_helper.List(vault, &protobuf.Entry{})
	assert.Error(t, err)
}

func TestGetErrors(t *testing.T) {
	vault := vault_helper.SetContext(t, bucket.Entry.GetName())

	err := vault_helper.Get(vault, "non-existent", &protobuf.Entry{})
	assert.Error(t, err)
}

func TestKeyError(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)
	vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		return vault_helper.Put(b, &protobuf.Card{Name: ""})
	})
}

func TestProtoErrors(t *testing.T) {
	vault := vault_helper.SetContext(t, bucket.Entry.GetName())

	name := "unformatted"
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.Entry.GetName())
		buf := make([]byte, 32)
		encBuf, _ := crypt.Encrypt(buf)
		return b.Put([]byte(name), encBuf)
	})
	assert.NoError(t, err, "Failed writing invalid type")

	err = vault_helper.Get(vault, name, &protobuf.Entry{})
	assert.Error(t, err)
	_, err = vault_helper.List(vault, &protobuf.Entry{})
	assert.Error(t, err)
}

func TestPutErrors(t *testing.T) {
	vault := vault_helper.SetContext(t, bucketName)

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

	vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for _, tc := range cases {
			t.Run(tc.desc, func(t *testing.T) {
				err := vault_helper.Put(b, &protobuf.Card{Name: tc.name})
				assert.Error(t, err)
			})
		}
		return nil
	})
}

func TestXorName(t *testing.T) {
	defer config.Reset()

	key := []byte{
		51, 0, 107, 95, 158, 240, 55, 129, 1, 249, 4,
		159, 37, 118, 174, 228, 69, 140, 141, 199, 105,
		124, 4, 120, 253, 220, 202, 0, 199, 47, 164, 134,
	}
	config.Set("auth.key", key)

	cases := []struct {
		name     string
		expected string
	}{
		{
			name: "test",
		},
		{
			name: "adidas",
		},
		{
			name: "github",
		},
		{
			name: "nike",
		},
		{
			name: "folder/test",
		},
		{
			name: "super_extralarge_name",
		},
		{
			name: "123456789",
		},
		{
			name: string([]byte{51, 0, 107, 95, 158, 240}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			xorName := vault_helper.XorName([]byte(tc.name))
			assert.NotEqual(t, tc.name, xorName)

			gotName := vault_helper.XorName(xorName)
			assert.Equal(t, tc.name, string(gotName))
		})
	}
}

func createRecord(t *testing.T, vault *bolt.DB, record vault_helper.Record) {
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(vault_helper.GetBucketName(record))
		return vault_helper.Put(b, record)
	})
	assert.NoError(t, err)
}
