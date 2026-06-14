// Package riddle demonstrates a knowledge-derived key: the (normalized) answer
// to a riddle IS the key that decrypts a sealed token. There is no stored hash
// to compare against — a wrong answer simply derives a different key and the
// authenticated decryption fails. Tutor-only tooling.
package riddle

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
)

// riddleFile is the on-disk, shareable riddle. None of it reveals the answer.
type riddleFile struct {
	Prompt     string     `json:"prompt"`
	Argon2     argon2Spec `json:"argon2"`
	Salt       []byte     `json:"salt"`
	Nonce      []byte     `json:"nonce"`
	Ciphertext []byte     `json:"ciphertext"`
}

type argon2Spec struct {
	Iterations uint32 `json:"iterations"`
	Memory     uint32 `json:"memory"`
	Threads    uint32 `json:"threads"`
}

func defaultArgon2() argon2Spec {
	return argon2Spec{Iterations: 3, Memory: 64 << 10, Threads: 4}
}

// NewCmd returns the riddle command tree.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "riddle",
		Short: "Seal and solve knowledge-derived-key riddles (tutor only)",
		Example: `
* Seal a secret behind an answer
sphinx riddle seal --prompt "I speak without a mouth" --answer "an echo" --secret <token> --out echo.riddle

* Solve it
sphinx riddle solve echo.riddle`,
	}
	cmd.AddCommand(newSealCmd(), newSolveCmd())
	return cmd
}

// normalize makes the answer robust to incidental variation: lowercase, trim,
// and collapse internal whitespace runs to single spaces. The same normalized
// string must be reproduced exactly to derive the key.
func normalize(answer string) string {
	return strings.ToLower(strings.Join(strings.Fields(answer), " "))
}

// deriveKey turns a normalized answer + salt into a 32-byte AES key.
func deriveKey(answer string, salt []byte, p argon2Spec) []byte {
	return argon2.IDKey([]byte(normalize(answer)), salt, p.Iterations, p.Memory, uint8(p.Threads), 32)
}

// seal encrypts secret so that only the answer can recover it.
func seal(prompt, answer string, secret []byte) (riddleFile, error) {
	p := defaultArgon2()
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return riddleFile{}, err
	}

	gcm, err := newGCM(deriveKey(answer, salt, p))
	if err != nil {
		return riddleFile{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return riddleFile{}, err
	}

	return riddleFile{
		Prompt:     prompt,
		Argon2:     p,
		Salt:       salt,
		Nonce:      nonce,
		Ciphertext: gcm.Seal(nil, nonce, secret, nil),
	}, nil
}

// solve attempts to recover the secret using answer. A wrong answer yields an
// error, not a different-but-plausible plaintext (GCM authentication fails).
func solve(rf riddleFile, answer string) ([]byte, error) {
	gcm, err := newGCM(deriveKey(answer, rf.Salt, rf.Argon2))
	if err != nil {
		return nil, err
	}
	if len(rf.Nonce) != gcm.NonceSize() {
		return nil, errors.New("corrupt riddle: bad nonce")
	}
	secret, err := gcm.Open(nil, rf.Nonce, rf.Ciphertext, nil)
	if err != nil {
		return nil, errors.New("wrong answer")
	}
	return secret, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
