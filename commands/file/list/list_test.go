package list

import (
	"testing"
	"time"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	"github.com/5-bare-bones/5bb__sphinx/vault/file"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	vault := command_helper.SetContext(t)

	err := file.Create(vault, &protobuf.File{Name: "test.txt"})
	assert.NoError(t, err, "Failed creating file")

	cases := []struct {
		desc   string
		name   string
		filter string
	}{
		{
			desc: "List one",
			name: "test.txt",
		},
		{
			desc:   "Filter by name",
			name:   "test*",
			filter: "true",
		},
		{
			desc: "List all",
			name: "",
		},
	}

	cmd := NewCmd(vault)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := cmd.Flags()
			cmd.SetArgs([]string{tc.name})
			f.Set("filter", tc.filter)

			err := cmd.Execute()
			assert.NoError(t, err)
		})
	}
}

func TestListErrors(t *testing.T) {
	vault := command_helper.SetContext(t)

	err := file.Create(vault, &protobuf.File{Name: "test.txt"})
	assert.NoError(t, err, "Failed creating file")

	cases := []struct {
		desc   string
		name   string
		filter string
	}{
		{
			desc: "File does not exist",
			name: "non-existent",
		},
		{
			desc:   "No files found",
			name:   "non-existent",
			filter: "true",
		},
		{
			desc:   "Filter syntax error",
			name:   "[error",
			filter: "true",
		},
	}

	cmd := NewCmd(vault)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			f := cmd.Flags()
			cmd.SetArgs([]string{tc.name})
			f.Set("filter", tc.filter)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}

func TestPrintFile(t *testing.T) {
	cases := []struct {
		desc string
		size int64
	}{
		{
			desc: "Bytes",
			size: 100,
		},
		{
			desc: "Kilo bytes",
			size: KB,
		},
		{
			desc: "Mega bytes",
			size: MB,
		},
		{
			desc: "Giga bytes",
			size: GB,
		},
		{
			desc: "Tera bytes",
			size: TB,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			printFile(&protobuf.FileCheap{Size: tc.size})
		})
	}
}

func TestPrintFileUpdatedAt(t *testing.T) {
	cases := []struct {
		desc string
		time int64
	}{
		{
			desc: "Without updated at",
			time: time.Time{}.Unix(),
		},
		{
			desc: "With updated at",
			time: 100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			printFile(&protobuf.FileCheap{UpdatedAt: tc.time})
		})
	}
}

func TestPostRun(t *testing.T) {
	NewCmd(nil).PostRun(nil, nil)
}
