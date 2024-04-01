package ddb

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"strconv"
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
func (d *DynamoDBStorage) Store(keyPath string, value string, version int) error {
	kvStoreTableName := d.tablePrefix + "kv_store" // Change to your DynamoDB table name

	// Marshal key-value pair to DynamoDB attribute values
	av, err := dynamodbattribute.MarshalMap(map[string]interface{}{
		"key_path": keyPath,
		"version":  fmt.Sprintf("%019d", version),
		"value":    value,
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
func (d *DynamoDBStorage) Retrieve(keyPath string, version int) (string, error) {
	kvStoreTableName := d.tablePrefix + "kv_store"

	// Create input for GetItem operation
	input := &dynamodb.GetItemInput{
		TableName: aws.String(kvStoreTableName), // Change to your DynamoDB table name
		Key: map[string]*dynamodb.AttributeValue{
			"key_path": {
				S: aws.String(keyPath),
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
	kvStoreTableName := d.tablePrefix + "kv_store"

	// Define input for query
	input := &dynamodb.QueryInput{
		TableName:              aws.String(kvStoreTableName), // Change to your DynamoDB table name
		KeyConditionExpression: aws.String("#key_path = :key_path"),
		ExpressionAttributeNames: map[string]*string{
			"#key_path": aws.String("key_path"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":key_path": {
				S: aws.String(key),
			},
		},
		ScanIndexForward: aws.Bool(false), // Sort results in descending order
		Limit:            aws.Int64(1),    // Limit to 1 item
	}

	// Execute query
	result, err := d.svc.Query(input)
	if err != nil {
		return 0, fmt.Errorf("failed to query DynamoDB: %v", err)
	}

	// Check if any items were returned
	if len(result.Items) == 0 {
		return 0, fmt.Errorf("no versions found for key: %s", key)
	}

	// Unmarshal the version attribute of the first item
	var versionStr string
	if err := dynamodbattribute.Unmarshal(result.Items[0]["version"], &versionStr); err != nil {
		return 0, fmt.Errorf("failed to unmarshal version attribute: %v", err)
	}

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return -1, err
	}

	return version, err

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
		var resourceNotFoundException *dynamodb.ResourceNotFoundException
		if errors.As(err, &resourceNotFoundException) {
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
				AttributeName: aws.String("key_path"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("version"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("key_path"),
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
