package vault

import (
	"encoding/base64"
	"encoding/json"
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

// Store stores the key-value pair in the vault.
func (ov *OwlVault) Store(key string, data map[string]interface{}) error {
	// Check if version exists
	version, err := ov.storage.LatestVersion(key)
	if err != nil {
		return err
	}

	version += 1

	b, err := json.Marshal(&data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %w", err)
	}

	// Implement logic to store key-value pair in the storage backend
	encryptedValue, err := ov.encryptor.Encrypt(b)
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

// RetrieveVersion retrieves the value for the specified key and version from the vault.
func (ov *OwlVault) RetrieveVersion(key string, version int) (map[string]interface{}, error) {
	var data map[string]interface{}

	// Implement logic to retrieve value from the storage backend
	base64Value, err := ov.storage.Retrieve(key, version)
	if err != nil {
		return nil, err
	}

	// Decode the base64-encoded value
	encryptedValue, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 value: %v", err)
	}

	// Decrypt the retrieved value
	decrypted, err := ov.encryptor.Decrypt(encryptedValue)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(decrypted, &data)

	return data, nil
}

// RetrieveLatestVersion retrieves the value for the specified key and latest version from the vault.
func (ov *OwlVault) RetrieveLatestVersion(key string) (map[string]interface{}, error) {
	var data map[string]interface{}

	version, err := ov.storage.LatestVersion(key)
	if err != nil {
		return nil, err
	}

	if version < 1 {
		return nil, fmt.Errorf("key_path not present in the valut")
	}

	// Implement logic to retrieve value from the storage backend
	base64Value, err := ov.storage.Retrieve(key, version)
	if err != nil {
		return nil, err
	}

	// Decode the base64-encoded value
	encryptedValue, err := base64.StdEncoding.DecodeString(base64Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 value: %v", err)
	}

	// Decrypt the retrieved value
	decrypted, err := ov.encryptor.Decrypt(encryptedValue)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(decrypted, &data)

	return data, nil
}

// Additional methods for OwlVault can be added as needed.
