package main

import (
	"bytes"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"
	"github.com/5-bare-bones/5bb__sphinx/vault/totp"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestXorNames(t *testing.T) {
	vault := command_helper.SetContext(t)
	key := []byte("0123456789")
	name := "test"
	xorName := "DTAG"
	config.Set("auth.key", key)

	createRecords(t, vault, name)

	buf := bytes.NewBufferString("y")
	err := xorNames(vault, buf)
	assert.NoError(t, err)

	_, err = entry.Get(vault, xorName)
	assert.NoError(t, err)

	_, err = card.Get(vault, xorName)
	assert.NoError(t, err)

	_, err = file.GetCheap(vault, xorName)
	assert.NoError(t, err)

	_, err = totp.Get(vault, xorName)
	assert.NoError(t, err)

	_, err = entry.Get(vault, name)
	assert.Error(t, err)

	_, err = card.Get(vault, name)
	assert.Error(t, err)

	_, err = file.Get(vault, name)
	assert.Error(t, err)

	_, err = totp.Get(vault, name)
	assert.Error(t, err)
}

func createRecords(t *testing.T, vault *bolt.DB, name string) {
	err := entry.Create(vault, &protobuf.Entry{Name: name})
	assert.NoError(t, err)

	err = card.Create(vault, &protobuf.Card{Name: name})
	assert.NoError(t, err)

	err = file.Create(vault, &protobuf.File{Name: name})
	assert.NoError(t, err)

	err = totp.Create(vault, &protobuf.TOTP{Name: name})
	assert.NoError(t, err)
}
