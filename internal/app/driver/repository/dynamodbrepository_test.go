package repository

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/pkg/config"
	"testing"
	"time"
)

func TestNewTimerEventDbItem(t *testing.T) {
	caseTime := time.Now()
	caseEvent := &enterpriserule.TimerEvent{
		UserId:           "test_user",
		NotificationTime: caseTime,
		IntervalMin:      3,
	}
	got := NewTimerEventDbItem(caseEvent)
	assert.Equal(t, caseEvent.UserId, got.UserId)
	assert.Equal(t, caseEvent.NotificationTime.Format(time.RFC3339), got.NotificationTime)
	assert.Equal(t, caseEvent.IntervalMin, got.IntervalMin)
}

func TestTimerEventDbItem_TimerEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseTime := time.Now().UTC().Truncate(time.Second)
		want := &enterpriserule.TimerEvent{
			UserId:           "test_user",
			NotificationTime: caseTime,
			IntervalMin:      3,
		}
		got, err := NewTimerEventDbItem(want).TimerEvent()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ng:notification time", func(t *testing.T) {
		caseTime := time.Now().Truncate(time.Second)
		event := &enterpriserule.TimerEvent{
			UserId:           "test_user",
			NotificationTime: caseTime,
			IntervalMin:      3,
		}
		item := NewTimerEventDbItem(event)
		item.NotificationTime = "invalid time format"
		got, err := item.TimerEvent()
		assert.Nil(t, got)
		assert.Error(t, err)
	})
}

func TestNewDynamoDbRepository(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		repo := NewDynamoDbRepository(nil)
		concrete, ok := repo.(*DynamoDbRepository)
		assert.True(t, ok)
		assert.IsType(t, &DynamoDbWrapperAdapter{}, concrete.wrp)
	})

	t.Run("mock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := NewMockDynamoDbWrapper(ctrl)
		repo := NewDynamoDbRepository(mock)
		concrete, ok := repo.(*DynamoDbRepository)
		assert.True(t, ok)
		assert.IsType(t, mock, concrete.wrp)
	})
}

