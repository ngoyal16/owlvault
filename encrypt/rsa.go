package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// RSAEncryptor implements the Encryptor interface for RSA encryption.
type RSAEncryptor struct {
	privateKey []byte
	publicKey  []byte
}

// NewRSAEncryptor creates a new instance of RSAEncryptor.
func NewRSAEncryptor() (*RSAEncryptor, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: publicKeyBytes})
	return &RSAEncryptor{
		privateKey: privateKeyPEM,
		publicKey:  publicKeyPEM,
	}, nil
}

// Encrypt encrypts the data using RSA encryption.
func (e *RSAEncryptor) Encrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(e.publicKey)
	if block == nil {
		return nil, errors.New("failed to decode public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// Decrypt decrypts the data using RSA encryption.
func (e *RSAEncryptor) Decrypt(data []byte) ([]byte, error) {
	block, _ := pem.Decode(e.privateKey)
	if block == nil {
		return nil, errors.New("failed to decode private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
