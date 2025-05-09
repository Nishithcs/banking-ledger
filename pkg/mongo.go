package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)


type MongoDBClient interface {
	Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error)
	Find(ctx context.Context, collection string, filter interface{}) (*mongo.Cursor, error)
	Disconnect(ctx context.Context) error
}

type MongoDBWrapper struct {
	client   *mongo.Client
	database string
}

func NewMongoDBClient(ctx context.Context, uri string, dbName string) (MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoDBWrapper{
		client:   client,
		database: dbName,
	}, nil
}

func (m *MongoDBWrapper) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	return m.client.Database(m.database).Collection(collection).InsertOne(ctx, document)
}

func (m *MongoDBWrapper) Find(ctx context.Context, collection string, filter interface{}) (*mongo.Cursor, error) {
	return m.client.Database(m.database).Collection(collection).Find(ctx, filter)
}

func (m *MongoDBWrapper) Disconnect(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
