package repository

import (
	"errors"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/config"
	"testing"
	"time"

	"slacktimer/internal/app/util/di"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewTimerEventDbItem(t *testing.T) {
	caseTime := time.Now()
	caseEvent := &enterpriserule.TimerEvent{
		UserID:           "test_user",
		NotificationTime: caseTime,
		IntervalMin:      3,
	}
	got := NewTimerEventDbItem(caseEvent)
	assert.Equal(t, caseEvent.UserID, got.UserID)
	assert.Equal(t, caseEvent.NotificationTime.Format(time.RFC3339), got.NotificationTime)
	assert.Equal(t, caseEvent.IntervalMin, got.IntervalMin)
}

func TestTimerEventDbItem_TimerEvent(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		caseTime := time.Now().UTC().Truncate(time.Second)
		want := &enterpriserule.TimerEvent{
			UserID:           "test_user",
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
			UserID:           "test_user",
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

func TestNewDynamoDb(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mw := NewMockDynamoDbWrapper(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		concrete := NewDynamoDb()
		assert.IsType(t, mw, concrete.wrp)
	})

	t.Run("mock", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mw := NewMockDynamoDbWrapper(ctrl)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		concrete := NewDynamoDb()
		assert.IsType(t, mw, concrete.wrp)
	})
}

func TestDynamoDb_FindTimerEvent(t *testing.T) {
	t.Run("ng:Query error", func(t *testing.T) {
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return("dummy")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(gomock.Any()).Return(nil, caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEvent("dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:Query returns two items", func(t *testing.T) {
		caseUserID := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserID),
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
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(caseInput).Return(caseItem, nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEvent(caseUserID)
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
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return("dummy")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(gomock.Any()).Return(caseItem, nil)
		mw.EXPECT().UnmarshalListOfMaps(caseItem.Items, gomock.Any()).Return(caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEvent("dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok:Query returns 0 item", func(t *testing.T) {
		caseUserID := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserID),
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
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(caseInput).Return(caseItem, nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEvent(caseUserID)
		assert.Nil(t, got)
		assert.NoError(t, err)
	})

	t.Run("ok: Query returns one item", func(t *testing.T) {
		caseUserID := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":userid": {
					S: aws.String(caseUserID),
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
			UserID:           caseUserID,
			NotificationTime: time.Now().Format(time.RFC3339),
			Dummy:            1,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(caseInput).Return(caseItem, nil)
		mw.EXPECT().UnmarshalListOfMaps(caseItem.Items, gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			events := out.(*[]TimerEventDbItem)
			*events = make([]TimerEventDbItem, 1)
			(*events)[0] = *caseDbItem
			return nil
		})

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEvent(caseUserID)
		assert.NoError(t, err)

		want, err := caseDbItem.TimerEvent()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func TestDynamoDb_FindTimerEventByTime(t *testing.T) {
	t.Run("ng:Query", func(t *testing.T) {
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return("dummy")
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), "").Return("dummy")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(gomock.Any()).Return(nil, caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEventByTime(time.Now(), time.Now().Add(100))
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
					S: aws.String(caseFrom.Format(time.RFC3339)),
				},
				":to": {
					S: aws.String(caseTo.Format(time.RFC3339)),
				},
			},
			KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime BETWEEN :from AND :to"),
			TableName:              aws.String(caseTableName),
			IndexName:              aws.String("dummy"),
		}
		caseErr := errors.New("dummy error")
		caseItem := &dynamodb.QueryOutput{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), "").Return("dummy")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(caseInput).Return(caseItem, nil)
		mw.EXPECT().UnmarshalListOfMaps(caseItem.Items, gomock.Any()).Return(caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEventByTime(caseFrom, caseTo)
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
					S: aws.String(caseFrom.Format(time.RFC3339)),
				},
				":to": {
					S: aws.String(caseTo.Format(time.RFC3339)),
				},
			},
			KeyConditionExpression: aws.String("Dummy = :dummy AND NotificationTime BETWEEN :from AND :to"),
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
				UserID:           "abc1",
				Dummy:            1,
				NotificationTime: caseTime,
			},
			{
				UserID:           "abc2",
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
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_NAME"), "").Return("dummy")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().Query(caseInput).Return(caseItem, nil)
		mw.EXPECT().UnmarshalListOfMaps(caseItem.Items, gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			events := out.(*[]TimerEventDbItem)
			*events = make([]TimerEventDbItem, 2)
			(*events)[0] = *caseDbItems[0]
			(*events)[1] = *caseDbItems[1]
			return nil
		})

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		got, err := repo.FindTimerEventByTime(caseFrom, caseTo)
		assert.NoError(t, err)
		assert.EqualValues(t, caseEvents, got)
	})
}

func TestDynamoDb_SaveTimerEvent(t *testing.T) {
	t.Run("ng:MarshalMap", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserID:           "test user",
			Dummy:            1,
			NotificationTime: time.Now().Format(time.RFC3339),
		}
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().MarshalMap(caseItem).Return(nil, caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		event, err := caseItem.TimerEvent()
		assert.NoError(t, err)
		got, err := repo.SaveTimerEvent(event)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:PutItem", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserID:           "test user",
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
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().MarshalMap(caseItem).Return(caseInput.Item, nil)
		mw.EXPECT().PutItem(caseInput).Return(nil, caseErr)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		event, err := caseItem.TimerEvent()
		got, err := repo.SaveTimerEvent(event)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseItem := &TimerEventDbItem{
			UserID:           "test user",
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
		c.EXPECT().Get(gomock.Eq("DYNAMODB_INDEX_PRIMARY_KEY_VALUE"), "").Return("1")
		c.EXPECT().Get("DYNAMODB_TABLE", "").Return(caseTableName)
		config.SetConfig(c)

		mw := NewMockDynamoDbWrapper(ctrl)
		mw.EXPECT().MarshalMap(caseItem).Return(caseInput.Item, nil)
		mw.EXPECT().PutItem(caseInput).Return(nil, nil)

		md := di.NewMockDI(ctrl)
		md.EXPECT().Get("repository.DynamoDbWrapper").Return(mw)
		di.SetDi(md)

		repo := NewDynamoDb()
		event, err := caseItem.TimerEvent()
		assert.NoError(t, err)
		got, err := repo.SaveTimerEvent(event)
		assert.NoError(t, err)
		assert.EqualValues(t, event, got)
	})
}
