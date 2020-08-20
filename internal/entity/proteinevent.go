package entity

import (
	"time"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"proteinreminder/internal/drivers"
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --------------------------------------------------------
// Repository Role

func GetProteinEvent(ctx context.Context, userId string) (*ProteinEvent, error) {
	db, err := drivers.GetMongoDb(ctx)
	if err != nil {
		return nil, err
	}
	defer drivers.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := drivers.GetMongoDbCollection(db)

	var value ProteinEvent
	filter := bson.M{"user_id": userId}
	if err := collection.FindOne(ctx, filter).Decode(&value); err != nil {
		value = ProteinEvent{
			userId: userId,
		}
	}
	return &value, nil
}

func FindProteinEventByTime(ctx context.Context, from, to time.Time) (error, []*ProteinEvent) {

	var results []*ProteinEvent

	db, err := drivers.GetMongoDb(ctx)
	if err != nil {
		return err, nil
	}
	defer drivers.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := drivers.GetMongoDbCollection(db)

	// Find ProteinEvent which event_time is between "from" and "to".
	filter := bson.D{
		{"$ge", bson.D{{"event_time", ""}}},
		{"$le", bson.D{{"event_time", ""}}},
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return err, nil
	}

	// Iterating the finding results.
	for cur.Next(ctx) {
		var elm ProteinEvent
		err := cur.Decode(&elm)
		if err != nil {
			return err, nil
		}

		results = append(results, &elm)
	}

	if err := cur.Err(); err != nil {
		return err, nil
	}

	return nil, results
}

// Save ProteinEvent to DB.
//
// Return error and the slice of ProteinEvent saved successfully.
func SaveProteinEvent(ctx context.Context, events []*ProteinEvent) (error, []*ProteinEvent) {
	db, err := drivers.GetMongoDb(ctx)
	if err != nil {
		return err, nil
	}
	defer drivers.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := drivers.GetMongoDbCollection(db)

	savedEvents := make([]*ProteinEvent, len(events))
	var filter bson.M
	opts := options.Update().SetUpsert(true)
	for _, event := range events {
		filter = bson.M{"user_id": event.userId}
		value := bson.D{{"$set", event}}
		_, err = collection.UpdateOne(ctx, filter, value, opts)
		if err == nil {
			savedEvents = append(savedEvents, event)
		}
	}
	return nil, savedEvents
}

// --------------------------------------------------------
// Entity

type ProteinEvent struct {
	userId               string        `bson:"_id"`
	UtcTimeToDrink       time.Time     `bson:"utc_time_to_drink"`
	DrinkTimeIntervalSec time.Duration `bson:"drink_time_interval_sec"`
}

func NewProteinEvent(userId string) (*ProteinEvent, error) {
	if userId == "" {
		return nil, errors.New("userId should be set")
	}

	p := &ProteinEvent{
		userId: userId,
	}
	return p, nil
}

