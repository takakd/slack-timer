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
	t.Run("ng:GetItem", func(t *testing.T) {
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().GetItem(gomock.Any()).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), "dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:UnmarshalMap", func(t *testing.T) {
		caseErr := errors.New("dummy error")
		caseItem := &dynamodb.GetItemOutput{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return("dummy")
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().GetItem(gomock.Any()).Return(caseItem, nil)
		s.EXPECT().UnmarshalMap(gomock.Eq(caseItem.Item), gomock.Any()).Return(caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), "dummy")
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseUserId := "abc123"
		caseTableName := "disable"
		caseInput := &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"UserId": {
					S: aws.String(caseUserId),
				},
			},
			TableName: aws.String(caseTableName),
		}
		caseItem := &dynamodb.GetItemOutput{}
		caseEvent := &enterpriserule.TimerEvent{
			UserId: caseUserId,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().GetItem(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalMap(gomock.Eq(caseItem.Item), gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			event := out.(*enterpriserule.TimerEvent)
			event.UserId = caseEvent.UserId
			return nil
		})

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEvent(context.TODO(), caseUserId)
		assert.NoError(t, err)
		assert.EqualValues(t, caseEvent, got)
	})
}

func TestDynamoDbRepository_FindTimerEventByTime(t *testing.T) {
	t.Run("ng:Query", func(t *testing.T) {
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

	t.Run("ng:UnmarshalListOfMaps", func(t *testing.T) {
		caseTableName := "disable"
		caseFrom := time.Now()
		caseTo := time.Now().Add(100)
		caseInput := &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":from": {
					S: aws.String(caseFrom.String()),
				},
				":to": {
					S: aws.String(caseTo.String()),
				},
			},
			KeyConditionExpression: aws.String("NotificationTime >= :from and NotificationTime <= :to"),
			TableName:              aws.String(caseTableName),
		}
		caseErr := errors.New("dummy error")
		caseItem := &dynamodb.GetItemOutput{}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Item), gomock.Any()).Return(caseErr)

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
				":from": {
					S: aws.String(caseFrom.String()),
				},
				":to": {
					S: aws.String(caseTo.String()),
				},
			},
			KeyConditionExpression: aws.String("NotificationTime >= :from and NotificationTime <= :to"),
			TableName:              aws.String(caseTableName),
		}
		caseItem := &dynamodb.GetItemOutput{}
		caseEvent := []*enterpriserule.TimerEvent{
			{
				UserId: "abc1",
			},
			{
				UserId: "abc2",
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().Query(gomock.Eq(caseInput)).Return(caseItem, nil)
		s.EXPECT().UnmarshalListOfMaps(gomock.Eq(caseItem.Item), gomock.Any()).DoAndReturn(func(_, out interface{}) interface{} {
			events := out.([]*enterpriserule.TimerEvent)
			events[0].UserId = caseEvent[0].UserId
			events[1].UserId = caseEvent[1].UserId
			return nil
		})

		repo := NewDynamoDbRepository(s)
		got, err := repo.FindTimerEventByTime(context.TODO(), caseFrom, caseTo)
		assert.NoError(t, err)
		assert.EqualValues(t, caseEvent, got)
	})
}

func TestDynamoDbRepository_SaveTimerEvent(t *testing.T) {
	t.Run("ng:MarshalMap", func(t *testing.T) {
		caseEvent := &enterpriserule.TimerEvent{}
		caseErr := errors.New("dummy error")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().MarshalMap(gomock.Eq(caseEvent)).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.SaveTimerEvent(context.TODO(), caseEvent)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ng:PutItem", func(t *testing.T) {
		caseEvent := &enterpriserule.TimerEvent{}
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
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().PutItem(gomock.Eq(caseInput)).Return(nil, caseErr)

		repo := NewDynamoDbRepository(s)
		got, err := repo.SaveTimerEvent(context.TODO(), caseEvent)
		assert.Nil(t, got)
		assert.Equal(t, caseErr, err)
	})

	t.Run("ok", func(t *testing.T) {
		caseEvent := &enterpriserule.TimerEvent{}
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
		c.EXPECT().Get(gomock.Eq("DYNAMODB_TABLE"), gomock.Eq("")).Return(caseTableName)
		config.SetConfig(c)

		s := NewMockDynamoDbWrapper(ctrl)
		s.EXPECT().PutItem(gomock.Eq(caseInput)).Return(nil, nil)

		repo := NewDynamoDbRepository(s)
		got, err := repo.SaveTimerEvent(context.TODO(), caseEvent)
		assert.NoError(t, err)
		assert.EqualValues(t, caseEvent, got)
	})
}
