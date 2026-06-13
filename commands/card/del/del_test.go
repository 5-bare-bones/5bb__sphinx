package del

import (
	"bytes"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	vault := command_helper.SetContext(t)

	names := []string{"test", "directory/test", "kure", "atoll"}
	for _, name := range names {
		err := card.Create(vault, &protobuf.Card{Name: name})
		assert.NoErrorf(t, err, "Failed creating %q", name)
	}

	cases := []struct {
		desc  string
		input string
		names []string
	}{
		{
			desc:  "Do not proceed",
			names: []string{"test"},
			input: "n",
		},
		{
			desc:  "Remove one card",
			names: []string{"test"},
			input: "y",
		},
		{
			desc:  "Remove multiple cards",
			names: []string{"kure", "atoll"},
			input: "y",
		},
		{
			desc:  "Remove directory",
			names: []string{"directory/"},
			input: "y",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.input)
			cmd := NewCmd(vault, buf)
			cmd.SetArgs(tc.names)

			err := cmd.Execute()
			assert.NoError(t, err)

			if tc.input == "y" {
				for _, name := range tc.names {
					_, err := card.Get(vault, name)
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestDeleteErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	name := "random"
	err := card.Create(vault, &protobuf.Card{Name: name})
	assert.NoErrorf(t, err, "Failed creating %q", name)

	cases := []struct {
		desc         string
		confirmation string
		names        []string
	}{
		{
			desc:  "Invalid name",
			names: []string{""},
		},
		{
			desc:         "Does not exists",
			names:        []string{"non-existent"},
			confirmation: "y",
		},
		{
			desc:  "Second name does not exist",
			names: []string{"random", "non-existent"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.confirmation)
			cmd := NewCmd(vault, buf)
			cmd.SetArgs(tc.names)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}
