package apprule

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
	if err != nil {
		t.Fatal("failed to cleanup test events")
		return
	}

	defer db.Close()

	for _, event := range events {
		_, err := db.NamedQueryContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE user_id=:user_id", tableName), event)
		if err != nil {
			t.Fatal(err)
		}
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
		{name: "OK", dbOk: true, collectionOk: true, srcStr: dsn},
		{name: "NG", dbOk: false, collectionOk: false, srcStr: "disabled source"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.TODO()
			db, err := getPostgresDb(ctx, c.srcStr)

			if c.dbOk {

				defer db.Close()

				if db == nil || err != nil {
					t.Error("should be able to connect")
					return
				}
			} else if !c.dbOk {
				if err == nil {
					t.Error("should not connect")
				}
				return
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

	t.Run("NG: GetPostgresDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return("disable")

		repo := NewPostgresRepository(c)
		event, err := repo.FindProteinEvent(ctx, userId)
		if event != nil {
			t.Error("event must be nil")
		}
		if err == nil {
			t.Error("error must be not nil")
		}
	})

	t.Run("NG: not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dsn, tableName := getTestPostgresEnv(t)
		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
			// For FindProteinEvent
			mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
			mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		)

		testEvents := makePostgresTestEvents()

		repo := NewPostgresRepository(mock)
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

		cleanupPostgresTestEvents(t, testEvents)
	})

	t.Run("OK", func(t *testing.T) {
		var call *gomock.Call

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		testEvents := makePostgresTestEvents()

		dsn, tableName := getTestPostgresEnv(t)
		mock := config.NewMockConfig(ctrl)
		// For SaveProteinEvent
		call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn)
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName).After(call)
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName).After(call)
		// For FindProteinEvent
		for i := 0; i < len(testEvents); i++ {
			call = mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn).After(call)
			call = mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName).After(call)
		}

		ctx := context.TODO()

		repo := NewPostgresRepository(mock)
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

		cleanupPostgresTestEvents(t, testEvents)
	})
}

func TestPostgresRepository_FindProteinEventByTime(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	setupPostgresTestDb(t)
	defer func() {
		cleanupPostgresTestDb(t)
	}()

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

	dsn, tableName := getTestPostgresEnv(t)
	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		// For FindProteinEventByTime
		mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
	)

	ctx := context.TODO()

	repo := NewPostgresRepository(mock)
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

	cleanupPostgresTestEvents(t, events)
}

func TestPostgresRepository_SaveProteinEvent(t *testing.T) {
	if doesSkipPostgresRepositoryTest(t) {
		return
	}

	setupPostgresTestDb(t)
	defer func() {
		cleanupPostgresTestDb(t)
	}()

	testEvents := makePostgresTestEvents()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dsn, tableName := getTestPostgresEnv(t)
	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("DATABASE_URL"), gomock.Eq("")).Return(dsn),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
		mock.EXPECT().Get(gomock.Eq("POSTGRES_TBL_PROTEINEVENT"), gomock.Eq("")).Return(tableName),
	)

	ctx := context.TODO()

	repo := NewPostgresRepository(mock)
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

	cleanupPostgresTestEvents(t, testEvents)
}
