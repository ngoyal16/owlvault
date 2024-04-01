package vault

import (
	"fmt"
	"github.com/ngoyal16/owlvault/storage"
)

// OwlVault represents the key vault service.
type OwlVault struct {
	storage storage.Storage
}

// NewOwlVault creates a new instance of OwlVault with the given storage.
func NewOwlVault(storage storage.Storage) *OwlVault {
	return &OwlVault{storage: storage}
}

// StoreKey stores the key-value pair in the vault.
func (ov *OwlVault) StoreKey(key, value string) error {

	version, err := ov.storage.LatestVersion(key)
	if err != nil {
		return err
	}

	version += 1

	// Implement logic to store key-value pair in the storage backend
	if err := ov.storage.Store(key, value, version); err != nil {
		return fmt.Errorf("failed to store key-value pair: %v", err)
	}
	return nil
}

// RetrieveKey retrieves the value for the specified key and version from the vault.
func (ov *OwlVault) RetrieveKey(key string, version int) (string, error) {
	// Implement logic to retrieve value from the storage backend
	retrieve, err := ov.storage.Retrieve(key, version)
	if err != nil {
		return "", err
	}

	return retrieve, nil
}

// Additional methods for OwlVault can be added as needed.
