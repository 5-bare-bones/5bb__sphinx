package forge

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleManifest = `
argon2:
  iterations: 1
  memory: 8192
  threads: 1
students:
  - name: alice
    passphrase: "correct horse battery staple"
    entries:
      - name: github
        username: alice
        password: hunter2
  - name: bob
    passphrase: "tower bridge canyon"
`

func runForgeCmd(t *testing.T, manifest, outDir string) error {
	t.Helper()
	cmd := NewCmd()
	cmd.SetArgs([]string{manifest, "-o", outDir})
	cmd.SetOut(io.Discard)
	return cmd.Execute()
}

func writeManifest(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "world.yaml")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestForgeMintsVaults(t *testing.T) {
	manifest := writeManifest(t, sampleManifest)
	out := t.TempDir()

	require.NoError(t, runForgeCmd(t, manifest, out))

	for _, name := range []string{"alice.db", "bob.db"} {
		_, err := os.Stat(filepath.Join(out, name))
		assert.NoError(t, err, "expected %s to be minted", name)
	}
}

func TestForgeRefusesOverwrite(t *testing.T) {
	manifest := writeManifest(t, sampleManifest)
	out := t.TempDir()

	require.NoError(t, runForgeCmd(t, manifest, out))
	// Second run hits the existing files and must refuse.
	err := runForgeCmd(t, manifest, out)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already exists")
}

func TestForgeRejectsDuplicateNames(t *testing.T) {
	dup := `
students:
  - name: alice
    passphrase: a
  - name: alice
    passphrase: b
`
	err := runForgeCmd(t, writeManifest(t, dup), t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}

func TestForgeRejectsMissingPassphrase(t *testing.T) {
	m := `
students:
  - name: alice
`
	err := runForgeCmd(t, writeManifest(t, m), t.TempDir())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "passphrase")
}
