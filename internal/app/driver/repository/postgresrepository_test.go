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
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/fileutil"
	"runtime"
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

func makePostgresTestEvents() []*enterpriserule.ProteinEvent {
	return []*enterpriserule.ProteinEvent{
		{
			"id1", time.Now().UTC(), 0,
		},
		{
			"id2", time.Now().UTC(), 0,
		},
	}
}

func cleanupPostgresTestEvents(t *testing.T, events []*enterpriserule.ProteinEvent) {
	if len(events) == 0 {
		return
	}

	ctx := context.TODO()

	srcStr, tableName := getTestPostgresEnv(t)
	db, err := getPostgresDb(ctx, srcStr)
	assert.NoError(t, err)

	defer db.Close()

	for _, event := range events {
		_, err := db.NamedQueryContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id=:user_id", tableName), event)
		assert.NoError(t, err)
	}
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
	envPath := filepath.Join(filepath.Dir(filePath), "../../../configs/.env.test")
	if fileutil.FileExists(envPath) {
		godotenv.Load(envPath)
	}
	dsn = os.Getenv("DATABASE_URL")
	tableName = os.Getenv("POSTGRES_TBL_PROTEINEVENT")
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

func TestPostgresRepository_FindProteinEvent(t *testing.T) {
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
		event, err := repo.FindProteinEvent(ctx, userId)
		assert.Nil(t, event)
		assert.Error(t, err)
	})

	t.Run("ng:not found", func(t *testing.T) {
		testEvents := makePostgresTestEvents()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
			// For FindProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		)

		config.SetConfig(mock)

		repo := NewPostgresRepository()
		_, err := repo.SaveProteinEvent(context.TODO(), testEvents)
		assert.NoError(t, err)

		ctx := context.TODO()
		userId := "nonexistent id"
		event, err := repo.FindProteinEvent(ctx, userId)
		assert.Nil(t, event)
		assert.NoError(t, err)

		cleanupPostgresTestEvents(t, testEvents)
	})

	t.Run("ok", func(t *testing.T) {
		ctx := context.TODO()
		testEvents := makePostgresTestEvents()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		// For SaveProteinEvent
		var call *gomock.Call
		call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn)
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName).After(call)
		// For FindProteinEvent
		for i := 0; i < len(testEvents); i++ {
			call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn).After(call)
			call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName).After(call)
		}
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		_, err := repo.SaveProteinEvent(ctx, testEvents)
		assert.NoError(t, err)

		for _, event := range testEvents {
			got, err := repo.FindProteinEvent(ctx, event.UserId)
			assert.NoError(t, err)
			assert.Equal(t, event, got)
		}

		cleanupPostgresTestEvents(t, testEvents)
	})
}

func TestPostgresRepository_FindProteinEventByTime(t *testing.T) {
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

		from := now.UTC()
		to := from.AddDate(0, 0, 1)
		events[0].UtcTimeToDrink = from
		events[1].UtcTimeToDrink = to
		events[2].UtcTimeToDrink = to.AddDate(0, 0, 1)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
			// For FindProteinEventByTime
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		)
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		_, err := repo.SaveProteinEvent(ctx, events)
		require.NoError(t, err)

		got, err := repo.FindProteinEventByTime(ctx, from, to)
		assert.NoError(t, err)

		wants := events[:2]
		assert.Equal(t, wants, got)

		cleanupPostgresTestEvents(t, events)
	})
}

func TestPostgresRepository_SaveProteinEvent(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	t.Run("ok", func(t *testing.T) {
		setupPostgresTestDb(t)
		defer func() {
			cleanupPostgresTestDb(t)
		}()

		ctx := context.TODO()
		testEvents := makePostgresTestEvents()
		dsn, tableName := getTestPostgresEnv(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		)
		config.SetConfig(mock)

		repo := NewPostgresRepository()
		savedEvents, err := repo.SaveProteinEvent(ctx, testEvents)
		assert.NoError(t, err)
		assert.Equal(t, testEvents, savedEvents)

		cleanupPostgresTestEvents(t, testEvents)
	})
}
