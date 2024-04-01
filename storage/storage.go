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
	Store(key, value string, version int) error

	// Retrieve retrieves the value for the specified key and version.
	Retrieve(key string, version int) (string, error)

	// LatestVersion returns the latest version of the value for the specified key.
	LatestVersion(key string) (int, error)

	Migrate() error // New method for migrations
}

// NewStorage initializes and returns the appropriate storage implementation based on the configuration.
func NewStorage(cfg *config.Config) (Storage, error) {
	switch cfg.Storage.Type {
	case "mysql":
		return mysql.NewMySQLStorage(cfg.Storage.MySQL.ConnectionString)
	case "ddb":
		return ddb.NewDynamoDBStorage(cfg.Storage.DDB.Region, cfg.Storage.DDB.TablePrefix)
	}

	return nil, fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
}
