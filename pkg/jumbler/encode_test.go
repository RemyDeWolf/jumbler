package jumbler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeFilename(t *testing.T) {
	tests := []struct {
		clear   string
		encoded string
	}{
		{
			clear: "", encoded: "",
		},
		{
			clear: "abc", encoded: "YWJj",
		},
		{
			clear: "Screen Shot 2022-06-14 at 11.33.51 AM.png", encoded: "U2NyZWVuIFNob3QgMjAyMi0wNi0xNCBhdCAxMS4zMy41MSBBTS5wbmc=",
		},
		{
			clear: "a/b/c", encoded: "YS9iL2M=",
		},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.encoded, encodeFilename([]byte(tc.clear)))
		clearBytes, err := decodeFilename(tc.encoded)
		assert.NoError(t, err)
		assert.Equal(t, tc.clear, string(clearBytes))
	}
}
