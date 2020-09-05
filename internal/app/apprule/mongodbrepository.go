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
}

func NewMongoDbRepository() Repository {
	return &MongoDbRepository{}
}

func (r *MongoDbRepository) FindProteinEvent(ctx context.Context, userId string) (*enterpriserule.ProteinEvent, error) {
	db, err := driver.GetMongoDb(ctx, config.Get("MONGODB_URI"))
	if err != nil {
		return nil, err
	}
	defer driver.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := driver.GetMongoDbCollection(db, config.Get("MONGODB_COLLECTION"))

	var value enterpriserule.ProteinEvent
	filter := bson.M{"user_id": userId}
	result := collection.FindOne(ctx, filter)
	notFound := result.Err() == mongo.ErrNoDocuments
	if notFound {
		return nil, nil
	}

	if err := result.Decode(&value); err != nil {
		return nil, err
	}

	return &value, nil
}

func (r *MongoDbRepository) FindProteinEventByTime(ctx context.Context, from, to time.Time) ([]*enterpriserule.ProteinEvent, error) {
	var results []*enterpriserule.ProteinEvent

	db, err := driver.GetMongoDb(ctx, config.Get("MONGODB_URI"))
	if err != nil {
		return nil, err
	}
	defer driver.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := driver.GetMongoDbCollection(db, config.Get("MONGODB_COLLECTION"))

	// Find ProteinEvent which event_time is between "from" and "to".
	filter := bson.D{
		{"utc_time_to_drink", bson.D{{"$gte", from}}},
		{"utc_time_to_drink", bson.D{{"$lte", to}}},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	// Iterating the finding results.
	for cur.Next(ctx) {
		var elm enterpriserule.ProteinEvent
		err := cur.Decode(&elm)
		if err != nil {
			return nil, err
		}

		results = append(results, &elm)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Save ProteinEvent to DB.
//
// Return error and the slice of ProteinEvent saved successfully.
func (r *MongoDbRepository) SaveProteinEvent(ctx context.Context, events []*enterpriserule.ProteinEvent) ([]*enterpriserule.ProteinEvent, error) {
	db, err := driver.GetMongoDb(ctx, config.Get("MONGODB_URI"))
	if err != nil {
		return nil, err
	}
	defer driver.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := driver.GetMongoDbCollection(db, config.Get("MONGODB_COLLECTION"))

	savedEvents := make([]*enterpriserule.ProteinEvent, 0, len(events))
	var filter bson.M
	opts := options.Update().SetUpsert(true)
	for _, event := range events {
		filter = bson.M{"user_id": event.UserId}
		value := bson.D{{"$set", event}}
		_, err = collection.UpdateOne(ctx, filter, value, opts)
		if err == nil {
			savedEvents = append(savedEvents, event)
		}
	}
	return savedEvents, nil
}
