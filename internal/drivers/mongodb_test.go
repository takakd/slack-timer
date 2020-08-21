package drivers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"testing"
)

func TestGetMongoDbCollection(t *testing.T) {
	// NOTE: Not need, cause using commandline mode.
	//_, filePath, _, _ := runtime.Caller(0)
	//
	//envPath := filepath.Join(filepath.Dir(filePath), "../../configs/.env")
	//if fileutil.FileExists(envPath) {
	//	godotenv.Load(envPath)
	//}

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

			os.Setenv("MONGODB_URI", c.mongoUri)

			db, err := GetMongoDb(ctx)
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
