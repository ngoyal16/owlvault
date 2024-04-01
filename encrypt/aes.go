package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AESEncryptor implements the Encryptor interface for AES encryption.
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor creates a new instance of AESEncryptor.
func NewAESEncryptor() (*AESEncryptor, error) {
	key := make([]byte, 32) // AES-256
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return &AESEncryptor{key: key}, nil
}

// Encrypt encrypts the data using AES encryption.
func (e *AESEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return ciphertext, nil
}

// Decrypt decrypts the data using AES encryption.
func (e *AESEncryptor) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}
