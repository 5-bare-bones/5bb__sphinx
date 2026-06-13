package main

import (
	"bytes"
	"testing"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/db/entry"
	"github.com/5-bare-bones/5bb__sphinx/db/file"
	"github.com/5-bare-bones/5bb__sphinx/db/totp"
	"github.com/5-bare-bones/5bb__sphinx/pb"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestXorNames(t *testing.T) {
	db := cmdutil.SetContext(t)
	key := []byte("0123456789")
	name := "test"
	xorName := "DTAG"
	config.Set("auth.key", key)

	createRecords(t, db, name)

	buf := bytes.NewBufferString("y")
	err := xorNames(db, buf)
	assert.NoError(t, err)

	_, err = entry.Get(db, xorName)
	assert.NoError(t, err)

	_, err = card.Get(db, xorName)
	assert.NoError(t, err)

	_, err = file.GetCheap(db, xorName)
	assert.NoError(t, err)

	_, err = totp.Get(db, xorName)
	assert.NoError(t, err)

	_, err = entry.Get(db, name)
	assert.Error(t, err)

	_, err = card.Get(db, name)
	assert.Error(t, err)

	_, err = file.Get(db, name)
	assert.Error(t, err)

	_, err = totp.Get(db, name)
	assert.Error(t, err)
}

func createRecords(t *testing.T, db *bolt.DB, name string) {
	err := entry.Create(db, &pb.Entry{Name: name})
	assert.NoError(t, err)

	err = card.Create(db, &pb.Card{Name: name})
	assert.NoError(t, err)

	err = file.Create(db, &pb.File{Name: name})
	assert.NoError(t, err)

	err = totp.Create(db, &pb.TOTP{Name: name})
	assert.NoError(t, err)
}
