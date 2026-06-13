package edit

import (
	"testing"

	cmdutil "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"

	"github.com/stretchr/testify/assert"
)

func TestEditErrors(t *testing.T) {
	db := cmdutil.SetContext(t)

	cases := []struct {
		desc string
		path string
	}{
		{
			desc: "Invalid filename",
			path: "\\",
		},
	}

	cmd := NewCmd(db)

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			config.SetFilename(tc.path)

			err := cmd.Execute()
			assert.Error(t, err)
		})
	}
}
