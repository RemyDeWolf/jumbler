package jumbler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"path/filepath"
	"strings"
)

func encryptFilename(filename string, passphrase string) (string, error) {
	//encrypt the base name
	encryptedData, err := encrypt([]byte(filepath.Base(filename)), passphrase)
	if err != nil {
		return "", err
	}
	//encode in base64 and remove special chars not allowed for filename
	encodedBasename := encodeFilename(encryptedData)
	return filepath.Join(filepath.Dir(filename), encodedBasename+JumblerExt), nil
}

func decryptFilename(filename string, passphrase string) (string, error) {
	//remove extension
	filename = strings.TrimSuffix(filename, JumblerExt)
	//decode to encrypted data
	encryptedData, err := decodeFilename(filepath.Base(filename))
	if err != nil {
		return "", err
	}
	//decrypt and return full path
	decryptedData, err := decrypt([]byte(encryptedData), passphrase)
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(filename), string(decryptedData)), nil

}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
