package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Nishithcs/banking-ledger/internal/processor"
	internal "github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// TransactionHistoryItem represents a single transaction in the history
type TransactionHistoryItem struct {
	TransactionID           string    `json:"id"`
	Amount                  float64   `json:"amount"`
	TransactionType         string    `json:"type"`
	Status                  string    `json:"status"`
	Timestamp               time.Time `json:"timestamp"`
	Balance float64   `json:"balance"`
	Description             string    `json:"description,omitempty"`
}

// TransactionHistoryResponse represents the response structure for transaction history
type TransactionHistoryResponse struct {
	AccountNumber string                   `json:"accountNumber"`
	Transactions  []TransactionHistoryItem `json:"transactions"`
	TotalCount    int                      `json:"totalCount"`
}

// GetTransactionHistoryHandler returns a handler for querying transaction history
func GetTransactionHistoryHandler(ctx context.Context, mongoDbClient internal.MongoDBClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("accountNumber")
		if accountNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode": http.StatusBadRequest,
				"error":     "Account number is required",
			})
			return
		}
		
		// Mongo Db changes
		filter := bson.M{"account_number": accountNumber}

		cursor, err := mongoDbClient.Find(ctx, "transactions", filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": http.StatusInternalServerError,
				"error":     "Failed to search transaction history: " + err.Error(),
			})
		}
		defer cursor.Close(ctx)

		var results []processor.TransactionDocument
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": http.StatusInternalServerError,
				"error":     "Failed to search transaction history: " + err.Error(),
			})
		}

		for _, tx := range results {
			fmt.Printf("Transaction: %+v\n", tx)
		}

		// Extract transactions from the response
		transactions := []TransactionHistoryItem{}
		for _, tx := range results {
			transactionHistoryItem := TransactionHistoryItem{
				TransactionID:           tx.TransactionID,
				Amount:                  tx.Amount,
				TransactionType:         tx.Type,
				Status:                  tx.Status,
				Timestamp:               tx.Timestamp,
				Balance: tx.Balance,
			}

			transactions = append(transactions, transactionHistoryItem)
		}

		// mongo changes done

		// Build the response
		response := TransactionHistoryResponse{
			AccountNumber: accountNumber,
			Transactions:  transactions,
			TotalCount:    len(results),
		}

		c.JSON(http.StatusOK, response)
	}
}
