package apprule

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/fileutil"
	"proteinreminder/internal/pkg/testutil"
	"runtime"
	"testing"
	"time"
)

func makeTestEvents() []*enterpriserule.ProteinEvent {
	return []*enterpriserule.ProteinEvent{
		{
			"id1", time.Now().UTC(), 0,
		},
		{
			"id2", time.Now().UTC(), 0,
		},
	}
}

func cleanupTestEvents(t *testing.T, events []*enterpriserule.ProteinEvent) {
	if len(events) == 0 {
		return
	}
	url, colName := getTestMongoDbEnv()
	ctx := context.TODO()
	db, err := getMongoDb(ctx, url)
	if err != nil {
		t.Fatal("failed to cleanup test events")
		return
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(e error) {
		return
	})

	col := getMongoCollection(db, colName)

	ids := make([]interface{}, len(events))
	for i, event := range events {
		ids[i] = event.UserId
	}
	filter := bson.D{{
		"user_id",
		bson.D{{
			"$in",
			bson.A(ids),
		}},
	}}
	ret, err := col.DeleteMany(ctx, filter)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
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
	testDbUrl, _ := getTestMongoDbEnv()

	cases := []struct {
		name               string
		dbOk, collectionOk bool
		mongoUri           string
	}{
		{name: "OK", dbOk: true, collectionOk: true, mongoUri: testDbUrl},
		{name: "NG", dbOk: false, collectionOk: false, mongoUri: "disabled uri"},
	}

	ctx := context.TODO()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var client *mongo.Client
			defer func() {
				if client != nil {
					if disCnnErr := client.Disconnect(ctx); disCnnErr != nil {
						t.Error(disCnnErr)
					}
				}
			}()

			db, err := getMongoDb(ctx, c.mongoUri)
			if c.dbOk && err != nil {
				t.Error("should be able to connect")
				return
			} else if !c.dbOk {
				if err == nil {
					t.Error("should not connect")
				}
				return
			}

			client = db.Client()
		})
	}
}

func TestMongoDbRepository_FindProteinEvent(t *testing.T) {
	t.Run("NG: GetDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"
		mongodbUri := "test mongo url"

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(mongodbUri)

		repo := NewMongoDbRepository(c)
		event, err := repo.FindProteinEvent(ctx, userId)
		if event != nil {
			t.Error("event must be nil")
		}
		if err == nil {
			t.Error("error must be not nil")
		}
	})

	t.Run("NG: not found", func(t *testing.T) {
		testEvents := makeTestEvents()
		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol),
			// For FindProteinEvent
			mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl),
			mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol),
		)

		repo := NewMongoDbRepository(mock)
		_, err := repo.SaveProteinEvent(context.TODO(), testEvents)
		if err != nil {
			t.Error(err)
			return
		}

		ctx := context.TODO()
		userId := "nonexistent id"
		event, err := repo.FindProteinEvent(ctx, userId)
		if event != nil {
			t.Error("event must be nil")
		}
		if err != nil {
			t.Error("error must be nil")
		}

		cleanupTestEvents(t, testEvents)
	})

	t.Run("OK", func(t *testing.T) {
		testEvents := makeTestEvents()
		testDbUrl, testDbCol := getTestMongoDbEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)

		var call *gomock.Call
		// For SaveProteinEvent
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl)
		call = mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol).After(call)
		// For FindProteinEvent
		for i := 0; i < len(testEvents); i++ {
			call = mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl).After(call)
			call = mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol).After(call)
		}

		ctx := context.TODO()

		repo := NewMongoDbRepository(mock)
		_, err := repo.SaveProteinEvent(ctx, testEvents)
		if err != nil {
			t.Error(err)
			return
		}

		for _, event := range testEvents {
			got, err := repo.FindProteinEvent(ctx, event.UserId)
			if err != nil {
				t.Error(err)
			}

			if !got.Equal(event) {
				t.Error(testutil.MakeTestMessageWithGotWant(got, event))
			}
		}

		cleanupTestEvents(t, testEvents)
	})
}

func TestMongoDbRepository_FindProteinEventByTime(t *testing.T) {
	now := time.Now().UTC()
	events := []*enterpriserule.ProteinEvent{
		{
			"tid1", now, 2,
		},
		{
			"tid2", now, 2,
		},
		{
			"tid3", now, 2,
		},
	}

	from := time.Now().UTC()
	to := from.AddDate(0, 0, 1)

	events[0].UtcTimeToDrink = from
	events[1].UtcTimeToDrink = to
	events[2].UtcTimeToDrink = to.AddDate(0, 0, 1)

	testDbUrl, testDbCol := getTestMongoDbEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl),
		mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol),
		// For FindProteinEventByTime
		mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl),
		mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol),
	)

	repo := NewMongoDbRepository(mock)
	ctx := context.TODO()
	_, err := repo.SaveProteinEvent(ctx, events)
	if err != nil {
		t.Error(err)
		return
	}

	got, err := repo.FindProteinEventByTime(ctx, from, to)
	if err != nil {
		t.Error(err)
	}
	wants := events[:2]
	if len(got) != 2 || !wants[0].Equal(got[0]) || !wants[1].Equal(got[1]) {
		t.Error(testutil.MakeTestMessageWithGotWant(got[0], wants[0]))
		t.Error(testutil.MakeTestMessageWithGotWant(got[1], wants[1]))
	}

	cleanupTestEvents(t, events)
}

func TestMongoDbRepository_SaveProteinEvent(t *testing.T) {
	testEvents := makeTestEvents()

	testDbUrl, testDbCol := getTestMongoDbEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(testDbUrl),
		mock.EXPECT().Get(gomock.Eq("MONGODB_COLLECTION")).Return(testDbCol),
	)

	repo := NewMongoDbRepository(mock)

	ctx := context.TODO()

	savedEvents, err := repo.SaveProteinEvent(ctx, testEvents)
	if err != nil {
		t.Error(err)
		return
	}

	for i, savedEvent := range savedEvents {
		if !testEvents[i].Equal(savedEvent) {
			t.Error(testutil.MakeTestMessageWithGotWant(savedEvent, testEvents[i]))
		}
	}

	cleanupTestEvents(t, testEvents)
}
