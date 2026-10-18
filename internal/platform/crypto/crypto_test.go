package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

const testKey = "0123456789abcdef0123456789abcdef"

func TestEncryptDecryptRoundTrip(t *testing.T) {
	cases := map[string]string{
		"empty":       "",
		"short":       "hunter2",
		"block":       "0123456789abcdef",
		"block-plus":  "0123456789abcdef!",
		"multiblock":  "the quick brown fox jumps over the lazy dog",
		"unicode":     "密码口令",
		"long-random": string(bytes.Repeat([]byte("x"), 1000)),
	}
	for name, plaintext := range cases {
		t.Run(name, func(t *testing.T) {
			enc, err := EncryptAES([]byte(plaintext), testKey)
			require.NoError(t, err)
			require.Zero(t, len(enc)%16)
			if plaintext != "" {
				require.NotEqual(t, []byte(plaintext), enc)
			}

			dec, err := DecryptAES(enc, testKey)
			require.NoError(t, err)
			require.Equal(t, plaintext, string(dec))
		})
	}
}

// TestEncryptIsDeterministic locks in the ECB contract: identical input under an
// identical key must produce identical ciphertext so stored credentials stay
// comparable across restarts.
func TestEncryptIsDeterministic(t *testing.T) {
	a, err := EncryptAES([]byte("secret"), testKey)
	require.NoError(t, err)
	b, err := EncryptAES([]byte("secret"), testKey)
	require.NoError(t, err)
	require.Equal(t, a, b)
}

func TestDecryptRejectsMisalignedCiphertext(t *testing.T) {
	_, err := DecryptAES([]byte("not a multiple of 16"), testKey)
	require.Error(t, err)
}

func TestPkcs7PadAlwaysPads(t *testing.T) {
	// A full block still gains a whole block of padding, so the length is
	// unambiguous on decryption.
	padded := pkcs7Pad([]byte("0123456789abcdef"))
	require.Len(t, padded, 32)
	require.Equal(t, byte(16), padded[16])

	require.Equal(t, []byte("0123456789abcde\x01"), pkcs7Pad([]byte("0123456789abcde")))
}

func TestKeyUsesFirstBlock(t *testing.T) {
	// Only the first 16 bytes matter; a longer key must not change the result.
	long := testKey + "extra-key-material"
	a, err := EncryptAES([]byte("secret"), testKey)
	require.NoError(t, err)
	b, err := EncryptAES([]byte("secret"), long)
	require.NoError(t, err)
	require.Equal(t, a, b)
}
