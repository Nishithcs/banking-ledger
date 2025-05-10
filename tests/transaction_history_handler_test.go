package handlers_test

import (
	"context"
	"testing"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
)

// MockMongoDBClient is a mock implementation of the MongoDBClient interface
type MockMongoDBClient struct {
	mock.Mock
}

func (m *MockMongoDBClient) Insert(ctx context.Context, collection string, document interface{}) (*mongo.InsertOneResult, error) {
	args := m.Called(ctx, collection, document)
	return args.Get(0).(*mongo.InsertOneResult), args.Error(1)
}

// func (m *MockMongoDBClient) Find(ctx context.Context, collection string, filter interface{}) (*mongo.Cursor, error) {
// 	args := m.Called(ctx, collection, filter)
// 	return args.Get(0).(*mongo.Cursor), args.Error(1)
// }
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


// TestGetTransactions_MongoDBFailure tests the case when MongoDB fails to return results
func TestGetTransactions_MongoDBFailure(t *testing.T) {
	// Create a mock MongoDB client
	mockMongoDbClient := new(MockMongoDBClient)

	// Mock the Find method to return an error
	mockMongoDbClient.On("Find", mock.Anything, "transactions", bson.M{"account_number": "12345"}).
		Return(nil, assert.AnError)

	// Set up the Gin router
	r := gin.Default()
	r.GET("/transactions/:accountNumber", handlers.GetTransactions(context.Background(), mockMongoDbClient))

	// Perform the request
	req, _ := http.NewRequest(http.MethodGet, "/transactions/12345", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to search transaction history")
}