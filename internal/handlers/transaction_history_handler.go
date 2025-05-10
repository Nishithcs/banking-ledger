package handlers

import (
	"context"
	"net/http"

	"github.com/Nishithcs/banking-ledger/internal/processor"
	"github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// GetTransactions returns the list of transactions made for an account
func GetTransactions(ctx context.Context, mongoDbClient pkg.MongoDBClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("accountNumber")
		if accountNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode": http.StatusBadRequest,
				"error":     "Account number is required",
			})
			return
		}
		
		filter := bson.M{"account_number": accountNumber}

		cursor, err := mongoDbClient.Find(ctx, "transactions", filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": http.StatusInternalServerError,
				"error":     "Failed to search transaction history: " + err.Error(),
			})
		}

		if cursor == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No cursor returned"})
			return
		}
		defer cursor.Close(ctx)

		var results []processor.TransactionDocument
		if err = cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errorCode": http.StatusInternalServerError,
				"error":     "Failed to search transaction history: " + err.Error(),
			})
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
				Balance:                 tx.Balance,
			}
			transactions = append(transactions, transactionHistoryItem)
		}

		// Build the response
		response := TransactionHistoryResponse{
			AccountNumber: accountNumber,
			Transactions:  transactions,
			TotalCount:    len(results),
		}

		c.JSON(http.StatusOK, response)
	}
}
