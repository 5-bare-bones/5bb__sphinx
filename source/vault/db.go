package vault_helper

import (
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/crypt"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"

	"github.com/awnumar/memguard"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

const nullChar = string('\x00')

// Record is an interface that all sphinx objects implement.
type Record interface {
	GetName() string
	proto.Message
}

// Get retrieves a record from the vault, decrypts it and loads it into record.
func Get(vault *bolt.DB, name string, record Record) error {
	return vault.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(GetBucketName(record))

		xorName := XorName([]byte(name))
		encRecord := b.Get(xorName)
		if encRecord == nil {
			return errors.Errorf("record %q does not exist", name)
		}

		decRecord, err := crypt.Decrypt(encRecord)
		if err != nil {
			return errors.Wrap(err, "decrypt record")
		}

		if err := proto.Unmarshal(decRecord, record); err != nil {
			return errors.Wrap(err, "unmarshal record")
		}

		return nil
	})
}

// GetBucketName returns the bucket name depending on the type of the record passed.
func GetBucketName(r Record) []byte {
	switch r.(type) {
	case *protobuf.Card:
		return bucket.Card.GetName()
	case *protobuf.Entry:
		return bucket.Entry.GetName()
	case *protobuf.File, *protobuf.FileCheap:
		return bucket.File.GetName()
	case *protobuf.TOTP:
		return bucket.TOTP.GetName()
	default:
		memguard.SafePanic("invalid object: " + r.GetName())
		return nil
	}
}

// List returns a list of decrypted records from the vault.
func List[R Record](vault *bolt.DB, record R) ([]R, error) {
	tx, err := vault.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	b := tx.Bucket(GetBucketName(record))
	records := make([]R, 0, b.Stats().KeyN)

	err = b.ForEach(func(_, v []byte) error {
		decRecord, err := crypt.Decrypt(v)
		if err != nil {
			return errors.Wrap(err, "decrypt record")
		}

		if err := proto.Unmarshal(decRecord, record); err != nil {
			return errors.Wrap(err, "unmarshal record")
		}

		records = append(records, record)
		// Allocate a new protobuf object of type R
		record = record.ProtoReflect().New().Interface().(R)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return records, nil
}

// ListNames returns a list with all the records names.
func ListNames(vault *bolt.DB, bucketName []byte) ([]string, error) {
	tx, err := vault.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// b will be nil only if the user attempts to add
	// a record on registration as this method is being used
	// in checks previous to a command execution
	b := tx.Bucket(bucketName)
	if b == nil {
		return nil, nil
	}

	names := make([]string, 0, b.Stats().KeyN)
	_ = b.ForEach(func(k, _ []byte) error {
		// Xor record name to get the original one
		name := XorName(k)
		names = append(names, string(name))
		return nil
	})

	sort.Strings(names)

	return names, nil
}

// Put encrypts and saves a record into the vault.
func Put(b *bolt.Bucket, record Record) error {
	name := strings.ReplaceAll(record.GetName(), nullChar, "")
	if name == "" {
		return errors.New("record name is empty")
	}

	buf, err := proto.Marshal(record)
	if err != nil {
		return errors.Wrap(err, "marshal record")
	}

	encRecord, err := crypt.Encrypt(buf)
	if err != nil {
		return errors.Wrap(err, "encrypt record")
	}

	xorName := XorName([]byte(name))
	if err := b.Put(xorName, encRecord); err != nil {
		return errors.Wrap(err, "store record")
	}

	return nil
}

// Remove removes records from the vault.
func Remove(vault *bolt.DB, bucketName []byte, names ...string) error {
	if len(names) == 0 {
		return nil
	}

	return vault.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		for _, name := range names {
			xorName := XorName([]byte(name))
			if err := b.Delete(xorName); err != nil {
				return errors.Wrapf(err, "delete record %q", name)
			}
		}

		return nil
	})
}

// SetContext creates a bucket and its context to test the vault operations.
func SetContext(t testing.TB, bucketName []byte) *bolt.DB {
	vaultFile, err := os.CreateTemp("", "*")
	assert.NoError(t, err)

	db, err := bolt.Open(vaultFile.Name(), 0o600, &bolt.Options{Timeout: 1 * time.Second})
	assert.NoError(t, err, "Failed connecting to the vault")

	config.Reset()
	// Reduce argon2 parameters to speed up tests
	auth := map[string]interface{}{
		"password":   memguard.NewEnclave([]byte("1")),
		"iterations": 1,
		"memory":     1,
		"threads":    1,
		"key":        []byte("01234567890123456789012345678901"),
	}
	config.Set("auth", auth)

	err = db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(bucketName)
		if _, err := tx.CreateBucketIfNotExists(bucketName); err != nil {
			return errors.Wrapf(err, "couldn't create %q bucket", bucketName)
		}
		return nil
	})
	assert.NoError(t, err)

	t.Cleanup(func() {
		err := db.Close()
		assert.NoError(t, err, "Failed closing the vault")
	})

	return db
}

// XorName does the bitwise xor operation between a name and the authentication key.
func XorName(name []byte) []byte {
	key, ok := config.Get("auth.key").([]byte)
	if !ok {
		memguard.SafeExit(1)
	}
	xor := make([]byte, len(name))

	for i := range name {
		xor[i] = key[i%len(key)] ^ name[i]
	}

	return xor
}
