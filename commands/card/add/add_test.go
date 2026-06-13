package add

import (
	"bytes"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	vault := command_helper.SetContext(t)

	cases := []struct {
		desc string
		name string
	}{
		{
			desc: "Add",
			name: "test",
		},
		{
			desc: "Add2",
			name: "test2",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString("type\n123456789\n1234\n2021/06\nnotes<\n")
			cmd := NewCmd(vault, buf)
			cmd.SetArgs([]string{tc.name})

			err := cmd.Execute()
			assert.NoError(t, err)

			_, err = card.Get(vault, tc.name)
			assert.NoError(t, err, "Card wasn't created correctly")
		})
	}
}

func TestAddErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	err := card.Create(vault, &protobuf.Card{Name: "test"})
	assert.NoError(t, err)

	cases := []struct {
		desc string
		name string
	}{
		{
			desc: "Already exists",
			name: "test",
		},
		{
			desc: "Invalid name",
			name: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString("type\n123456789\n1234\n2021/06\nnotes<\n")
			cmd := NewCmd(vault, buf)
			cmd.SetArgs([]string{tc.name})

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestInput(t *testing.T) {
	vault := command_helper.SetContext(t)

	expected := &protobuf.Card{
		Name:         "test",
		Type:         "type",
		Number:       "123456789",
		SecurityCode: "1234",
		ExpireDate:   "2021/06",
		Notes:        "notes",
	}

	buf := bytes.NewBufferString("type\n123456789\n1234\n2021/06\nnotes<")

	got, err := input(vault, "test", buf)
	assert.NoError(t, err, "Failed creating the card")

	assert.Equal(t, expected, got)
}
