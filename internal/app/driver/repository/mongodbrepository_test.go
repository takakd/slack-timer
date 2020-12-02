package repository

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"runtime"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/util/config"
	"slacktimer/internal/pkg/fileutil"
	"testing"
	"time"
)

func makeTestEvent() *enterpriserule.TimerEvent {
	return &enterpriserule.TimerEvent{
		UserId:           "id1",
		NotificationTime: time.Now().UTC(),
		IntervalMin:      0,
	}
}

func cleanupTestEvents(t *testing.T, event *enterpriserule.TimerEvent) {
	if event == nil {
		return
	}
	url, colName := getTestMongoDbEnv()
	ctx := context.TODO()
	db, err := getMongoDb(ctx, url)
	require.NoError(t, err)
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(e error) {
		return
	})

	col := getMongoCollection(db, colName)

	filter := bson.D{{
		"user_id",
		bson.D{{
			"$in",
			event.UserId,
		}},
	}}
	_, err = col.DeleteMany(ctx, filter)
	assert.NoError(t, err)
}

func doesSkipMongoDbRepositoryTest(t *testing.T) bool {
	url, _ := getTestMongoDbEnv()
	isSkip := url == "" || url == "skip"
	if isSkip {
		t.Skip("skip MongoDB test")
	}
	return isSkip
}

func getTestMongoDbEnv() (url, collection string) {
	// NOTE: Also use commandline argument
	_, filePath, _, _ := runtime.Caller(0)
	// e.g. internal/configs/.env.test
	envPath := filepath.Join(filepath.Dir(filePath), "../../../configs/.env.test")
	if fileutil.FileExists(envPath) {
		godotenv.Load(envPath)
	}
	url = os.Getenv("MONGODB_URI")
	collection = os.Getenv("MONGODB_COLLECTION")
	return
}

// Test getMongoDb, getMongoCollection, and disconnectMongoDbClientFunc together.
func Test_getMongoDb(t *testing.T) {
	if doesSkipMongoDbRepositoryTest(t) {
		return
	}

	testDbUrl, _ := getTestMongoDbEnv()

	cases := []struct {
		name               string
		dbOk, collectionOk bool
		mongoUri           string
	}{
		{name: "ok", dbOk: true, collectionOk: true, mongoUri: testDbUrl},
		{name: "ng", dbOk: false, collectionOk: false, mongoUri: "disabled uri"},
	}

	ctx := context.TODO()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var client *mongo.Client
			defer func() {
				if client != nil {
					assert.NoError(t, client.Disconnect(ctx))
				}
			}()

			db, err := getMongoDb(ctx, c.mongoUri)
			if c.dbOk {
				assert.NoError(t, err)
			} else if !c.dbOk {
				assert.Error(t, err)
			}

			client = db.Client()
		})
	}
}

func TestMongoDbRepository_FindTimerEvent(t *testing.T) {
	if doesSkipMongoDbRepositoryTest(t) {
		return
	}

	t.Run("ng:GetDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"
		mongodbUri := "test mongo url"

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(mongodbUri)
		config.SetConfig(c)

		repo := NewMongoDbRepository()
		event, err := repo.FindTimerEvent(ctx, userId)
		assert.NotNil(t, event)
		assert.NoError(t, err)
	})

	t.Run("ng:not found", func(t *testing.T) {
		testEvent := makeTestEvent()
		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol),
			// For FindTimerEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol),
		)
		config.SetConfig(mock)

		repo := NewMongoDbRepository()
		_, err := repo.SaveTimerEvent(context.TODO(), testEvent)
		assert.NoError(t, err)

		ctx := context.TODO()
		userId := "nonexistent id"
		event, err := repo.FindTimerEvent(ctx, userId)
		assert.NotNil(t, event)
		assert.NoError(t, err)

		cleanupTestEvents(t, testEvent)
	})

	t.Run("ok", func(t *testing.T) {
		testEvent := makeTestEvent()
		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)

		var call *gomock.Call
		// For SaveTimerEvent
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl)
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol).After(call)
		// For FindTimerEvent
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl).After(call)
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol).After(call)

		config.SetConfig(mock)

		ctx := context.TODO()

		repo := NewMongoDbRepository()
		_, err := repo.SaveTimerEvent(ctx, testEvent)
		assert.NoError(t, err)

		got, err := repo.FindTimerEvent(ctx, testEvent.UserId)
		assert.NoError(t, err)
		assert.Equal(t, testEvent, got)

		cleanupTestEvents(t, testEvent)
	})
}

func TestMongoDbRepository_FindTimerEventByTime(t *testing.T) {
	if doesSkipMongoDbRepositoryTest(t) {
		return
	}

	t.Run("ok", func(t *testing.T) {
		now := time.Now().UTC()
		events := []*enterpriserule.TimerEvent{
			{
				UserId: "tid1", NotificationTime: now, IntervalMin: 2,
			},
			{
				UserId: "tid2", NotificationTime: now, IntervalMin: 2,
			},
			{
				UserId: "tid3", NotificationTime: now, IntervalMin: 2,
			},
		}

		from := time.Now().UTC()
		to := from.AddDate(0, 0, 1)

		events[0].NotificationTime = from
		events[1].NotificationTime = to
		events[2].NotificationTime = to.AddDate(0, 0, 1)

		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol),
			// For FindTimerEventByTime
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol),
		)

		config.SetConfig(mock)

		repo := NewMongoDbRepository()
		ctx := context.TODO()
		for _, event := range events {
			_, err := repo.SaveTimerEvent(ctx, event)
			assert.NoError(t, err)
		}

		got, err := repo.FindTimerEventByTime(ctx, from, to)
		assert.NoError(t, err)

		wants := events[:2]
		// TODO: check
		assert.Equal(t, wants, got)
		//if len(got) != 2 || !wants[0].Equal(got[0]) || !wants[1].Equal(got[1]) {
		//	t.Error(testutil.MakeTestMessageWithGotWant(got[0], wants[0]))
		//	t.Error(testutil.MakeTestMessageWithGotWant(got[1], wants[1]))
		//}

		for _, event := range events {
			cleanupTestEvents(t, event)
		}
	})
}

func TestMongoDbRepository_SaveTimerEvent(t *testing.T) {
	if doesSkipMongoDbRepositoryTest(t) {
		return
	}

	t.Run("ok", func(t *testing.T) {
		testEvent := makeTestEvent()

		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI"), gomock.Eq("")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION"), gomock.Eq("")).Return(testDbCol),
		)

		config.SetConfig(mock)

		repo := NewMongoDbRepository()

		ctx := context.TODO()

		savedEvent, err := repo.SaveTimerEvent(ctx, testEvent)
		assert.NoError(t, err)
		assert.Equal(t, testEvent, savedEvent)
		cleanupTestEvents(t, testEvent)
	})
}
