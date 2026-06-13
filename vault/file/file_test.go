package file

import (
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/crypt"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/awnumar/memguard"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestFile(t *testing.T) {
	vault := setContext(t)

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)

	_, err := gw.Write([]byte("content"))
	assert.NoError(t, err, "Failed compressing content")

	err = gw.Close()
	assert.NoError(t, err, "Failed closing gzip writer")

	f := &protobuf.File{
		Name:      "test",
		Content:   buf.Bytes(),
		CreatedAt: 0,
	}
	updatedName := "tested"

	t.Run("Create", create(vault, f))
	t.Run("Get", get(vault, f))
	t.Run("Rename", rename(vault, f.Name, updatedName))
	t.Run("Get cheap", getCheap(vault, updatedName))
	f.Name = updatedName
	t.Run("List", list(vault))
	t.Run("List names", listNames(vault, updatedName))
	t.Run("Remove", remove(vault, updatedName))
}

func create(vault *bolt.DB, f *protobuf.File) func(*testing.T) {
	return func(t *testing.T) {
		err := Create(vault, f)
		assert.NoError(t, err)
	}
}

func get(vault *bolt.DB, expected *protobuf.File) func(*testing.T) {
	return func(t *testing.T) {
		got, err := Get(vault, expected.Name)
		assert.NoError(t, err)

		// Using proto.Equal fails because of differing content buffers
		assert.Equal(t, expected.Name, got.Name)
	}
}

func rename(vault *bolt.DB, name, updatedName string) func(*testing.T) {
	return func(t *testing.T) {
		err := Rename(vault, name, updatedName)
		assert.NoError(t, err)

		_, err = GetCheap(vault, name)
		assert.Error(t, err)
	}
}

func getCheap(vault *bolt.DB, expectedName string) func(*testing.T) {
	return func(t *testing.T) {
		gotName, err := GetCheap(vault, expectedName)
		assert.NoError(t, err)

		assert.Equal(t, expectedName, gotName.Name)
	}
}

func list(vault *bolt.DB) func(*testing.T) {
	return func(t *testing.T) {
		files, err := List(vault)
		assert.NoError(t, err)

		assert.NotZero(t, len(files), "Expected one or more files")
	}
}

func listNames(vault *bolt.DB, expectedName string) func(*testing.T) {
	return func(t *testing.T) {
		files, err := ListNames(vault)
		assert.NoError(t, err)

		assert.NotZero(t, len(files), "Expected one or more files")

		gotName := files[0]
		assert.Equal(t, expectedName, gotName)
	}
}

func remove(vault *bolt.DB, name string) func(*testing.T) {
	return func(t *testing.T) {
		err := Remove(vault, name)
		assert.NoError(t, err)
	}
}

func TestRemoveNone(t *testing.T) {
	vault := vault_helper.SetContext(t, bucket.File.GetName())

	err := Remove(vault)
	assert.NoError(t, err)
}

func TestCreateErrors(t *testing.T) {
	vault := vault_helper.SetContext(t, bucket.File.GetName())

	err := Create(vault, &protobuf.File{})
	assert.Error(t, err)
}

func TestGetError(t *testing.T) {
	vault := setContext(t)

	_, err := Get(vault, "non-existent")
	assert.Error(t, err)
}

func TestGetCheapError(t *testing.T) {
	vault := setContext(t)

	_, err := GetCheap(vault, "non-existent")
	assert.Error(t, err)
}

func TestRenameError(t *testing.T) {
	vault := setContext(t)

	err := Rename(vault, "non-existent", "")
	assert.Error(t, err)
}

func TestCryptErrors(t *testing.T) {
	vault := setContext(t)

	name := "crypt-errors"
	err := Create(vault, &protobuf.File{Name: name})
	assert.NoError(t, err)

	// Try to get the file with other password
	config.Set("auth.password", memguard.NewEnclave([]byte("invalid")))

	_, err = Get(vault, name)
	assert.Error(t, err)

	_, err = GetCheap(vault, name)
	assert.Error(t, err)

	_, err = List(vault)
	assert.Error(t, err)

	err = Rename(vault, name, "fail")
	assert.Error(t, err)
}

func TestProtoErrors(t *testing.T) {
	vault := setContext(t)

	name := "unformatted"
	err := vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.File.GetName())
		buf := make([]byte, 64)
		encBuf, _ := crypt.Encrypt(buf)
		return b.Put([]byte(name), encBuf)
	})
	assert.NoError(t, err, "Failed writing invalid type")

	_, err = Get(vault, name)
	assert.Error(t, err)

	_, err = GetCheap(vault, name)
	assert.Error(t, err)

	_, err = List(vault)
	assert.Error(t, err)

	err = Rename(vault, name, "fail")
	assert.Error(t, err)
}

func TestKeyError(t *testing.T) {
	vault := setContext(t)

	err := Create(vault, &protobuf.File{})
	assert.Error(t, err)

	err = Rename(vault, "", "")
	assert.Error(t, err)
}

func TestCompression(t *testing.T) {
	content := []byte{1, 2, 3}

	compressed, err := compress(content)
	assert.NoError(t, err)

	assert.NotEqual(t, content, compressed)

	decompressed, err := decompress(compressed)
	assert.NoError(t, err)

	assert.Equal(t, content, decompressed)
}

func TestCompressNil(t *testing.T) {
	var content []byte

	compressed, err := compress(content)
	assert.NoError(t, err)

	assert.Equal(t, content, compressed)
}

func TestDecompressNil(t *testing.T) {
	var content []byte

	decompressed, err := decompress(content)
	assert.NoError(t, err)

	assert.Equal(t, content, decompressed)
}

func setContext(t testing.TB) *bolt.DB {
	return vault_helper.SetContext(t, bucket.File.GetName())
}
