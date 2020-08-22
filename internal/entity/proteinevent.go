package entity

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"proteinreminder/internal/drivers"
	"time"
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
		return nil, err
	}

	return &value, nil
}

func FindProteinEventByTime(ctx context.Context, from, to time.Time) ([]*ProteinEvent, error) {

	var results []*ProteinEvent

	db, err := drivers.GetMongoDb(ctx)
	if err != nil {
		return nil, err
	}
	defer drivers.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := drivers.GetMongoDbCollection(db)

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
		var elm ProteinEvent
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
func SaveProteinEvent(ctx context.Context, events []*ProteinEvent) (error, []*ProteinEvent) {
	db, err := drivers.GetMongoDb(ctx)
	if err != nil {
		return err, nil
	}
	defer drivers.DisConnectMongoDbClientFunc(ctx, db.Client(), func(err error) {
		return
	})()

	collection := drivers.GetMongoDbCollection(db)

	savedEvents := make([]*ProteinEvent, 0, len(events))
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
	return nil, savedEvents
}

// --------------------------------------------------------
// Entity

type ProteinEvent struct {
	UserId               string        `bson:"user_id"`
	UtcTimeToDrink       time.Time     `bson:"utc_time_to_drink"`
	DrinkTimeIntervalSec time.Duration `bson:"drink_time_interval_sec"`
}

func NewProteinEvent(userId string) (*ProteinEvent, error) {
	if userId == "" {
		return nil, errors.New("userId should be set")
	}

	p := &ProteinEvent{
		UserId: userId,
	}
	return p, nil
}

func (p *ProteinEvent) Equal(another *ProteinEvent) bool {
	if p == nil || another == nil {
		return false
	}
	if p.UserId != another.UserId {
		return false
	}
	if p.DrinkTimeIntervalSec != another.DrinkTimeIntervalSec {
		return false
	}
	if p.UtcTimeToDrink.Second() != another.UtcTimeToDrink.Second() {
		return false
	}
	return true
}
