package move

import (
	"strconv"
	"strings"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestMv(t *testing.T) {
	vault := command_helper.SetContext(t)

	oldName := "test.txt"
	newName := "renamed-test" // omit extension on purpose
	createFile(t, vault, oldName)

	cmd := NewCmd(vault)
	cmd.SetArgs([]string{oldName, newName})

	err := cmd.Execute()
	assert.NoError(t, err)

	_, err = file.GetCheap(vault, newName+".txt")
	assert.NoError(t, err, "Failed getting the renamed file")
}

func TestMvDir(t *testing.T) {
	vault := command_helper.SetContext(t)

	oldDir := "directory/"
	for i := range 2 {
		createFile(t, vault, oldDir+strconv.Itoa(i))
	}

	newDir := "folder/"
	cmd := NewCmd(vault)
	cmd.SetArgs([]string{oldDir, newDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	names, err := file.ListNames(vault)
	assert.NoError(t, err)

	for _, name := range names {
		if !strings.HasPrefix(name, newDir) {
			t.Errorf("%q wasn't moved into %q", name, newDir)
		}
	}
}

func TestMvFileIntoDir(t *testing.T) {
	vault := command_helper.SetContext(t)

	filename := "directory/test.csv"
	newDir := "folder/"
	createFile(t, vault, filename)

	cmd := NewCmd(vault)
	cmd.SetArgs([]string{filename, newDir})

	err := cmd.Execute()
	assert.NoError(t, err)

	newName := newDir + strings.Split(filename, "/")[1]
	_, err = file.GetCheap(vault, newName)
	assert.NoError(t, err, "Failed getting the renamed file")
}

func TestMvErrors(t *testing.T) {
	vault := command_helper.SetContext(t)
	createFile(t, vault, "exists")
	dir := "dir/"
	createFile(t, vault, dir+"1")

	cases := []struct {
		desc    string
		oldName string
		newName string
	}{
		{
			desc:    "Invalid new name",
			oldName: "exists",
			newName: "",
		},
		{
			desc:    "New name already exists",
			oldName: "exists",
			newName: "exists",
		},
		{
			desc:    "File does not exist",
			oldName: "non-existent",
			newName: "test.txt",
		},
		{
			desc:    "Dir into file",
			oldName: dir,
			newName: "test.txt",
		},
	}

	cmd := NewCmd(vault)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cmd.SetArgs([]string{tc.oldName, tc.newName})
			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestMissingArguments(t *testing.T) {
	vault := command_helper.SetContext(t)

	cmd := NewCmd(vault)
	cmd.SetArgs([]string{"oldName"})
	err := cmd.Execute()
	assert.Error(t, err)
}

func createFile(t *testing.T, vault *bolt.DB, name string) {
	err := file.Create(vault, &protobuf.File{Name: name})
	assert.NoError(t, err, "Failed creating file")
}
