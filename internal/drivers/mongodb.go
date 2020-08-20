package drivers

import (
	"net/url"
	"strings"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMongoDb(ctx context.Context) (*mongo.Database, error) {
	config := GetEnvConfig()
	mongoUrl := config.Get("MONGODB_URI")
	elements, err := url.Parse(mongoUrl)
	if err != nil {
		return nil, err
	}

	fmt.Println(elements)

	// e.g. mongodb://.../dbname -> dbname
	dbName := strings.TrimLeft(elements.Path, "/")

	if mongoUrl == "" || dbName == "" {
		panic(fmt.Sprintf("should be set MONGODB_URI. url=%s, name=%s", mongoUrl, dbName))
	}

	// e.g. mongodb://.../dbname -> mongodb//...

	clientOpts := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
		return nil, err
	}

	db := client.Database(dbName)
	return db, nil
}

func GetMongoDbCollection(db *mongo.Database) *mongo.Collection {
	config := GetEnvConfig()
	return db.Collection(config.Get("MONGODB_COLLECTION"))
}

func DisConnectMongoDbClientFunc(ctx context.Context, client *mongo.Client, f func(error)) func() {
	return func() {
		if client != nil {
			if err := client.Disconnect(ctx); err != nil {
				f(err)
			}
		}
		f(nil)
	}
}
