package jumbler

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		clear      string
		passphrase string
		error      string
	}{
		{
			clear: "", passphrase: "",
		},
		{
			clear: "", passphrase: "kjhdskjd",
		},
		{
			clear: "dumpling", passphrase: "abc",
		},
		{
			clear: "", passphrase: "abc",
		},
	}

	for _, tc := range tests {
		encryptGot, err := encrypt([]byte(tc.clear), tc.passphrase)
		if tc.error != "" {
			assert.ErrorContains(t, err, tc.error)
			continue
		}
		assert.NoError(t, err)

		clearGot, err := decrypt([]byte(encryptGot), tc.passphrase)
		assert.NoError(t, err)
		assert.Equal(t, tc.clear, string(clearGot))
	}
}

func TestEncryptDecryptFilename(t *testing.T) {
	tests := []struct {
		clear      string
		passphrase string
		error      string
	}{
		{
			clear: "dumpling", passphrase: "abc",
		},
		{
			clear: "jumbler-file.pdf", passphrase: "strong-password-#@$#",
		},
	}

	for _, tc := range tests {
		encryptGot, err := encryptFilename(tc.clear, tc.passphrase)
		if tc.error != "" {
			assert.ErrorContains(t, err, tc.error)
			continue
		}
		assert.NoError(t, err)
		assert.False(t, strings.Contains(encryptGot, "/"))

		clearGot, err := decryptFilename(encryptGot, tc.passphrase)
		assert.NoError(t, err)
		assert.Equal(t, tc.clear, string(clearGot))
	}
}
