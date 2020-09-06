package driver

import "go.mongodb.org/mongo-driver/mongo"

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/url"
	"strings"
)

// Define this interface to use gomock.
type MongoDbConnector interface {
	GetDb(ctx context.Context, mongoUrl string) (db *mongo.Database, err error)
	GetCollection(db *mongo.Database, name string) (col *mongo.Collection)
	DisConnectClientFunc(ctx context.Context, client *mongo.Client, f func(error)) (fnc func())
}

// Define these interface to use gomock.
type MongoDatabase interface {
	Client() *mongo.Client
}
type MongoCollection interface {
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult
}
type MongoSingleResult interface {
	Err() error
	Decode(v interface{}) error
}
type MongoClient interface {
}

// Implements MongoDb interface.
type ConcreteMongoDbConnector struct {
}

func NewMongoDbConnector() (db MongoDbConnector) {
	db = &ConcreteMongoDbConnector{}
	return
}

func (c *ConcreteMongoDbConnector) GetDb(ctx context.Context, mongoUrl string) (db *mongo.Database, err error) {
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

func (c *ConcreteMongoDbConnector) GetCollection(db *mongo.Database, name string) (col *mongo.Collection) {
	col = db.Collection(name)
	return
}

func (c *ConcreteMongoDbConnector) DisConnectClientFunc(ctx context.Context, client *mongo.Client, f func(error)) (fnc func()) {
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
