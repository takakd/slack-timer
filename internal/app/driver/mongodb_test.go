package driver

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"path/filepath"
	"proteinreminder/internal/pkg/fileutil"
	"runtime"
	"testing"
)

func TestGetMongoDbCollection(t *testing.T) {
	// NOTE: Also use commandline argument
	_, filePath, _, _ := runtime.Caller(0)
	// e.g. internal/configs/.env.test
	envPath := filepath.Join(filepath.Dir(filePath), "../../../configs/.env.test")
	if fileutil.FileExists(envPath) {
		godotenv.Load(envPath)
	}

	cases := []struct {
		name               string
		dbOk, collectionOk bool
		mongoUri           string
	}{
		{name: "ok", dbOk: true, collectionOk: true, mongoUri: os.Getenv("MONGODB_URI")},
		{name: "ng db", dbOk: false, collectionOk: false, mongoUri: "disabled uri"},
	}

	ctx := context.Background()

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

			conn := NewMongoDbConnector()

			db, err := conn.GetDb(ctx, c.mongoUri)
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
