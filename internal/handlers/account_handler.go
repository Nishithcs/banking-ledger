package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"math/rand"

	"github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateAccount creates a new Account Number and pushes the task to the queue. 
// The worker will do rest of the account creation.
func CreateAccount(ctx context.Context, messageQueue pkg.MessageQueue, queueName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accountRequestJson AccountRequest

		if err := c.ShouldBindJSON(&accountRequestJson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode": http.StatusBadRequest,
				"error":     err.Error(),
			})
			return
		}

		// Process the account creation request
		response, err := accountRequestJson.createAccount(ctx, messageQueue, queueName)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": http.StatusInternalServerError,
				"error":     err.Error(),
			})
			return
		}

		// Return a successful response with tracking information
		c.JSON(http.StatusAccepted, gin.H{
			"referenceID": response.ReferenceID,
			"createdAt":   response.CreatedAt,
			"accountNumber": response.AccountNumber,
		})
	}
}

// createAccount processes an account creation request by:
// 1. Generating a unique reference ID
// 2. Publishing the request to a message queue for asynchronous processing
// 3. Returning a response with tracking information
func (a *AccountRequest) createAccount(ctx context.Context, messageQueue pkg.MessageQueue, queueName string) (AccountResponse, error) {
	// Generate a unique reference ID for tracking this request
	a.ReferenceID = uuid.New().String()

	// Generate a unique account number
	timestamp := time.Now().Unix()
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(10000) // 4 digits: 0000–9999
	
	// 10 digits = last 6 of timestamp + 4-digit random
	a.AccountNumber= fmt.Sprintf("%06d%04d", timestamp%1000000, randomNumber)

	// Create JSON payload
	requestByteArray, err := json.Marshal(a)
	if err != nil {
		// Handle marshaling error
		return AccountResponse{}, err
	}

	// Publish message to RabbitMQ
	err = messageQueue.Publish(queueName, requestByteArray)

	if err != nil {
		// Handle publishing error
		return AccountResponse{}, err
	}

	// Return response with tracking ID and timestamp
	return AccountResponse{
		AccountNumber: a.AccountNumber,
		ReferenceID: a.ReferenceID,
		CreatedAt:   time.Now(),
	}, nil
}
