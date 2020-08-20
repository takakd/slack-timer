package drivers

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	ctx := context.Background()

	db, err := GetMongoDb(ctx)
	if err != nil {
		t.Error(err)
	}
	collection := GetMongoDbCollection(db)

	t.Log("connected to mongodb.")

	filter := bson.M{"user_id": "123"}
	opts := options.Update().SetUpsert(true)
	//value := bson.M{"$set":bson.M{"test1": "abc"}}
	value := bson.M{"$set": bson.M{"test1": "abc"}}
	_, err = collection.UpdateOne(ctx, filter, value, opts)
	if err != nil {
		t.Error(err)
	}
}