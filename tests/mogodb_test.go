package handlers_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockMongoDBClient is a mock implementation of the MongoDBClient interface
type MockMongoDBClient struct {
	mock.Mock
}

func (m *MockMongoDBClient) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, collection, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

func (m *MockMongoDBClient) Find(ctx context.Context, collection string, filter interface{}) (*mongo.Cursor, error) {
	args := m.Called(ctx, collection, filter)

	var cursor *mongo.Cursor
	if args.Get(0) != nil {
		cursor = args.Get(0).(*mongo.Cursor)
	}
	return cursor, args.Error(1)
}

func (m *MockMongoDBClient) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}