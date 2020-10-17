package apprule

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"proteinreminder/internal/app/enterpriserule"
	"proteinreminder/internal/pkg/config"
	"strings"
	"time"
)

// Implements Repository interface with MongoDB.
type MongoDbRepository struct {
	config config.Config
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

func NewMongoDbRepository(config config.Config) Repository {
	return &MongoDbRepository{
		config,
	}
}

// Find protein event by user id.
func (r *MongoDbRepository) FindProteinEvent(ctx context.Context, userId string) (event *enterpriserule.ProteinEvent, err error) {
	db, err := getMongoDb(ctx, r.config.Get("MONGODB_URI", ""))
	if err != nil {
		return
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, r.config.Get("MONGODB_COLLECTION", ""))

	var value enterpriserule.ProteinEvent
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

// Find protein event from "from" to "to".
func (r *MongoDbRepository) FindProteinEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.ProteinEvent, err error) {
	db, err := getMongoDb(ctx, r.config.Get("MONGODB_URI", ""))
	if err != nil {
		return
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, r.config.Get("MONGODB_COLLECTION", ""))

	// Find ProteinEvent which event_time is between "from" and "to".
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
		var elm enterpriserule.ProteinEvent
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

// Save ProteinEvent to DB.
//
// Return error and the slice of ProteinEvent saved successfully.
func (r *MongoDbRepository) SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) (saved []*enterpriserule.ProteinEvent, err error) {
	db, err := getMongoDb(ctx, r.config.Get("MONGODB_URI", ""))
	if err != nil {
		return nil, err
	}
	defer disconnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := getMongoCollection(db, r.config.Get("MONGODB_COLLECTION", ""))

	saved = make([]*enterpriserule.ProteinEvent, 0, len(events))
	var filter bson.M
	opts := options.Update().SetUpsert(true)
	for _, event := range events {
		filter = bson.M{"user_id": event.UserId}
		value := bson.D{{"$set", event}}
		_, err = collection.UpdateOne(ctx, filter, value, opts)
		if err == nil {
			saved = append(saved, event)
		}
	}
	return
}
