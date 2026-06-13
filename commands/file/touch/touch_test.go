package touch

import (
	"bytes"
	"os"
	"testing"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/db/file"
	"github.com/5-bare-bones/5bb__sphinx/pb"

	"github.com/stretchr/testify/assert"
	bolt "go.etcd.io/bbolt"
)

func TestTouch(t *testing.T) {
	db := cmdutil.SetContext(t)
	createTestFiles(t, db)

	rootDir, err := os.Getwd()
	assert.NoError(t, err)

	cases := []struct {
		desc  string
		path  string
		names []string
	}{
		{
			desc:  "Create one file",
			names: []string{"test.txt"},
			path:  "testdata/create/one",
		},
		{
			desc: "Create all files",
			path: "testdata/create",
		},
		{
			desc:  "Create multiple files",
			names: []string{"testAll/file.csv", "testAll/file.pdf"},
			path:  "testdata/create",
		},
		{
			desc:  "Create a directory",
			names: []string{"testdata/subfolder/"},
			path:  "testdata/create",
		},
	}

	cmd := NewCmd(db)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			var stderr bytes.Buffer
			cmd.SetErr(&stderr)
			cmd.SetArgs(tc.names)
			f := cmd.Flags()
			f.Set("path", "../"+tc.path)

			err := cmd.Execute()
			assert.NoError(t, err)

			assert.LessOrEqual(t, stderr.Len(), 0, "Expected nothing on stderr")

			// Go back to the directory where the test was executed
			os.Chdir(rootDir)
		})
	}

	err = os.RemoveAll("../testdata/create")
	assert.NoError(t, err, "Failed removing the folder containing created files")
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}

func createTestFiles(t *testing.T, db *bolt.DB) {
	t.Helper()

	names := []string{"test.txt", "testAll/file.csv", "testAll/file.pdf", "testdata/subfolder/file.txt"}
	for _, name := range names {
		err := file.Create(db, &pb.File{Name: cmdutil.NormalizeName(name)})
		assert.NoErrorf(t, err, "Failed creating %q", name)
	}
}
