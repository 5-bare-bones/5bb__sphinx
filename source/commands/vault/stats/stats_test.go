package stats

import (
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"

	"github.com/stretchr/testify/assert"
)

func TestStats(t *testing.T) {
	vault := command_helper.SetContext(t)

	t.Run("Success", func(t *testing.T) {
		cmd := NewCmd(vault)
		err := cmd.Execute()
		assert.NoError(t, err)
	})

	t.Run("vault connection closed", func(t *testing.T) {
		vault.Close()
		err := NewCmd(vault).Execute()
		assert.Error(t, err)
	})
}
