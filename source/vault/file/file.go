package file

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/5-bare-bones/5bb__sphinx/crypt"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	vault_helper "github.com/5-bare-bones/5bb__sphinx/vault"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

// Create a new file with its content compressed.
func Create(vault *bolt.DB, file *protobuf.File) error {
	compressedContent, err := compress(file.Content)
	if err != nil {
		return err
	}
	file.Content = compressedContent

	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.File.GetName())
		return vault_helper.Put(b, file)
	})
}

// Get retrieves the file with the specified name.
func Get(vault *bolt.DB, name string) (*protobuf.File, error) {
	file := &protobuf.File{}
	if err := vault_helper.Get(vault, name, file); err != nil {
		return nil, err
	}

	decompressedContent, err := decompress(file.Content)
	if err != nil {
		return nil, err
	}
	file.Content = decompressedContent

	return file, nil
}

// GetCheap is like Get but without getting the file content.
func GetCheap(vault *bolt.DB, name string) (*protobuf.FileCheap, error) {
	file := &protobuf.FileCheap{}
	if err := vault_helper.Get(vault, name, file); err != nil {
		return nil, err
	}

	return file, nil
}

// List returns a slice with all the files stored in the file bucket.
func List(vault *bolt.DB) ([]*protobuf.File, error) {
	tx, err := vault.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket(bucket.File.GetName())
	files := make([]*protobuf.File, 0, b.Stats().KeyN)

	err = b.ForEach(func(k, v []byte) error {
		file := &protobuf.File{}

		decFile, err := crypt.Decrypt(v)
		if err != nil {
			return errors.Wrap(err, "decrypt file")
		}

		if err := proto.Unmarshal(decFile, file); err != nil {
			return errors.Wrap(err, "unmarshal file")
		}

		decompressedContent, err := decompress(file.Content)
		if err != nil {
			return err
		}
		file.Content = decompressedContent
		files = append(files, file)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ListNames returns a slice with all the files names.
func ListNames(vault *bolt.DB) ([]string, error) {
	return vault_helper.ListNames(vault, bucket.File.GetName())
}

// Remove removes one or more files from the vault.
func Remove(vault *bolt.DB, names ...string) error {
	return vault_helper.Remove(vault, bucket.File.GetName(), names...)
}

// Rename recreates a file with a new key and deletes the old one.
func Rename(vault *bolt.DB, oldName, newName string) error {
	return vault.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket.File.GetName())

		file, err := Get(vault, oldName)
		if err != nil {
			return err
		}

		compressedContent, err := compress(file.Content)
		if err != nil {
			return err
		}
		file.Name = newName
		file.Content = compressedContent

		if err := vault_helper.Put(b, file); err != nil {
			return err
		}

		xorName := vault_helper.XorName([]byte(oldName))
		return b.Delete(xorName)
	})
}

func compress(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, nil
	}

	var gzipBuf bytes.Buffer
	gw := gzip.NewWriter(&gzipBuf)

	if _, err := gw.Write(content); err != nil {
		return nil, errors.Wrap(err, "compress content")
	}

	if err := gw.Close(); err != nil {
		return nil, errors.Wrap(err, "close gzip writer")
	}

	return gzipBuf.Bytes(), nil
}

func decompress(content []byte) ([]byte, error) {
	if len(content) == 0 {
		return nil, nil
	}

	compressed := bytes.NewBuffer(content)
	gr, err := gzip.NewReader(compressed)
	if err != nil {
		return nil, errors.Wrap(err, "decompress content")
	}
	defer gr.Close()

	var decompressed bytes.Buffer
	if _, err = io.Copy(&decompressed, gr); err != nil {
		return nil, errors.Wrap(err, "copy decompressed content")
	}

	return decompressed.Bytes(), nil
}
