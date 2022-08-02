package jumbler

import (
	"encoding/base64"
	"strings"
)

//encode in base64 and remove special chars not allowed for filename
func encodeFilename(b []byte) string {
	result := base64.StdEncoding.EncodeToString(b)
	result = handleSpecialChar(result, true)
	return result
}

func decodeFilename(s string) ([]byte, error) {
	s = handleSpecialChar(s, false)
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return data, nil
}

var specialCars = map[string]string{"/": "{slash}"}

func handleSpecialChar(s string, encode bool) string {
	for k, v := range specialCars {
		if !encode {
			k, v = v, k
		}
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
