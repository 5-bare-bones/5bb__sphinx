package edit

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestEditErrors(t *testing.T) {
	vault := command_helper.SetContext(t)
	createEntry(t, vault, "test")

	cases := []struct {
		set  func()
		desc string
		name string
		it   bool
	}{
		{
			desc: "Invalid name",
			name: "",
		},
		{
			desc: "Executable not found",
			name: "test",
			it:   true,
			set: func() {
				config.Set("editor", "non-existent")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.set != nil {
				tc.set()
			}
			cmd := NewCmd(vault)
			cmd.SetArgs([]string{tc.name})
			f := cmd.Flags()
			f.Set("it", strconv.FormatBool(tc.it))

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestCreateTempFile(t *testing.T) {
	e := &protobuf.Entry{
		Name:    "test-create-file",
		Expires: "Never",
	}

	filename, err := createTempFile(e)
	assert.NoError(t, err, "Failed creating the file")
	defer os.Remove(filename)

	content, err := os.ReadFile(filename)
	assert.NoError(t, err, "Failed reading the file")

	var got protobuf.Entry
	err = json.Unmarshal(content, &got)
	assert.NoError(t, err, "Failed reading the file")

	assert.Equal(t, e, &got)
}

func TestReadTmpFile(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		_, err := readTmpFile("testdata/test_read.json")
		assert.NoError(t, err)
	})

	t.Run("Errors", func(t *testing.T) {
		dir, _ := os.Getwd()
		os.Chdir("testdata")
		// Go back to the initial directory
		defer os.Chdir(dir)

		cases := []struct {
			desc     string
			filename string
		}{
			{
				desc:     "Does not exists",
				filename: "does_not_exists.json",
			},
			{
				desc:     "EOF",
				filename: "test_read_EOF.json",
			},
		}

		for _, tc := range cases {
			t.Run(tc.desc, func(t *testing.T) {
				_, err := readTmpFile(tc.filename)
				assert.Error(t, err)
			})
		}
	})
}

func TestUpdateEntry(t *testing.T) {
	vault := command_helper.SetContext(t)
	name := "test_update"
	createEntry(t, vault, name)

	newName := "new_name"
	newEntry := &protobuf.Entry{
		Name:     newName,
		Username: "test",
		Password: "q8rvq63r/q",
		URL:      "https://www.github.com/5-bare-bones/5bb__sphinx",
		Expires:  "02/12/2023",
		Notes:    "",
	}

	err := updateEntry(vault, name, newEntry)
	assert.NoError(t, err)

	e, err := entry.Get(vault, newName)
	assert.NoError(t, err)

	assert.NotEqual(t, newEntry, e)

	t.Run("Invalid name", func(t *testing.T) {
		newEntry.Name = ""
		err := updateEntry(vault, "fail", newEntry)
		assert.Error(t, err)
	})
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}

func createEntry(t *testing.T, vault *bolt.DB, name string) {
	t.Helper()
	e := &protobuf.Entry{
		Name:    name,
		Expires: "Never",
	}

	err := entry.Create(vault, e)
	assert.NoError(t, err)
}
