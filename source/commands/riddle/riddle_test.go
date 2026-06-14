package riddle

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSealSolveRoundTrip(t *testing.T) {
	rf, err := seal("I speak without a mouth", "an echo", []byte("tier-3-token"))
	require.NoError(t, err)

	got, err := solve(rf, "an echo")
	require.NoError(t, err)
	assert.Equal(t, "tier-3-token", string(got))
}

func TestSolveNormalizesAnswer(t *testing.T) {
	rf, err := seal("prompt", "An  ECHO ", []byte("secret"))
	require.NoError(t, err)

	// Different casing/spacing/trailing space must still solve it.
	got, err := solve(rf, "  an echo")
	require.NoError(t, err)
	assert.Equal(t, "secret", string(got))
}

func TestSolveWrongAnswerFails(t *testing.T) {
	rf, err := seal("prompt", "the moon", []byte("secret"))
	require.NoError(t, err)

	_, err = solve(rf, "the sun")
	require.EqualError(t, err, "wrong answer")
}

func TestNormalize(t *testing.T) {
	cases := map[string]string{
		"  Hello   World  ": "hello world",
		"ECHO":              "echo",
		"a\tb\nc":           "a b c",
	}
	for in, want := range cases {
		assert.Equal(t, want, normalize(in))
	}
}
