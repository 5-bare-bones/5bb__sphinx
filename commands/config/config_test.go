package config

import (
	"testing"

	command_helper "github.com/5-bare-bones/5bb__sphinx/commands"
	"github.com/5-bare-bones/5bb__sphinx/config"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	vault := command_helper.SetContext(t)
	config.SetFilename("./testdata/mock_config.yaml")

	cmd := NewCmd(vault)
	err := cmd.Execute()
	assert.NoError(t, err, "Failed reading config")
}

func TestReadError(t *testing.T) {
	vault := command_helper.SetContext(t)
	config.SetFilename("")

	cmd := NewCmd(vault)
	err := cmd.Execute()
	assert.Error(t, err)
}
