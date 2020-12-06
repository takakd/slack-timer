package repository

import (
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Wrap AWS SDK for Unit test.
type DynamoDbWrapper interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error
	UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

// Use this if the argument of "NewDynamoDb" is not passed.
type DynamoDbWrapperAdapter struct {
	svc *dynamodb.DynamoDB
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return d.svc.GetItem(input)
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return d.svc.Query(input)
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return d.svc.PutItem(input)
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalListOfMaps(l, out)
}

// Dispatch simply.
func (d DynamoDbWrapperAdapter) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}
