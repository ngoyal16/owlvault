package encryptors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// AESCipher represents an AES encryption and decryption cipher.
type AESCipher struct {
	key   []byte
	block cipher.Block
}

// NewAESCipher creates a new instance of the AES encryption and decryption cipher.
func NewAESCipher(key []byte) (*AESCipher, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &AESCipher{key: key, block: block}, nil
}

// Encrypt encrypts data using AES encryption.
func (c *AESCipher) Encrypt(data string) (string, error) {
	plaintext := []byte(data)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(c.block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts data using AES decryption.
func (c *AESCipher) Decrypt(ciphertext string) (string, error) {
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if len(decodedCiphertext) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}
	iv := decodedCiphertext[:aes.BlockSize]
	ciphertext = decodedCiphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(c.block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	return string(ciphertext), nil
}
