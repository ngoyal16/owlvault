package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // Import MySQL driver
)

// MySQLStorage implements the Storage interface for MySQL database.
type MySQLStorage struct {
	db *sql.DB
	// You can include additional fields if needed
}

// NewMySQLStorage creates a new instance of MySQLStorage.
func NewMySQLStorage(connectionString string) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &MySQLStorage{db: db}, nil
}

// Migrate runs the necessary database migrations for MySQLStorage.
func (m *MySQLStorage) Migrate() error {
	// Implement database migrations specific to MySQL
	_, err := m.db.Exec(`
        CREATE TABLE IF NOT EXISTS key_value_pairs (
            id INT AUTO_INCREMENT PRIMARY KEY,
            key VARCHAR(255) NOT NULL,
            value TEXT NOT NULL,
            version INT NOT NULL
        );
    `)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}
	return nil
}

// Store stores the key-value pair with the specified version and timestamp.
func (m *MySQLStorage) Store(key, value string, version int) error {
	_, err := m.db.Exec("INSERT INTO data (key, value, version) VALUES (?, ?, ?, ?)", key, value, version)
	return err
}

// Retrieve retrieves the value for the specified key and version.
func (m *MySQLStorage) Retrieve(key string, version int) (string, error) {
	var value string
	err := m.db.QueryRow("SELECT value FROM data WHERE key = ? AND version = ?", key, version).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

// LatestVersion returns the latest version of the value for the specified key.
func (m *MySQLStorage) LatestVersion(key string) (int, error) {
	var latestVersion int
	err := m.db.QueryRow("SELECT MAX(version) FROM data WHERE key = ?", key).Scan(&latestVersion)
	if err != nil {
		return 0, err
	}
	return latestVersion, nil
}
