package apprule

import (
	"context"
	"testing"
	"github.com/golang/mock/gomock"
	"proteinreminder/internal/app/driver"
	"github.com/pkg/errors"
	"proteinreminder/internal/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

//func testEvents() []*enterpriserule.ProteinEvent {
//	return []*enterpriserule.ProteinEvent{
//		{
//			"id1", time.Now().UTC(), 0,
//		},
//		{
//			"id2", time.Now().UTC(), 0,
//		},
//	}
//}

func TestFindProteinEvent(t *testing.T) {
	//testEvents := testEvents()
	//
	//repo := NewMongoDbRepository()
	//_, err := repo.SaveProteinEvent(context.TODO(), testEvents)
	//if err != nil {
	//	t.Error(err)
	//	return
	//}

	t.Run("NG: GetDb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"
		mongodbUri := "test mongo url"
		wantErr := errors.New("test: fail")

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(mongodbUri)

		m := driver.NewMockMongoDbConnector(ctrl)
		m.EXPECT().GetDb(gomock.Eq(ctx), gomock.Eq(mongodbUri)).Return(nil, wantErr)
		m.EXPECT().DisConnectClientFunc(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(0)
		m.EXPECT().GetCollection(gomock.Any(), gomock.Any()).MaxTimes(0)

		repo := NewMongoDbRepository(m, c)
		event, err := repo.FindProteinEvent(ctx, userId)
		if event != nil {
			t.Error("event must be nil")
		}
		if err != wantErr {
			t.Error("error must be not nil")
		}
	})

	t.Run("OK: not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.TODO()
		userId := "abc123"
		mongodbUri := "test mongo url"
		mongoCol := "test mongo col"

		c := config.NewMockConfig(ctrl)
		c.EXPECT().Get(gomock.Eq("MONGODB_URI")).Return(mongodbUri)

		client := driver.NewMockMongoClient(ctrl)
		db := driver.NewMockMongoDatabase(ctrl)
		db.EXPECT().Client().Return(client).MaxTimes(1)

		mRet := driver.NewMockMongoSingleResult(ctrl)
		mRet.EXPECT().Err().Return(mongo.ErrNoDocuments)

		col := driver.NewMockMongoCollection(ctrl)
		col.EXPECT().FindOne(gomock.Eq(ctx), gomock.Eq(bson.M{"uesr_id":userId})).Return(mRet)

		conn := driver.NewMockMongoDbConnector(ctrl)
		conn.EXPECT().GetDb(gomock.Any(), gomock.Any()).Return(db, nil)
		conn.EXPECT().DisConnectClientFunc(gomock.Eq(ctx), gomock.Eq(client), gomock.Any()).MaxTimes(1)
		conn.EXPECT().GetCollection(gomock.Eq(db), gomock.Eq(mongoCol)).Return(col)

		repo := NewMongoDbRepository(conn, c)
		event, err := repo.FindProteinEvent(ctx, userId)
		if event != nil {
			t.Error("event must be nil")
		}
		if err != nil {
			t.Error("error must be nil")
		}
	})



	//
	//// collection find one error
	//// result.DecodeError
	//
	//t.Run("OK", func(t *testing.T) {
	//	for _, event := range testEvents {
	//		repo := NewMongoDbRepository()
	//		got, err := repo.FindProteinEvent(context.TODO(), event.UserId)
	//		if err != nil {
	//			t.Error(err)
	//		}
	//
	//		if !got.Equal(event) {
	//			t.Error(testutil.MakeTestMessageWithGotWant(got, event))
	//		}
	//	}
	//})
	//
	//t.Run("NG", func(t *testing.T) {
	//	repo := NewMongoDbRepository()
	//	event, err := repo.FindProteinEvent(context.TODO(), "disable id")
	//	if event != nil || err != nil {
	//		t.Error("should not exist")
	//	}
	//})
}

//func TestFindProteinEventByTime(t *testing.T) {
//	now := time.Now().UTC()
//	events := []*enterpriserule.ProteinEvent{
//		{
//			"id1", now, 2,
//		},
//		{
//			"id2", now, 2,
//		},
//		{
//			"id3", now, 2,
//		},
//	}
//
//	from := time.Now().UTC()
//	to := from.AddDate(0, 0, 1)
//
//	events[0].UtcTimeToDrink = from
//	events[1].UtcTimeToDrink = to
//	events[2].UtcTimeToDrink = to.AddDate(0, 0, 1)
//
//	repo := NewMongoDbRepository()
//
//	repo.SaveProteinEvent(context.TODO(), events)
//
//	got, err := repo.FindProteinEventByTime(context.TODO(), from, to)
//	if err != nil {
//		t.Error(err)
//	}
//	wants := events[:2]
//	if len(got) != 2 || !wants[0].Equal(got[0]) || !wants[1].Equal(got[1]) {
//		t.Error(testutil.MakeTestMessageWithGotWant(got[0], wants[0]))
//		t.Error(testutil.MakeTestMessageWithGotWant(got[1], wants[1]))
//	}
//}
//
//func TestSaveProteinEvent(t *testing.T) {
//	testEvents := testEvents()
//	repo := NewMongoDbRepository()
//	savedEvents, err := repo.SaveProteinEvent(context.TODO(), testEvents)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//
//	for i, savedEvent := range savedEvents {
//		if !testEvents[i].Equal(savedEvent) {
//			t.Error(testutil.MakeTestMessageWithGotWant(savedEvent, testEvents[i]))
//		}
//	}
//}
