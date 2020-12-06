// package repository provides features that persist data.
package repository

import (
	"context"
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/config"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Implements Repository interface with DynamoDB.
type DynamoDb struct {
	wrp DynamoDbWrapper
}

var _ updatetimerevent.Repository = (*DynamoDb)(nil)
var _ enqueueevent.Repository = (*DynamoDb)(nil)
var _ notifyevent.Repository = (*DynamoDb)(nil)

type DbItemState string

// DAO for repository
type TimerEventDbItem struct {
	UserId           string `dynamodbav:"UserId"`
	NotificationTime string `dynamodbav:"NotificationTime"`
	IntervalMin      int    `dynamodbav:"IntervalMin"`
	// Ref. https://forums.aws.amazon.com/thread.jspa?threadID=330244&tstart=0
	State string `dynamodbav:"State"`

	// Not set a value to this field, because this is set by internal for sorting.
	Dummy int `dynamodbav:"Dummy"`
}

func NewTimerEventDbItem(event *enterpriserule.TimerEvent) *TimerEventDbItem {
	t := &TimerEventDbItem{
		UserId:           event.UserId,
		NotificationTime: event.NotificationTime.Format(time.RFC3339),
		IntervalMin:      event.IntervalMin,
		State:            string(event.State),
	}
	return t
}

func (t TimerEventDbItem) TimerEvent() (*enterpriserule.TimerEvent, error) {
	e := &enterpriserule.TimerEvent{
		UserId:      t.UserId,
		IntervalMin: t.IntervalMin,
		State:       enterpriserule.TimerEventState(t.State),
	}

	var err error
	e.NotificationTime, err = time.Parse(time.RFC3339, t.NotificationTime)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// Set wrp to null. In case unit test, set mock interface.
func NewDynamoDb(wrp DynamoDbWrapper) *DynamoDb {
	if wrp == nil {
		wrp = &DynamoDbWrapperAdapter{
			svc: dynamodb.New(session.New()),
		}
	}
	return &DynamoDb{
		wrp: wrp,
	}
}

// Find timer event by user id.
func (r DynamoDb) FindTimerEvent(ctx context.Context, userId string) (event *enterpriserule.TimerEvent, err error) {
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
	event, err = events[0].TimerEvent()
	return
}

// Find timer event from "from" to "to".
func (r DynamoDb) FindTimerEventByTime(ctx context.Context, from, to time.Time) (events []*enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dummy": {
				N: aws.String(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE")),
			},
			":from": {
				S: aws.String(aws.Time(from).Format(time.RFC3339)),
			},
			":to": {
				S: aws.String(aws.Time(to).Format(time.RFC3339)),
			},
		},
		KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime BETWEEN :from AND :to"),
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
		events[i], err = v.TimerEvent()
		if err != nil {
			events = nil
			break
		}
	}
	return
}

// Save TimerEvent to DB.
// Return error and saved event successfully.
func (r DynamoDb) SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error) {
	dbItem := NewTimerEventDbItem(event)

	dbItem.Dummy, err = strconv.Atoi(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"))
	if err != nil {
		return
	}

	item, err := r.wrp.MarshalMap(dbItem)
	if err != nil {
		return
	}

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
func (r DynamoDb) FindTimerEventsByTime(ctx context.Context, eventTime time.Time) (events []*enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":dummy": {
				N: aws.String(config.MustGet("DYNAMODB_INDEX_PRIMARY_KEY_VALUE")),
			},
			":eventTime": {
				S: aws.String(eventTime.Format(time.RFC3339)),
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
		events[i], err = v.TimerEvent()
		if err != nil {
			events = nil
			break
		}
	}
	return
}
