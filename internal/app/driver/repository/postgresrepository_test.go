package repository

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/pkg/config"
	"slacktimer/internal/pkg/fileutil"
	"testing"
	"time"
)

func execTestPostgresSql(t *testing.T, sqlFileName string) {
	t.Logf("exec %s", sqlFileName)

	ctx := context.TODO()

	src, _ := getTestPostgresEnv(t)
	db, err := getPostgresDb(ctx, src)
	require.NoError(t, err)

	defer db.Close()

	sqlByte, err := ioutil.ReadFile(sqlFileName)
	require.NoError(t, err)

	row, err := db.QueryContext(ctx, string(sqlByte))
	require.NoError(t, err)
	require.NotNil(t, row)
}

func setupPostgresTestDb(t *testing.T) {
	execTestPostgresSql(t, "testdata/setup.sql")
}

func cleanupPostgresTestDb(t *testing.T) {
	execTestPostgresSql(t, "testdata/cleanup.sql")
}

func makePostgresTestEvent() *enterpriserule.TimerEvent {
	return &enterpriserule.TimerEvent{
		UserId: "id1", NotificationTime: time.Now().UTC(), IntervalMin: 0,
	}
}

func cleanupPostgresTestEvent(t *testing.T, event *enterpriserule.TimerEvent) {
	if event == nil {
		return
	}

	ctx := context.TODO()

	srcStr, tableName := getTestPostgresEnv(t)
	db, err := getPostgresDb(ctx, srcStr)
	assert.NoError(t, err)

	defer db.Close()

	_, err = db.NamedQueryContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id=:user_id", tableName), event)
	assert.NoError(t, err)
}

func doesSkipPostgresRepositoryTest(t *testing.T) bool {
	src, _ := getTestPostgresEnv(t)
	isSkip := src == "" || src == "skip"
	if isSkip {
		t.Skip("skip Postgres test")
	}
	return isSkip
}

func getTestPostgresEnv(t *testing.T) (dsn, tableName string) {
	// NOTE: Also use commandline argument
	_, filePath, _, _ := runtime.Caller(0)
	// e.g. internal/configs/.env.test
	envPath := filepath.Join(filepath.Dir(filePath), "../../../../configs/.env.test")
	if fileutil.FileExists(envPath) {
		godotenv.Load(envPath)
	}
	dsn = os.Getenv("DATABASE_URL")
	tableName = os.Getenv("POSTGRES_TBL_TIMEREVENT")
	return
}

func Test_getPostgresDb(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	setupPostgresTestDb(t)
	defer func() {
		cleanupPostgresTestDb(t)
	}()

	dsn, _ := getTestPostgresEnv(t)

	cases := []struct {
		name               string
		dbOk, collectionOk bool
		srcStr             string
	}{
		{name: "ok", dbOk: true, collectionOk: true, srcStr: dsn},
		{name: "ng", dbOk: false, collectionOk: false, srcStr: "disabled source"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			db, err := getPostgresDb(ctx, c.srcStr)

			if c.dbOk {

				defer db.Close()

				assert.NotNil(t, db)
				assert.NoError(t, err)
			} else if !c.dbOk {
				assert.Error(t, err)
			}
		})
	}
}

func TestPostgresRepository_FindTimerEvent(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	setupPostgresTestDb(t)
	defer func() {
		cleanupPostgresTestDb(t)
	}()

	t.Run("ng:GetPostgresDb", func(t *testing.T) {
		ctx := context.TODO()
		userId := "abc123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return("disable")
		config.SetConfig(c)

		repo := NewPostgresRepository()
		event, err := repo.FindTimerEvent(ctx, userId)
		assert.Nil(t, event)
		assert.Error(t, err)
	})

	t.Run("ng:not found", func(t *testing.T) {
		testEvent := makePostgresTestEvent()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName),
			// For FindTimerEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName),
		)

		config.SetConfig(mock)

		repo := NewPostgresRepository()
		_, err := repo.SaveTimerEvent(context.TODO(), testEvent)
		assert.NoError(t, err)

		ctx := context.TODO()
		userId := "nonexistent id"
		event, err := repo.FindTimerEvent(ctx, userId)
		assert.Nil(t, event)
		assert.NoError(t, err)

		cleanupPostgresTestEvent(t, testEvent)
	})

	t.Run("ok", func(t *testing.T) {
		ctx := context.TODO()
		testEvent := makePostgresTestEvent()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		// For SaveTimerEvent
		var call *gomock.Call
		call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn)
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName).After(call)
		// For FindTimerEvent
		call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn).After(call)
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName).After(call)
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		_, err := repo.SaveTimerEvent(ctx, testEvent)
		assert.NoError(t, err)

		got, err := repo.FindTimerEvent(ctx, testEvent.UserId)
		assert.NoError(t, err)
		assert.Equal(t, testEvent, got)

		cleanupPostgresTestEvent(t, testEvent)
	})
}

func TestPostgresRepository_FindTimerEventByTime(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}
	t.Run("ok", func(t *testing.T) {
		setupPostgresTestDb(t)
		defer func() {
			cleanupPostgresTestDb(t)
		}()

		ctx := context.TODO()
		dsn, tableName := getTestPostgresEnv(t)

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

		from := now.UTC()
		to := from.AddDate(0, 0, 1)
		events[0].NotificationTime = from
		events[1].NotificationTime = to
		events[2].NotificationTime = to.AddDate(0, 0, 1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName),
			// For FindTimerEventByTime
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName),
		)
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		for _, event := range events {
			_, err := repo.SaveTimerEvent(ctx, event)
			require.NoError(t, err)
		}

		got, err := repo.FindTimerEventByTime(ctx, from, to)
		assert.NoError(t, err)

		wants := events[:2]
		assert.Equal(t, wants, got)

		for _, event := range events {
			cleanupPostgresTestEvent(t, event)
		}
	})
}

func TestPostgresRepository_SaveTimerEvent(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	t.Run("ok", func(t *testing.T) {
		setupPostgresTestDb(t)
		defer func() {
			cleanupPostgresTestDb(t)
		}()

		ctx := context.TODO()
		testEvent := makePostgresTestEvent()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveTimerEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_TIMEREVENT"), gomock.Eq("")).Return(tableName),
		)
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		savedEvent, err := repo.SaveTimerEvent(ctx, testEvent)
		assert.NoError(t, err)
		assert.Equal(t, testEvent, savedEvent)

		cleanupPostgresTestEvent(t, testEvent)
	})
}
