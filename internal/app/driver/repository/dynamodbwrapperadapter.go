package repository

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDbWrapperAdapter dispatches to AWS SDK DynamoDB methods.
type DynamoDbWrapperAdapter struct {
	svc *dynamodb.DynamoDB
}

var _ DynamoDbWrapper = (*DynamoDbWrapperAdapter)(nil)

// NewDynamoDbWrapperAdapter create new struct.
func NewDynamoDbWrapperAdapter() *DynamoDbWrapperAdapter {
	return &DynamoDbWrapperAdapter{
		svc: dynamodb.New(session.New()),
	}
}

// GetItem dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return d.svc.GetItem(input)
}

// Query dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return d.svc.Query(input)
}

// PutItem dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return d.svc.PutItem(input)
}

// UnmarshalMap dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}

// UnmarshalListOfMaps dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalListOfMaps(l, out)
}

// MarshalMap dispatches SDK's method simply.
func (d DynamoDbWrapperAdapter) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}
