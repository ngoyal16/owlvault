package vault

import (
	"crypto/hmac"
	"crypto/sha256"
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
func (ov *OwlVault) Store(keyPath string, data map[string]interface{}) error {
	// Check if version exists
	version, err := ov.storage.LatestVersion(keyPath)
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

	// Calculate HMAC of the value
	hmacValue := ov.generateHMAC(b)

	// Convert encrypted value to base64 encoding
	base64Value := base64.StdEncoding.EncodeToString(encryptedValue)
	base64HMAC := base64.StdEncoding.EncodeToString(hmacValue)

	// Implement logic to store key-value pair in the storage backend
	if err := ov.storage.Store(keyPath, base64Value, base64HMAC, version); err != nil {
		return fmt.Errorf("failed to store key-value pair: %v", err)
	}
	return nil
}

// RetrieveVersion retrieves the value for the specified key and version from the vault.
func (ov *OwlVault) RetrieveVersion(keyPath string, version int) (map[string]interface{}, error) {
	var data map[string]interface{}

	// Implement logic to retrieve value from the storage backend
	base64Value, base64HMAC, err := ov.storage.Retrieve(keyPath, version)
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

	// Calculate HMAC of the decrypted value
	expectedHMAC := ov.generateHMAC(decrypted)

	// Decode the base64-encoded HMAC
	storedHMAC, err := base64.StdEncoding.DecodeString(base64HMAC)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 HMAC: %v", err)
	}

	// Compare the calculated HMAC with the stored HMAC
	if !hmac.Equal(expectedHMAC, storedHMAC) {
		return nil, fmt.Errorf("HMAC validation failed")
	}

	_ = json.Unmarshal(decrypted, &data)

	return data, nil
}

// RetrieveLatestVersion retrieves the value for the specified key and latest version from the vault.
func (ov *OwlVault) RetrieveLatestVersion(keyPath string) (map[string]interface{}, error) {
	version, err := ov.storage.LatestVersion(keyPath)
	if err != nil {
		return nil, err
	}

	if version < 1 {
		return nil, fmt.Errorf("key_path not present in the valut")
	}

	return ov.RetrieveVersion(keyPath, version)
}

// Additional methods for OwlVault can be added as needed.
func (ov *OwlVault) generateHMAC(data []byte) []byte {
	// Calculate HMAC of the decrypted value
	h := hmac.New(sha256.New, []byte("your-hmac-key")) // Replace "your-hmac-key" with your actual HMAC key
	h.Write(data)
	return h.Sum(nil)
}
