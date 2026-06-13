package vault

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/5-bare-bones/5bb__sphinx/config"
	"github.com/5-bare-bones/5bb__sphinx/crypt"
	"github.com/5-bare-bones/5bb__sphinx/protobuf"
	authVault "github.com/5-bare-bones/5bb__sphinx/vault/auth"
	"github.com/5-bare-bones/5bb__sphinx/vault/card"
	"github.com/5-bare-bones/5bb__sphinx/vault/entry"

	"github.com/awnumar/memguard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

// login reproduces auth.Login non-interactively: derive the password key,
// decrypt the stored vault key, and load both into config — exactly what a
// fresh student session does when opening a forged vault.
func login(t *testing.T, vault *bolt.DB, passphrase string) {
	t.Helper()
	config.Set("auth", nil)
	config.Set("auth.key", nil)

	params, err := authVault.GetParams(vault)
	require.NoError(t, err)

	config.Set("auth", map[string]interface{}{
		"password":   memguard.NewEnclave([]byte(passphrase)),
		"iterations": params.Argon2.Iterations,
		"memory":     params.Argon2.Memory,
		"threads":    params.Argon2.Threads,
	})

	key, err := crypt.Decrypt(params.AuthKey)
	require.NoError(t, err, "wrong passphrase would fail here")
	config.Set("auth.key", key)
}

func TestCreateRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alice.db")
	const pass = "correct horse battery staple"

	seed := Seed{
		Entries: []*protobuf.Entry{
			{Name: "github", Username: "alice", Password: "hunter2", URL: "https://github.com"},
		},
		Cards: []*protobuf.Card{
			{Name: "visa", Type: "Visa", Number: "4111111111111111", SecurityCode: "123", ExpireDate: "12/2030"},
		},
	}

	require.NoError(t, Create(path, pass, DefaultArgon2(), seed))

	vault, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	require.NoError(t, err)
	defer vault.Close()

	// Open it as a brand-new session would.
	login(t, vault, pass)

	e, err := entry.Get(vault, "github")
	require.NoError(t, err)
	assert.Equal(t, "alice", e.Username)
	assert.Equal(t, "hunter2", e.Password)

	c, err := card.Get(vault, "visa")
	require.NoError(t, err)
	assert.Equal(t, "4111111111111111", c.Number)
}

func TestCreateRejectsEmptyPassphrase(t *testing.T) {
	err := Create(filepath.Join(t.TempDir(), "x.db"), "", DefaultArgon2(), Seed{})
	require.Error(t, err)
}

func TestWrongPassphraseFailsLogin(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bob.db")
	require.NoError(t, Create(path, "right answer", DefaultArgon2(), Seed{}))

	vault, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	require.NoError(t, err)
	defer vault.Close()

	config.Set("auth", nil)
	config.Set("auth.key", nil)
	params, err := authVault.GetParams(vault)
	require.NoError(t, err)
	config.Set("auth", map[string]interface{}{
		"password":   memguard.NewEnclave([]byte("WRONG answer")),
		"iterations": params.Argon2.Iterations,
		"memory":     params.Argon2.Memory,
		"threads":    params.Argon2.Threads,
	})
	_, err = crypt.Decrypt(params.AuthKey)
	require.Error(t, err)
}
