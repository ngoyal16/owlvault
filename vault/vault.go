package vault

import (
	"encoding/base64"
	"fmt"

	"github.com/ngoyal16/owlvault/encrypt"
	"github.com/ngoyal16/owlvault/storage"
)

// OwlVault represents the key vault service.
type OwlVault struct {
	encryptor encrypt.Encryptor
	storage   storage.Storage
}

// NewOwlVault creates a new instance of OwlVault with the given storage.
func NewOwlVault(storage storage.Storage, encryptor encrypt.Encryptor) *OwlVault {
	return &OwlVault{
		storage:   storage,
		encryptor: encryptor,
	}
}

// StoreKey stores the key-value pair in the vault.
func (ov *OwlVault) StoreKey(key, value string) error {

	version, err := ov.storage.LatestVersion(key)
	if err != nil {
		return err
	}

	version += 1

	// Implement logic to store key-value pair in the storage backend
	encryptedValue, err := ov.encryptor.Encrypt([]byte(value))
	if err != nil {
		return err
	}

	// Convert encrypted value to base64 encoding
	base64Value := base64.StdEncoding.EncodeToString(encryptedValue)

	// Implement logic to store key-value pair in the storage backend
	if err := ov.storage.Store(key, base64Value, version); err != nil {
		return fmt.Errorf("failed to store key-value pair: %v", err)
	}
	return nil
}

// RetrieveKey retrieves the value for the specified key and version from the vault.
func (ov *OwlVault) RetrieveKey(key string, version int) (string, error) {
	// Implement logic to retrieve value from the storage backend
	base64Value, err := ov.storage.Retrieve(key, version)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded value
	encryptedValue, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 value: %v", err)
	}

	// Decrypt the retrieved value
	decryptedValue, err := ov.encryptor.Decrypt(encryptedValue)
	if err != nil {
		return "", err
	}

	return string(decryptedValue), nil
}

// Additional methods for OwlVault can be added as needed.
