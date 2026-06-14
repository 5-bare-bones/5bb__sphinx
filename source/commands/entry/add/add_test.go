package add

import (
	"bytes"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	vault := command_helper.SetContext(t)

	cases := []struct {
		desc   string
		name   string
		length string
		format string
	}{
		{
			desc:   "Add",
			name:   "test",
			length: "18",
			format: "1,2,3,4,5",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString("username\nurl\n2024-05-03\nnotes<")
			cmd := NewCmd(vault, buf)
			cmd.SetArgs([]string{tc.name})

			f := cmd.Flags()
			f.Set("length", tc.length)
			f.Set("format", tc.format)

			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestAddErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	err := entry.Create(vault, &protobuf.Entry{Name: "test"})
	assert.NoError(t, err)

	cases := []struct {
		desc   string
		name   string
		length string
		levels string
	}{
		{
			desc: "Invalid name",
			name: "",
		},
		{
			desc:   "Invalid length",
			name:   "test-invalid-length",
			length: "0",
		},
		{
			desc:   "Already exists",
			name:   "test",
			length: "10",
		},
		{
			desc:   "Invalid levels",
			name:   "test-levels",
			length: "6",
			levels: "1,2,3,7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			buf := bytes.NewBufferString("username\nurl\n2024-05-03\nnotes<")
			cmd := NewCmd(vault, buf)
			cmd.SetArgs([]string{tc.name})
			f := cmd.Flags()
			f.Set("length", tc.length)
			f.Set("levels", tc.levels)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestGenPassword(t *testing.T) {
	cases := []struct {
		opts *addOptions
		desc string
	}{
		{
			desc: "Generate password",
			opts: &addOptions{
				length: 10,
				levels: []int{1, 2, 3},
				repeat: false,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := genPassword(tc.opts)
			assert.NoError(t, err, "Failed generating password")
		})
	}
}

func TestGenPasswordErrors(t *testing.T) {
	cases := []struct {
		opts *addOptions
		desc string
	}{
		{
			desc: "Invalid length",
			opts: &addOptions{
				length: 0,
				levels: []int{1, 2, 3},
			},
		},
		{
			desc: "Include and exclude error",
			opts: &addOptions{
				length:  5,
				include: "a",
				exclude: "a",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := genPassword(tc.opts)
			assert.Error(t, err)
		})
	}
}

func TestEntryInput(t *testing.T) {
	expected := &protobuf.Entry{
		Name:     "test",
		Username: "username",
		URL:      "url",
		Notes:    "notes",
		Expires:  "2024-05-03",
	}

	buf := bytes.NewBufferString("username\nurl\n2024-05-03\nnotes<")
	got, err := entryInput(buf, "test", false)
	assert.NoError(t, err, "Failed creating entry")

	// As it's randomly generated, use the same one
	expected.Password = got.Password
	assert.Equal(t, expected, got)
}

func TestInvalidExpirationTime(t *testing.T) {
	buf := bytes.NewBufferString("username\nurl\nnotes\ninvalid<\n")

	_, err := entryInput(buf, "test", false)
	assert.Error(t, err)
}

func TestPostRun(t *testing.T) {
	NewCmd(nil, nil).PostRun(nil, nil)
}
