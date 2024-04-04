package awskms

import (
	"context"
	"fmt"
	"sync"
	"time"

	bigcache "github.com/allegro/bigcache/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// AWSKMSKeyProvider implements the KeyProvider interface for retrieving keys from AWS KMS.
type AWSKMSKeyProvider struct {
	sync.RWMutex

	// AWS KMS specific fields
	region string

	keyId string

	svc *kms.KMS

	isEncKeyExists bool
	encKey         *kms.GenerateDataKeyOutput
	keyCacheStore  *bigcache.BigCache
}

func NewAWSKMSKeyProvider(region string, keyId string) (*AWSKMSKeyProvider, error) {
	// Initialize DynamoDB client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := kms.New(sess)

	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("error encountered while creating kms key provider cache: %v", err)
	}

	return &AWSKMSKeyProvider{
		region:        region,
		keyId:         keyId,
		svc:           svc,
		keyCacheStore: cache,
	}, nil
}

// GenerateKey generates a new encryption key using AWS KMS.
func (kp *AWSKMSKeyProvider) GenerateKey() ([]byte, []byte, []byte, error) {
	var resp *kms.GenerateDataKeyOutput
	var err error

	if kp.isEncKeyExists {
		resp = kp.encKey
	} else {
		// Call AWS KMS API to generate a new data key
		resp, err = kp.svc.GenerateDataKey(&kms.GenerateDataKeyInput{
			KeyId:         aws.String(kp.keyId),
			NumberOfBytes: aws.Int64(64),
		})
		if err != nil {
			return nil, nil, nil, err
		}

		kp.Lock()
		kp.encKey = resp
		kp.isEncKeyExists = true
		kp.Unlock()
	}

	return resp.Plaintext[:32], resp.Plaintext[32:], resp.CiphertextBlob, nil
}

// RetrieveKey retrieves the encryption key from AWS KMS.
func (kp *AWSKMSKeyProvider) RetrieveKey(ctBlob []byte) ([]byte, []byte, error) {
	kp.Lock()
	defer kp.Unlock()

	key, err := kp.keyCacheStore.Get(string(ctBlob))
	if err == nil {
		encKey := key[:32]
		hmacKey := key[32:]

		return encKey, hmacKey, nil
	}

	// Call AWS KMS API to decrypt the encrypted key
	resp, err := kp.svc.Decrypt(&kms.DecryptInput{
		CiphertextBlob: ctBlob,
	})
	if err != nil {
		return nil, nil, err
	}

	key = resp.Plaintext

	err = kp.keyCacheStore.Set(string(ctBlob), key)
	if err != nil {
		return nil, nil, err
	}

	return key[:32], key[32:], nil
}
