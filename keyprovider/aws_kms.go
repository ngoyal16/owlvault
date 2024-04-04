package keyprovider

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

// KeyManagementService defines methods for encryption and decryption.
type KeyManagementService interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

// AWSKMSProvider represents an AWS KMS provider for encryption and decryption.
type AWSKMSProvider struct {
	client *kms.Client
}

// NewAWSKMSProviderOption represents an option for configuring the creation of a new AWS KMS provider.
type NewAWSKMSProviderOption func(*aws.Config)

// WithConfig sets the AWS config for the AWS KMS provider.
func WithConfig(cfg aws.Config) NewAWSKMSProviderOption {
	return func(c *aws.Config) {
		*c = cfg
	}
}

// NewAWSKMSProvider creates a new instance of the AWS KMS provider.
func NewAWSKMSProvider(region string, opts ...NewAWSKMSProviderOption) (*AWSKMSProvider, error) {
	var cfg aws.Config
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.Region == "" {
		// Load default AWS config if region is not provided
		loadedCfg, err := awsConfig.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, err
		}
		cfg.Region = loadedCfg.Region
	}

	cfg.Region = region

	client := kms.NewFromConfig(cfg)

	return &AWSKMSProvider{
		client: client,
	}, nil
}

// Encrypt encrypts data using the AWS KMS service.
func (kp *AWSKMSProvider) Encrypt(data []byte) ([]byte, error) {
	input := &kms.EncryptInput{
		KeyId:     aws.String("default"), // Update with your KMS key ID
		Plaintext: data,
	}

	result, err := kp.client.Encrypt(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result.CiphertextBlob, nil
}

// Decrypt decrypts data using the AWS KMS service.
func (kp *AWSKMSProvider) Decrypt(data []byte) ([]byte, error) {
	input := &kms.DecryptInput{
		CiphertextBlob: data,
	}

	result, err := kp.client.Decrypt(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result.Plaintext, nil
}
