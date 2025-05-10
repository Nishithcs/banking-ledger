package handlers_test

import (
	"context"
	"testing"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"net/http/httptest"
)


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