func TestDynamoDbRepository_FindTimerEvent(t *testing.T) {
	t.Run("ng:Query error", func(t *testing.T) {
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Any()).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), "dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:Query returns two items", func(t *testing.T) {
		caseUserId := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserId),
				},
			},
			KeyConditionExpression: aws.String("UserId = :userid"),
			TableName:              aws.String(caseTableName),
		}
		caseItem := &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					"item1": {
						S: aws.String("dummy"),
					},
				},
				{
					"item2": {
						S: aws.String("dummy"),
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), caseUserId)
		assert.Nil(t, got)
		assert.Error(t, err)
	})

	t.Run("ng:UnmarshalMap", func(t *testing.T) {
		caseErr := errors.New("dummy error")
		caseItem := &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					":dummy": {
						S: aws.String("1"),
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Any()).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Items), gomock.Any()).Return(caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), "dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok:Query returns 0 item", func(t *testing.T) {
		caseUserId := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserId),
				},
			},
			KeyConditionExpression: aws.String("UserId = :userid"),
			TableName:              aws.String(caseTableName),
		}
		caseItem := &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), caseUserId)
		assert.Nil(t, got)
		assert.NoError(t, err)
	})

	t.Run("ok: Query returns one item", func(t *testing.T) {
		caseUserId := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserId),
				},
			},
			KeyConditionExpression: aws.String("UserId = :userid"),
			TableName:              aws.String(caseTableName),
		}
		caseItem := &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					":dummy": {
						S: aws.String("1"),
					},
				},
			},
		}
		caseDbItem := &TimerEventDbItem{
			UserId:           caseUserId,
			NotificationTime: time.Now().Format(time.RFC3339),
			Dummy:            1,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Items), gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			events := out.(*[]TimerEventDbItem)
			*events = make([]TimerEventDbItem, 1)
			(*events)[0] = *caseDbItem
			return nil
		})

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), caseUserId)
		assert.NoError(t, err)

		want, err := caseDbItem.TimerEvent()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestDynamoDbRepository_FindTimerEventByTime(t *testing.T) {
	t.Run("ng:Query", func(t *testing.T) {
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return("dummy")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Any()).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEventByTime(context.TODO(), time.Now(), time.Now().Add(100))
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:UnmarshalListOfMaps", func(t *testing.T) {
		caseTableName := "disable"
		caseFrom := time.Now()
		caseTo := time.Now().Add(100)
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":dummy": {
					N: aws.String("1"),
				},
				":from": {
					S: aws.String(caseFrom.String()),
				},
				":to": {
					S: aws.String(caseTo.String()),
				},
			},
			KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime >= :from AND NotificationTime <= :to"),
			TableName:              aws.String(caseTableName),
			IndexName:              aws.String("dummy"),
		}
		caseErr := errors.New("dummy error")
		caseItem := &dynamodb.QueryOutput{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Items), gomock.Any()).Return(caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEventByTime(context.TODO(), caseFrom, caseTo)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseTableName := "disable"
		caseFrom := time.Now()
		caseTo := time.Now().Add(100)
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":dummy": {
					N: aws.String("1"),
				},
				":from": {
					S: aws.String(caseFrom.String()),
				},
				":to": {
					S: aws.String(caseTo.String()),
				},
			},
			KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime >= :from AND NotificationTime <= :to"),
			TableName:              aws.String(caseTableName),
			IndexName:              aws.String("dummy"),
		}
		caseItem := &dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{"dummy": {S: aws.String("dummy")}},
				{"dummy": {S: aws.String("dummy")}},
			},
		}
		caseTime := time.Now().Format(time.RFC3339)
		caseDbItems := []*TimerEventDbItem{
			{
				UserId:           "abc1",
				Dummy:            1,
				NotificationTime: caseTime,
			},
			{
				UserId:           "abc2",
				Dummy:            1,
				NotificationTime: caseTime,
			},
		}

		caseEvents := make([]*enterpriserule.TimerEvent, 2)
		var err error
		caseEvents[0], err = caseDbItems[0].TimerEvent()
		assert.NoError(t, err)
		caseEvents[1], err = caseDbItems[1].TimerEvent()
		assert.NoError(t, err)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Items), gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			events := out.(*[]TimerEventDbItem)
			*events = make([]TimerEventDbItem, 2)
			(*events)[0] = *caseDbItems[0]
			(*events)[1] = *caseDbItems[1]
			return nil
		})

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEventByTime(context.TODO(), caseFrom, caseTo)
		assert.NoError(t, err)
		assert.EqualValues(t, caseEvents, got)
	})
}

func TestDynamoDbRepository_SaveTimerEvent(t *testing.T) {
	t.Run("ng:MarshalMap", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserId:           "test user",
			Dummy:            1,
			NotificationTime: time.Now().Format(time.RFC3339),
		}
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().MarshalMap(gomock.Eq(caseItem)).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		event, err := caseItem.TimerEvent()
		assert.NoError(t, err)
		got, err := repo.SaveTimerEvent(context.TODO(), event)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:PutItem", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserId:           "test user",
			Dummy:            1,
			NotificationTime: time.Now().Format(time.RFC3339),
		}
		caseTableName := "disable"
		caseErr := errors.New("dummy error")
		caseInput := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"from": {
					S: aws.String(""),
				},
			},
			TableName: aws.String(caseTableName),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().MarshalMap(gomock.Eq(caseItem)).Return(caseInput.Item, nil)
		s.EXPECT().PutItem(gomock.Eq(caseInput)).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		event, err := caseItem.TimerEvent()
		got, err := repo.SaveTimerEvent(context.TODO(), event)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserId:           "test user",
			Dummy:            1,
			NotificationTime: time.Now().Format(time.RFC3339),
		}
		caseTableName := "disable"
		caseInput := &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"from": {
					S: aws.String(""),
				},
			},
			TableName: aws.String(caseTableName),
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), gomock.Eq("")).Return("1")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().MarshalMap(gomock.Eq(caseItem)).Return(caseInput.Item, nil)
		s.EXPECT().PutItem(gomock.Eq(caseInput)).Return(nil, nil)

		repo := NewDynamoDbRepository(s)
		event, err := caseItem.TimerEvent()
		assert.NoError(t, err)
		got, err := repo.SaveTimerEvent(context.TODO(), event)
		assert.NoError(t, err)
		assert.EqualValues(t, event, got)
	})
}
