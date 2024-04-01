package storage

import (
	"fmt"
	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/storage/ddb"
	"github.com/ngoyal16/owlvault/storage/mysql"
)

// Storage defines the interface for interacting with the storage backend.
type Storage interface {
	// Store stores the key-value pair with the specified version and timestamp.
	Store(keyPath string, contents string, hmac string, version int) error

	// Retrieve retrieves the value for the specified key and version.
	Retrieve(keyPath string, version int) (string, string, error)

	// LatestVersion returns the latest version of the value for the specified key.
	LatestVersion(keyPath string) (int, error)

	Migrate() error // New method for migrations
}

// StorageType represents the type of storage.
type StorageType string

const (
	// MYSQL represents the MySQL storage solution.
	MYSQL StorageType = "mysql"
	// DDB represents the DynamoDB storage solution.
	DDB StorageType = "dynamodb"
	// Add more storage solution as needed
)

// NewStorage initializes and returns the appropriate storage implementation based on the configuration.
func NewStorage(cfg *config.Config) (Storage, error) {
	var dbStorage Storage
	var err error

	storageType := StorageType(cfg.Storage.Type)

	switch storageType {
	case MYSQL:
		dbStorage, err = mysql.NewMySQLStorage(cfg.Storage.MySQL.ConnectionString)
	case DDB:
		dbStorage, err = ddb.NewDynamoDBStorage(cfg.Storage.DDB.Region, cfg.Storage.DDB.TablePrefix)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
	}

	if err != nil {
		return nil, err
	}

	err = dbStorage.Migrate()
	if err != nil {
		return nil, err
	}

	return dbStorage, nil
}
