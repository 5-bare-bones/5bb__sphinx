// Package vault is the shared "vault factory": it creates and populates a
// brand-new sphinx vault from a set of records and a passphrase, using the
// exact db packages a student's binary uses. A forged vault is therefore
// byte-for-byte the same kind of artifact as a hand-built one.
//
// It is the single place tutor tooling (forge) — and, later, the bestow/receive
// transfer commands — go through to mint a vault.
package vault

import (
	"crypto/rand"
	"time"

	"github.com/5-bare-bones/5bb__sphinx/config"
	authDB "github.com/5-bare-bones/5bb__sphinx/db/auth"
	"github.com/5-bare-bones/5bb__sphinx/db/card"
	"github.com/5-bare-bones/5bb__sphinx/db/entry"
	"github.com/5-bare-bones/5bb__sphinx/db/totp"
	"github.com/5-bare-bones/5bb__sphinx/pb"

	"github.com/awnumar/memguard"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// Argon2 holds the key-derivation parameters a minted vault is locked with.
type Argon2 struct {
	Iterations uint32
	Memory     uint32 // kibibytes
	Threads    uint32
}

// DefaultArgon2 returns parameters tuned for seeding many fixture vaults
// quickly (64 MiB) rather than kure's heavyweight interactive default (1 GiB).
// The student binary reads whatever parameters the vault was minted with, so
// this only affects minting speed, never compatibility.
func DefaultArgon2() Argon2 {
	return Argon2{Iterations: 1, Memory: 64 << 10, Threads: 4}
}

// Seed is the content a new vault is populated with.
type Seed struct {
	Entries []*pb.Entry
	Cards   []*pb.Card
	TOTPs   []*pb.TOTP
}

// Create mints a new vault file at path, protected by passphrase, and populates
// it with seed.
//
// It mutates process-global auth configuration, because the crypt layer reads
// the active password and parameters from config. Minting is therefore
// sequential-only: do not call Create concurrently.
func Create(path, passphrase string, p Argon2, seed Seed) (err error) {
	if passphrase == "" {
		return errors.New("empty passphrase")
	}

	loadAuthConfig(passphrase, p)

	db, err := bolt.Open(path, 0o600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return errors.Wrap(err, "creating vault file")
	}
	defer func() {
		if cerr := db.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Random 32-byte vault key: obfuscates record names and is itself stored
	// encrypted under the passphrase-derived key.
	key := make([]byte, 32)
	if _, err = rand.Read(key); err != nil {
		return errors.Wrap(err, "generating vault key")
	}
	config.Set("auth.key", key)

	if err = authDB.Register(db, key, authDB.Params{Argon2: authDB.Argon2(p)}); err != nil {
		return errors.Wrap(err, "registering auth")
	}

	if len(seed.Entries) > 0 {
		if err = entry.Create(db, seed.Entries...); err != nil {
			return errors.Wrap(err, "seeding entries")
		}
	}
	for _, c := range seed.Cards {
		if err = card.Create(db, c); err != nil {
			return errors.Wrap(err, "seeding card")
		}
	}
	for _, t := range seed.TOTPs {
		if err = totp.Create(db, t); err != nil {
			return errors.Wrap(err, "seeding totp")
		}
	}

	return nil
}

// loadAuthConfig sets the process-global auth state that crypt and the db
// helpers read from. Mirrors auth.setAuthToConfig/setKeyToConfig.
func loadAuthConfig(passphrase string, p Argon2) {
	config.Set("auth", map[string]interface{}{
		"password":   memguard.NewEnclave([]byte(passphrase)),
		"iterations": p.Iterations,
		"memory":     p.Memory,
		"threads":    p.Threads,
	})
}
