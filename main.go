package main

import (
	"github.com/ngoyal16/owlvault/encrypt"
	"log"

	"github.com/ngoyal16/owlvault/config"
	"github.com/ngoyal16/owlvault/storage"
	"github.com/ngoyal16/owlvault/vault"
)

func main() {
	// Read configurations
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	// Initialize storage based on configuration
	dbStorage, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Initialize encryptor based on configuration
	encryptor, err := encrypt.NewEncryptor(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize encryptor: %v", err)
	}

	// Initialize OwlVault with the chosen storage implementation
	owlVault := vault.NewOwlVault(dbStorage, encryptor)

	// Example usage: Store and retrieve data
	key := "example_key"
	value := "example_value"
	err = owlVault.StoreKey(key, value)
	if err != nil {
		log.Fatalf("Failed to store data: %v", err)
	}

	storedValue, err := owlVault.RetrieveKey(key, 1)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}
	log.Printf("Retrieved value: %s", storedValue)
}
