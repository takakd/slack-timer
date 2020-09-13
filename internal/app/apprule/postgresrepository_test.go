package apprule

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/pgmock"
	"github.com/jackc/pgx/pgproto3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"net"
	"os"
	"path/filepath"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/config"
	"proteinreminder/internal/pkg/fileutil"
	"proteinreminder/internal/pkg/testutil"
	"runtime"
	"strings"
	"testing"
	"time"
	"fmt"
)

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
	srcStr := getTestPostgreEnv()
	ctx := context.TODO()
	db, err := getPostgreDb(ctx, srcStr)
	if err != nil {
		t.Fatal("failed to cleanup test events")
		return
	}

	for _, event := range events {
		_, err := db.NamedQueryContext(ctx, "DELETE FROM protein_event WHERE user_id=:user_id", event.UserId)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func isSkipPostgresRepositoryTest() bool {
	src := getTestPostgreEnv()
	return src == "" || src == "skip"
}

func getTestPostgreEnv() (sourceStr string) {
	// NOTE: Also use commandline argument
	_, filePath, _, _ := runtime.Caller(0)
	// e.g. internal/configs/.env.test
	envPath := filepath.Join(filepath.Dir(filePath), "../../../configs/.env.test")
	if fileutil.FileExists(envPath) {
		godotenv.Load(envPath)
	}
	sourceStr = os.Getenv("POSTGRES_DATASOURCE")
	return
}

func Test_getPostgreDbEx(t *testing.T) {

	// TEST
	t.Run("NG: GetPostgreDb1", func(t *testing.T) {
		script := &pgmock.Script{
			Steps: pgmock.AcceptUnauthenticatedConnRequestSteps(),
		}
		//script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Query{String: "select 42"}))
		//script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.RowDescription{
		//	Fields: []pgproto3.FieldDescription{
		//		pgproto3.FieldDescription{
		//			Name:                 []byte("?column?"),
		//			TableOID:             0,
		//			TableAttributeNumber: 0,
		//			DataTypeOID:          23,
		//			DataTypeSize:         4,
		//			TypeModifier:         -1,
		//			Format:               0,
		//		},
		//	},
		//}))
		//script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.DataRow{
		//	Values: [][]byte{[]byte("42")},
		//}))
		//script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.CommandComplete{CommandTag: []byte("SELECT 1")}))
		//script.Steps = append(script.Steps, pgmock.SendMessage(&pgproto3.ReadyForQuery{TxStatus: 'I'}))
		script.Steps = append(script.Steps, pgmock.ExpectMessage(&pgproto3.Terminate{}))

		ln, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer ln.Close()

		serverErrChan := make(chan error, 1)
		go func() {
			defer close(serverErrChan)

			conn, err := ln.Accept()
			if err != nil {
				serverErrChan <- err
				return
			}
			defer conn.Close()

			err = conn.SetDeadline(time.Now().Add(time.Second))
			if err != nil {
				serverErrChan <- err
				return
			}

			// TODO
			backend, err := pgproto3.NewBackend(conn, conn)
			require.NoError(t, err)

			t.Log("go func")

			err = script.Run(backend)
			if err != nil {
				serverErrChan <- err
				return
			}
		}()

		parts := strings.Split(ln.Addr().String(), ":")
		host := parts[0]
		port := parts[1]
		connStr := fmt.Sprintf("sslmode=disable host=%s port=%s", host, port)

		t.Log(connStr)

		ctx := context.TODO()
		db, err := getPostgreDb(ctx, connStr)
		fmt.Println(db)
		fmt.Println(err)

		//// Test
		//ctrl := gomock.NewController(t)
		//defer ctrl.Finish()
		//
		//ctx := context.TODO()
		//userId := "abc123"
		//
		//c := config.NewMockConfig(ctrl)
		////c.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return("disable")
		//c.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(connStr)
		//
		//repo := NewPostgresRepository(c)
		//event, err := repo.FindProteinEvent(ctx, userId)
		//if event != nil {
		//	t.Error("event must be nil")
		//}
		//if err == nil {
		//	t.Error("error must be not nil")
		//}
	})

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()
	//pgConn, err := pgconn.Connect(ctx, connStr)
	//require.NoError(t, err)
	//results, err := pgConn.Exec(ctx, "select 42").ReadAll()
	//assert.NoError(t, err)
	//
	//assert.Len(t, results, 1)
	//assert.Nil(t, results[0].Err)
	//assert.Equal(t, "SELECT 1", string(results[0].CommandTag))
	//assert.Len(t, results[0].Rows, 1)
	//assert.Equal(t, "42", string(results[0].Rows[0][0]))
	//
	//pgConn.Close(ctx)
	//
	//assert.NoError(t, <-serverErrChan)

}

func Test_getPostgreDb(t *testing.T) {
	if isSkipPostgresRepositoryTest() {
		t.Skip("skip")
		return
	}

	testSrcStr := getTestPostgreEnv()

	cases := []struct {
		name               string
		dbOk, collectionOk bool
		srcStr             string
	}{
		{name: "OK", dbOk: true, collectionOk: true, srcStr: testSrcStr},
		{name: "NG", dbOk: false, collectionOk: false, srcStr: "disabled source"},
	}

	ctx := context.TODO()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			db, err := getPostgreDb(ctx, c.srcStr)
			if c.dbOk {
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
	if isSkipPostgresRepositoryTest() {
		t.Skip("skip")
		return
	}

	t.Run("NG: GetPostgreDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return("disable")

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
		testEvents := makePostgresTestEvents()
		testSrcStr := getTestPostgreEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)
		gomock.InOrder(
			// For SaveProteinEvent
			mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr),
			// For FindProteinEvent
			mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr),
		)

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
		testEvents := makePostgresTestEvents()
		testSrcStr := getTestPostgreEnv()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mock := config.NewMockConfig(ctrl)

		var call *gomock.Call
		// For SaveProteinEvent
		call = mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr)
		// For FindProteinEvent
		for i := 0; i < len(testEvents); i++ {
			call = mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr).After(call)
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
	if isSkipPostgresRepositoryTest() {
		t.Skip("skip")
		return
	}

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

	testSrcStr := getTestPostgreEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr),
		// For FindProteinEventByTime
		mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr),
	)

	repo := NewPostgresRepository(mock)
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

	cleanupPostgresTestEvents(t, events)
}

func TestPostgresRepository_SaveProteinEvent(t *testing.T) {
	if isSkipPostgresRepositoryTest() {
		t.Skip("skip")
		return
	}

	testEvents := makePostgresTestEvents()

	testSrcStr := getTestPostgreEnv()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := config.NewMockConfig(ctrl)
	gomock.InOrder(
		// For SaveProteinEvent
		mock.EXPECT().Get(gomock.Eq("POSTGRES_DATASOURCE")).Return(testSrcStr),
	)

	repo := NewPostgresRepository(mock)

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

	cleanupPostgresTestEvents(t, testEvents)
}
