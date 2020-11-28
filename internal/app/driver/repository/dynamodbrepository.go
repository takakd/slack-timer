package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/log"
	"strconv"
	"time"
)

// Wrapper interface to unit test.
type DynamoDbWrapper interface {
	GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
	PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
	UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error
	UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error
	MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error)
}

// Use this if the argument of "NewDynamoDbRepository" is not passed.
type DynamoDbWrapperAdapter struct {
	svc *dynamodb.DynamoDB
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	return d.svc.GetItem(input)
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	return d.svc.Query(input)
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	return d.svc.PutItem(input)
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) UnmarshalMap(m map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalMap(m, out)
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) UnmarshalListOfMaps(l []map[string]*dynamodb.AttributeValue, out interface{}) error {
	return dynamodbattribute.UnmarshalListOfMaps(l, out)
}

// Dispatch simply.
func (d *DynamoDbWrapperAdapter) MarshalMap(in interface{}) (map[string]*dynamodb.AttributeValue, error) {
	return dynamodbattribute.MarshalMap(in)
}

// Implements Repository interface with DynamoDB.
type DynamoDbRepository struct {
	wrp DynamoDbWrapper
}

type DbItemState string

// DAO for repository
type TimerEventDbItem struct {
	UserId           string `dynamodbav:"UserId"`
	NotificationTime int64  `dynamodbav:"NotificationTime"`
	IntervalMin      int    `dynamodbav:"IntervalMin"`
	// Ref. https://forums.aws.amazon.com/thread.jspa?threadID=330244&tstart=0
	State string `dynamodbav:"State"`

	// Not set a value to this field, because this is set by internal for sorting.
	Dummy int `dynamodbav:"Dummy"`
}

func NewTimerEventDbItem(event *enterpriserule.TimerEvent) *TimerEventDbItem {
	t := &TimerEventDbItem{
		UserId:           event.UserId,
		NotificationTime: event.NotificationTime.Unix(),
		IntervalMin:      event.IntervalMin,
		State:            string(event.State),
	}
	return t
}

func (t *TimerEventDbItem) TimerEvent() *enterpriserule.TimerEvent {
	e := &enterpriserule.TimerEvent{
		UserId:      t.UserId,
		IntervalMin: t.IntervalMin,
		State:       enterpriserule.TimerEventState(t.State),
	}
	e.NotificationTime = time.Unix(t.NotificationTime, 0)
	return e
}

// Set wrp to null. In case unit test, set mock interface.
func NewDynamoDbRepository(wrp DynamoDbWrapper) updatetimerevent.Repository {
	if wrp == nil {
		wrp = &DynamoDbWrapperAdapter{
			svc: dynamodb.New(session.New()),
		}
	}
	return &DynamoDbRepository{
		wrp: wrp,
	}
}

// Find timer event by user id.
func (r *DynamoDbRepository) FindTimerEvent(ctx context.Context, userId string) (event *enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userid": {
				S: aws.String(userId),
			},
		},
		KeyConditionExpression: aws.String("UserId = :userid"),
		TableName:              aws.String(config.MustGet("DYNAMODB_TABLE")),
	}
	result, err := r.wrp.Query(input)
	if err != nil {
		return
	}

	itemLen := len(result.Items)
	if itemLen == 0 {
		event = nil
		return
	} else if itemLen > 1 {
		event = nil
		err = fmt.Errorf("item should be one, but found two, user_id=%v", userId)
		return
	}

	var events []TimerEventDbItem
	err = r.wrp.UnmarshalListOfMaps(result.Items, &events)
	if err != nil {
		event = nil
		return
	}
	event = events[0].TimerEvent()
	return
}

// Find timer event from "from" to "to".
func (r *DynamoDbRepository) FindTimerEventByTime(ctx context.Context, from, to time.Time) (events []*enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dummy": {
				N: aws.String(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE")),
			},
			":from": {
				S: aws.String(aws.Time(from).String()),
			},
			":to": {
				S: aws.String(aws.Time(to).String()),
			},
		},
		KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime >= :from AND NotificationTime <= :to"),
		TableName:              aws.String(config.MustGet("DYNAMODB_TABLE")),
		IndexName:              aws.String(config.MustGet("DYNAMODB_INDEX_NAME")),
	}
	result, err := r.wrp.Query(input)
	if err != nil {
		return
	}

	var items []TimerEventDbItem
	err = r.wrp.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		events = nil
		return
	}

	events = make([]*enterpriserule.TimerEvent, len(result.Items))
	for i, v := range items {
		events[i] = v.TimerEvent()
	}
	return
}

// Save TimerEvent to DB.
// Return error and saved event successfully.
func (r *DynamoDbRepository) SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error) {
	dbItem := NewTimerEventDbItem(event)

	dbItem.Dummy, err = strconv.Atoi(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"))
	if err != nil {
		return
	}

	log.Debug(dbItem)
	item, err := r.wrp.MarshalMap(dbItem)
	if err != nil {
		return
	}

	log.Debug(item)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(config.MustGet("DYNAMODB_TABLE")),
	}
	_, err = r.wrp.PutItem(input)
	if err != nil {
		return
	}

	saved = event
	return
}

// Find timer event before eventTime.
func (r *DynamoDbRepository) FindTimerEventsByTime(ctx context.Context, eventTime time.Time) (events []*enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dummy": {
				N: aws.String(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE")),
			},
			":eventTime": {
				N: aws.String(fmt.Sprintf("%d", eventTime.Unix())),
			},
		},
		KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime <= :eventTime"),
		TableName:              aws.String(config.MustGet("DYNAMODB_TABLE")),
		IndexName:              aws.String(config.MustGet("DYNAMODB_INDEX_NAME")),
	}
	result, err := r.wrp.Query(input)
	if err != nil {
		return
	}

	var items []TimerEventDbItem
	err = r.wrp.UnmarshalListOfMaps(result.Items, &items)
	if err != nil {
		events = nil
		return
	}

	events = make([]*enterpriserule.TimerEvent, len(result.Items))
	for i, v := range items {
		events[i] = v.TimerEvent()
	}
	return
}
