package backup

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"
	bolt "go.etcd.io/bbolt"

	"github.com/stretchr/testify/assert"
)

func TestBackupFile(t *testing.T) {
	vault := command_helper.SetContext(t)
	filename := "backup-test"
	name := "test"

	t.Cleanup(func() {
		err := os.Remove(filename)
		assert.NoError(t, err, "Failed removing the vault backup")
	})

	err := entry.Create(vault, &protobuf.Entry{Name: name})
	assert.NoError(t, err)

	cmd := NewCmd(vault)
	f := cmd.Flags()
	f.Set("path", filename)

	err = cmd.Execute()
	assert.NoError(t, err, "Failed creating the vault backup")

	newDB, err := bolt.Open(filename, 0o600, bolt.DefaultOptions)
	assert.NoError(t, err)

	e, err := entry.Get(newDB, name)
	assert.NoError(t, err)

	assert.Equal(t, name, e.Name)

	err = newDB.Close()
	assert.NoError(t, err)
}

func TestBackupServer(t *testing.T) {
	vault := command_helper.SetContext(t)

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "localhost:4000/", nil)
	assert.NoError(t, err, "Failed sending the request")

	hf := httpBackup(vault)
	hf.ServeHTTP(rec, req)

	res := rec.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	gotCt := res.Header.Get("Content-Type")
	expectedCt := "application/octet-stream"
	assert.Equal(t, expectedCt, gotCt)
}

func TestBackupErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	cases := []struct {
		desc string
		port string
		http string
		path string
	}{
		{
			desc: "HTTP",
			http: "true",
			port: "0",
		},
		{
			desc: "Invalid path",
			path: "",
		},
		{
			desc: "Mkdir error",
			path: "backup.go/",
		},
		{
			desc: "Mutually exclusive flags",
			http: "true",
			path: "backup.db",
		},
	}

	cmd := NewCmd(vault)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := cmd.Flags()
			f.Set("path", tc.path)
			f.Set("http", tc.http)
			f.Set("port", tc.port)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestWriteTo(t *testing.T) {
	vault := command_helper.SetContext(t)

	var buf bytes.Buffer
	err := writeTo(vault, &buf)
	assert.NoError(t, err, "Failed writing vault")

	assert.NotEqual(t, buf.Len(), 0)
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}
