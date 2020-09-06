package apprule

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"proteinreminder/internal/app/driver"
	"proteinreminder/internal/app/enterpriserule"
	"time"
	"proteinreminder/internal/pkg/config"
)

// Implements Repository interface with MongoDB.
type MongoDbRepository struct {
	conn   driver.MongoDbConnector
	config config.Config
}

func NewMongoDbRepository(conn driver.MongoDbConnector, config config.Config) Repository {
	return &MongoDbRepository{
		conn,
		config,
	}
}

func (r *MongoDbRepository) FindProteinEvent(ctx context.Context, userId string) (event *enterpriserule.ProteinEvent, err error) {
	db, err := r.conn.GetDb(ctx, r.config.Get("MONGODB_URI"))
	if err != nil {
		return
	}
	defer r.conn.DisConnectClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := r.conn.GetCollection(db, r.config.Get("MONGODB_COLLECTION"))

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

func (r *MongoDbRepository) FindProteinEventByTime(ctx context.Context, from, to time.Time) (results []*enterpriserule.ProteinEvent, err error) {
	db, err := r.conn.GetDb(ctx, r.config.Get("MONGODB_URI"))
	if err != nil {
		return
	}
	defer r.conn.DisConnectClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := r.conn.GetCollection(db, r.config.Get("MONGODB_COLLECTION"))

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
func (r *MongoDbRepository) SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) (results []*enterpriserule.ProteinEvent, err error) {
	db, err := r.conn.GetDb(ctx, r.config.Get("MONGODB_URI"))
	if err != nil {
		return nil, err
	}
	defer r.conn.DisConnectClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := r.conn.GetCollection(db, r.config.Get("MONGODB_COLLECTION"))

	events = make([]*enterpriserule.ProteinEvent, 0, len(events))
	var filter bson.M
	opts := options.Update().SetUpsert(true)
	for _, event := range events {
		filter = bson.M{"user_id": event.UserId}
		value := bson.D{{"$set", event}}
		_, err = collection.UpdateOne(ctx, filter, value, opts)
		if err == nil {
			events = append(events, event)
		}
	}
	return
}
