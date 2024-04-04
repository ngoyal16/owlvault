package keyprovider

import (
	"fmt"
	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/keyprovider/awskms"
	"github.com/ngoyal16/owlvault/keyprovider/localfile"
)

type KeyProvider interface {
	// RetrieveKey retrieves the encryption key.
	RetrieveKey() ([]byte, error)
}

// KeyProviderType represents the type of key provider.
type KeyProviderType string

const (
	// LOCAL_PATH represents the LOCAL_PATH key provider solution.
	LOCAL KeyProviderType = "localfile"
	// AWSKMS represents the AWS KMS key provider solutions.
	AWSKMS KeyProviderType = "awskms"
)

// NewKeyProvider initalizes and returns the appropriate  key provider implementation based on the configuration.
func NewKeyProvider(cfg *config.Config) (KeyProvider, error) {
	var keyProvider KeyProvider
	var err error

	keyProviderType := KeyProviderType(cfg.KeyProvider.Type)

	switch keyProviderType {
	case LOCAL:
		keyProvider, err = localfile.NewLocalFileKeyProvider(cfg.KeyProvider.LocalFile.Path)
	case AWSKMS:
		keyProvider, err = awskms.NewAWSKMSKeyProvider()
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.KeyProvider.Type)
	}

	if err != nil {
		return nil, err
	}

	return keyProvider, nil
}
