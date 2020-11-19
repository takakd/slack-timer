package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"slacktimer/internal/app/enterpriserule"
	"slacktimer/internal/app/usecase/updatetimerevent"
	"slacktimer/internal/pkg/config"
	"strings"
	"time"
)

// Implements Repository interface with MongoDB.
type MongoDbRepository struct {
}

// Get mongo.Database
func getMongoDb(ctx context.Context, mongoUrl string) (db *mongo.Database, err error) {
	elements, err := url.Parse(mongoUrl)
	if err != nil {
		return
	}

	// e.g. mongodb://.../dbname -> dbname
	dbName := strings.TrimLeft(elements.Path, "/")

	if mongoUrl == "" || dbName == "" {
		panic(fmt.Sprintf("should be set MONGODB_URI. url=%s, name=%s", mongoUrl, dbName))
	}

	// e.g. mongodb://.../dbname -> mongodb//...

	clientOpts := options.Client().ApplyURI(mongoUrl)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println(err)
		return
	}

	db = client.Database(dbName)
	return
}

// Get mongo.Collection
func getMongoCollection(db *mongo.Database, name string) (col *mongo.Collection) {
	col = db.Collection(name)
	return
}

/*
	Disconnect MongoDB Client Function. Use this with defer.

	e.g.
	defer disconnectMongoDbClientFunc(ctx, client, func(err error){
		// something you like
	})()
*/
func disconnectMongoDbClientFunc(ctx context.Context, client *mongo.Client, f func(error)) (fnc func()) {
	fnc = func() {
		if client != nil {
			if err := client.Disconnect(ctx); err != nil {
				f(err)
			}
		}
		f(nil)
	}
	return
}

func NewMongoDbRepository() updatetimerevent.Repository {
	return &MongoDbRepository{}
}

// Find timer event by user id.
func (r *MongoDbRepository) FindTimerEvent(ctx context.Context, userId string) (event *enterpriserule.TimerEvent, err error) {
	db, err := getMongoDb(ctx, config.Get("MONGODB_URI", ""))
	if err != nil {
		return
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, config.Get("MONGODB_COLLECTION", ""))

	var value enterpriserule.TimerEvent
	filter := bson.M{"user_id": userId}
	result := collection.FindOne(ctx, filter)
	notFound := result.Err() == mongo.ErrNoDocuments
	if notFound {
		return
	}

	if err = result.Decode(&value); err != nil {
		return nil, nil
	}

	event = &value
	return
}

// Find timer event from "from" to "to".
func (r *MongoDbRepository) FindTimerEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.TimerEvent, err error) {
	db, err := getMongoDb(ctx, config.Get("MONGODB_URI", ""))
	if err != nil {
		return
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, config.Get("MONGODB_COLLECTION", ""))

	// Find TimerEvent which event_time is between "from" and "to".
	filter := bson.D{
		{"utc_time_to_drink", bson.D{{"$gte", from}}},
		{"utc_time_to_drink", bson.D{{"$lte", to}}},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return
	}
	defer cur.Close(ctx)
	// Iterating the finding results.
	for cur.Next(ctx) {
		var elm enterpriserule.TimerEvent
		err = cur.Decode(&elm)
		if err != nil {
			return
		}
		results = append(results, &elm)
	}

	if err = cur.Err(); err != nil {
		return
	}

	return
}

// Save TimerEvent to DB.
//
// Return error and the slice of TimerEvent saved successfully.
func (r *MongoDbRepository) SaveTimerEvent(ctx context.Context, event *enterpriserule.TimerEvent) (saved *enterpriserule.TimerEvent, err error) {
	db, err := getMongoDb(ctx, config.Get("MONGODB_URI", ""))
	if err != nil {
		return nil, err
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, config.Get("MONGODB_COLLECTION", ""))

	var filter bson.M
	opts := options.Update().SetUpsert(true)

	filter = bson.M{"user_id": event.UserId}
	value := bson.D{{"$set", event}}
	_, err = collection.UpdateOne(ctx, filter, value, opts)
	if err == nil {
		saved = event
	}
	return
}
