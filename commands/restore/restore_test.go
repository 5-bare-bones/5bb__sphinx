package restore

import (
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/bucket"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestLogs(t *testing.T) {
	vault := command_helper.SetContext(t)

	expected := &protobuf.Entry{
		Name:     "test",
		Username: "test@test.com",
		Password: "Pb*9' fxd%,IS:Zo_1JVw",
		Expires:  "Never",
	}
	err := entry.Create(vault, expected)
	assert.NoError(t, err)

	l, err := newLog(bucket.Entry.GetName())
	assert.NoError(t, err)
	defer l.Close()

	logs := []*log{l}
	err = writeLogs(vault, logs)
	assert.NoError(t, err, "Failed writing logs")

	err = readLogs(vault, logs)
	assert.NoError(t, err, "Failed reading logs")

	got, err := entry.Get(vault, expected.Name)
	assert.NoError(t, err, "Failed fetching entry")

	equal := proto.Equal(expected, got)
	assert.True(t, equal)
}

func TestReadLogs(t *testing.T) {
	vault := command_helper.SetContext(t)

	expected := &protobuf.Card{
		Name:         "testRead",
		Number:       "47964212",
		SecurityCode: "442",
	}
	err := card.Create(vault, expected)
	assert.NoError(t, err)

	l, err := newLog(bucket.Card.GetName())
	assert.NoError(t, err)
	defer l.Close()

	logs := []*log{l}
	err = readLogs(vault, logs)
	assert.NoError(t, err, "Failed reading logs")

	got, err := card.Get(vault, expected.Name)
	assert.NoError(t, err, "Failed fetching card")

	equal := proto.Equal(expected, got)
	assert.True(t, equal)
}

func TestWriteLogs(t *testing.T) {
	vault := command_helper.SetContext(t)

	l, err := newLog(bucket.Entry.GetName())
	assert.NoError(t, err)
	defer l.Close()

	err = writeLogs(vault, []*log{l})
	assert.NoError(t, err, "Failed writing logs")
}
