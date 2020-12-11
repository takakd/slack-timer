// Package repository provides features that persist event.
package repository

import (
	"fmt"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/enqueueevent"
	"slacktimer/internal/app/usecase/notifyevent"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/app/util/config"
	"strconv"
	"time"

	"slacktimer/internal/app/util/di"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// DynamoDb implements Repository interface with DynamoDB.
type DynamoDb struct {
	wrp DynamoDbWrapper
}

var _ updatetimerevent.Repository = (*DynamoDb)(nil)
var _ enqueueevent.Repository = (*DynamoDb)(nil)
var _ notifyevent.Repository = (*DynamoDb)(nil)

// DbItemState represents the type of Queueing state.
type DbItemState string

// TimerEventDbItem s DAO for repository.
type TimerEventDbItem struct {
	UserID           string `dynamodbav:"UserId"`
	Text             string `dynamodbav:"Text"`
	NotificationTime string `dynamodbav:"NotificationTime"`
	IntervalMin      int    `dynamodbav:"IntervalMin"`
	// Ref. https://forums.aws.amazon.com/thread.jspa?threadID=330244&tstart=0
	State string `dynamodbav:"State"`

	// Not set a value to this field, because this is set by internal for sorting.
	Dummy int `dynamodbav:"Dummy"`
}

// NewTimerEventDbItem creates new struct.
func NewTimerEventDbItem(event *enterpriserule.TimerEvent) *TimerEventDbItem {
	t := &TimerEventDbItem{
		UserID:           event.UserID(),
		Text:             event.Text,
		NotificationTime: event.NotificationTime.Format(time.RFC3339),
		IntervalMin:      event.IntervalMin,
		State:            string(event.State),
	}
	return t
}

// TimerEvent generates enterpriserule.TimerEvent struct.
func (t TimerEventDbItem) TimerEvent() (*enterpriserule.TimerEvent, error) {
	e, err := enterpriserule.NewTimerEvent(t.UserID)
	if err != nil {
		return nil, err
	}
	e.IntervalMin = t.IntervalMin
	e.State = enterpriserule.TimerEventState(t.State)
	e.Text = t.Text

	e.NotificationTime, err = time.Parse(time.RFC3339, t.NotificationTime)
	if err != nil {
		return nil, err
	}

	return e, nil
}

// NewDynamoDb creates new struct.
func NewDynamoDb() *DynamoDb {
	return &DynamoDb{
		wrp: di.Get("repository.DynamoDbWrapper").(DynamoDbWrapper),
	}
}

// FindTimerEvent finds an event by user id.
func (r DynamoDb) FindTimerEvent(userID string) (event *enterpriserule.TimerEvent, err error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":userid": {
				S: aws.String(userID),
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
		err = fmt.Errorf("item should be one, but found two, user_id=%v", userID)
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

// FindTimerEventByTime finds events from "from" to "to".
func (r DynamoDb) FindTimerEventByTime(from, to time.Time) (events []*enterpriserule.TimerEvent, err error) {
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

// SaveTimerEvent persists an event.
// Return nil error and saved event if it is successful, if not, return an error.
func (r DynamoDb) SaveTimerEvent(event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error) {
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

// FindTimerEventsByTime finds events before eventTime.
func (r DynamoDb) FindTimerEventsByTime(eventTime time.Time) (events []*enterpriserule.TimerEvent, err error) {
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
