package copy

import (
	"strconv"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	bolt "go.etcd.io/bbolt"

	"github.com/atotto/clipboard"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	if clipboard.Unsupported {
		t.Skip("No clipboard utilities available")
	}

	vault := command_helper.SetContext(t)
	e := createEntry(t, vault)

	cases := []struct {
		desc         string
		value        string
		timeout      string
		copyUsername bool
	}{
		{
			desc:  "Copy password",
			value: e.Password,
		},
		{
			desc:         "Copy username",
			value:        e.Username,
			copyUsername: true,
		},
		{
			desc:    "Copy with timeout",
			value:   "",
			timeout: "1ms",
		},
	}

	cmd := NewCmd(vault)
	f := cmd.Flags()

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd.SetArgs([]string{e.Name})
			f.Set("timeout", tc.timeout)
			f.Set("username", strconv.FormatBool(tc.copyUsername))

			err := cmd.Execute()
			assert.NoError(t, err)

			got, err := clipboard.ReadAll()
			assert.NoError(t, err, "Failed reading from clipboard")

			assert.Equal(t, tc.value, got)
		})
	}
}

func TestCopyWithConfigTimeout(t *testing.T) {
	if clipboard.Unsupported {
		t.Skip("No clipboard utilities available")
	}
	vault := command_helper.SetContext(t)
	e := createEntry(t, vault)

	config.Set("clipboard.timeout", "1ns")
	cmd := NewCmd(vault)
	cmd.SetArgs([]string{e.Name})

	err := cmd.Execute()
	assert.NoError(t, err, "Failed to copy password to clipboard")

	got, err := clipboard.ReadAll()
	assert.NoError(t, err, "Failed reading from clipboard")

	assert.Empty(t, got)
}

func TestCopyErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	cases := []struct {
		desc string
		name string
	}{
		{desc: "Non-existent", name: "non-existent"},
		{desc: "Invalid name", name: ""},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd := NewCmd(vault)
			cmd.SetArgs([]string{tc.name})

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}

func createEntry(t *testing.T, vault *bolt.DB) *protobuf.Entry {
	t.Helper()

	e := &protobuf.Entry{
		Name:     "test",
		Username: "Go",
		Password: "Gopher",
		Expires:  "Never",
	}
	err := entry.Create(vault, e)
	assert.NoError(t, err)

	return e
}
