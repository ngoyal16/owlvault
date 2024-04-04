package encrypt

import (
	"errors"
	"github.com/ngoyal16/owlvault/config"
)

// Encryptor defines the interface for encryption algorithms.
type Encryptor interface {
	Encrypt([]byte, []byte) ([]byte, error)
	Decrypt([]byte, []byte) ([]byte, error)
}

// EncryptorType represents the type of encryptor.
type EncryptorType string

const (
	// AES represents the AES encryption algorithm.
	AES EncryptorType = "aes"
	// RSA represents the RSA encryption algorithm.
	RSA EncryptorType = "rsa"
	// Add more encryption algorithms as needed
)

// NewEncryptor creates a new instance of an encryptor based on the provided type.
func NewEncryptor(cfg *config.Config) (Encryptor, error) {
	encryptorType := EncryptorType(cfg.Encryptor.Type)
	switch encryptorType {
	case AES:
		return NewAESEncryptor()
	case RSA:
		return NewRSAEncryptor()
	// Add cases for other encryption algorithms as needed
	default:
		return nil, errors.New("unsupported encryptor type")
	}
}
