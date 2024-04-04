package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents the configuration for the OwlVault service.
type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
	} `yaml:"server"`
	Encryptor struct {
		Type string `yaml:"type"`
	} `yaml:"encryptor"`
	KeyProvider struct {
		Type      string `yaml:"type"`
		LocalFile struct {
			Path string `yaml:"path"`
		} `yaml:"local_file"`
		AWSKMS struct {
			KeyArn string `yaml:"key_arn"`
		} `yaml:"aws_kms"`
	} `yaml:"key_provider"`
	Storage struct {
		Type  string `yaml:"type"`
		MySQL struct {
			ConnectionString string `yaml:"connection_string"`
		} `yaml:"mysql"`
		DDB struct {
			Region      string `yaml:"region"`
			TablePrefix string `yaml:"table_prefix"`
		} `yaml:"dynamodb"`
		// Add other storage types here
	} `yaml:"storage"`
}

// ReadConfig reads configuration from the specified YAML file path provided by the environment variable.
func ReadConfig() (*Config, error) {
	// Get the config file path from the environment variable
	configPath := os.Getenv("OWLVAULT_CONFIG_PATH")
	if configPath == "" {
		return nil, fmt.Errorf("environment variable OWLVAULT_CONFIG_PATH is not set")
	}

	// Get the absolute path of the configuration file
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	// Check if the file exists
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist at path: %s", absPath)
	}

	// Read YAML configuration file
	configFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML into Config struct
	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
