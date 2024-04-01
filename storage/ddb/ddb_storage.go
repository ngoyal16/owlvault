package ddb

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDBStorage implements the Storage interface for DynamoDB.
type DynamoDBStorage struct {
	svc *dynamodb.DynamoDB

	tablePrefix string
}

// NewDynamoDBStorage creates a new instance of DynamoDBStorage.
func NewDynamoDBStorage(region string, tablePrefix string) (*DynamoDBStorage, error) {
	// Initialize DynamoDB client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}
	svc := dynamodb.New(sess)

	return &DynamoDBStorage{
		svc:         svc,
		tablePrefix: tablePrefix,
	}, nil
}

// Store stores the key-value pair with the specified version.
func (d *DynamoDBStorage) Store(key, value string, version int) error {
	kvStoreTableName := d.tablePrefix + "kv_store" // Change to your DynamoDB table name

	// Marshal key-value pair to DynamoDB attribute values
	av, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		"key":     key,
		"version": fmt.Sprintf("%019d", version),
		"value":   value,
	})
	if err != nil {
		return err
	}

	// Create input for PutItem operation
	input := &dynamodb.PutItemInput{
		TableName: aws.String(kvStoreTableName), // Change to your DynamoDB table name
		Item:      av,
	}

	// Execute PutItem operation
	_, err = d.svc.PutItem(input)
	if err != nil {
		return fmt.Errorf("failed to store item: %v", err)
	}
	return nil
}

// Retrieve retrieves the value for the specified key and version.
func (d *DynamoDBStorage) Retrieve(key string, version int) (string, error) {
	kvStoreTableName := d.tablePrefix + "kv_store"

	// Create input for GetItem operation
	input := &dynamodb.GetItemInput{
		TableName: aws.String(kvStoreTableName), // Change to your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"key": {
				S: aws.String(key),
			},
			"version": {
				S: aws.String(fmt.Sprintf("%019d", version)),
			},
		},
	}

	// Execute GetItem operation
	result, err := d.svc.GetItem(input)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve item: %v", err)
	}

	// Unmarshal retrieved item
	item := struct {
		Value string `json:"value"`
	}{}
	if err := dynamodbattribute.UnmarshalMap(result.Item, &item); err != nil {
		return "", fmt.Errorf("failed to unmarshal item: %v", err)
	}

	return item.Value, nil
}

// LatestVersion is not applicable for DynamoDB storage
func (d *DynamoDBStorage) LatestVersion(key string) (int, error) {
	return 0, fmt.Errorf("LatestVersion method is not applicable for DynamoDB storage")
}

// Migrate is not applicable for DynamoDB storage
func (d *DynamoDBStorage) Migrate() error {
	// Check if the table exists
	kvStoreTableName := d.tablePrefix + "kv_store" // Change to your DynamoDB table name
	exists, err := d.checkTableExists(kvStoreTableName)
	if err != nil {
		return fmt.Errorf("failed to check if table exists: %v", err)
	}
	if !exists {
		// Table does not exist, create it
		if err := d.createKVStoreTable(kvStoreTableName); err != nil {
			return fmt.Errorf("failed to create table: %v", err)
		}
		fmt.Printf("Table '%s' created successfully\n", kvStoreTableName)
	} else {
		fmt.Printf("Table '%s' already exists\n", kvStoreTableName)
	}

	return nil
}

// checkTableExists checks if the table exists in DynamoDB.
func (d *DynamoDBStorage) checkTableExists(tableName string) (bool, error) {
	// Describe table to check if it exists
	_, err := d.svc.DescribeTable(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
			return false, nil // Table does not exist
		}
		return false, err // Other error occurred
	}
	return true, nil // Table exists
}

func (d *DynamoDBStorage) createKVStoreTable(tableName string) error {
	// Define table schema
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("key"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("version"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5), // Adjust as needed
			WriteCapacityUnits: aws.Int64(5), // Adjust as needed
		},
		TableName: aws.String(tableName),
	}

	// Create table
	_, err := d.svc.CreateTable(input)
	return err
}
