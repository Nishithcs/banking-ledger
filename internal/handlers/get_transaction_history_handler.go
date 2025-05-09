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
	BalanceAfterTransaction float64   `json:"updatedBalance"`
	Description             string    `json:"description,omitempty"`
}

// EsResponseItem represents a single document in the Elasticsearch response
type EsResponseItem struct {
	TransactionID           string    `json:"transaction_id"`
	AccountNumber           string    `json:"account_number"`
	Amount                  float64   `json:"amount"`
	TransactionType         string    `json:"type"`
	Status                  string    `json:"status"`
	Description             string    `json:"description,omitempty"`
	Timestamp               time.Time `json:"timestamp"`
	BranchCode              string    `json:"branch_code"`
	BalanceAfterTransaction float64   `json:"balance_after_transaction"`
}

// TransactionHistoryResponse represents the response structure for transaction history
type TransactionHistoryResponse struct {
	AccountNumber string                   `json:"accountNumber"`
	Transactions  []TransactionHistoryItem `json:"transactions"`
	TotalCount    int                      `json:"totalCount"`
	CurrentPage   int                      `json:"currentPage"`
}

// GetTransactionHistoryHandler returns a handler for querying transaction history
func GetTransactionHistoryHandler(ctx context.Context, esClient internal.ElasticsearchClient, mongoDbClient internal.MongoDBClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNumber := c.Param("accountNumber")
		if accountNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode": http.StatusBadRequest,
				"error":     "Account number is required",
			})
			return
		}

		// // Parse query parameters
		// page := 1
		// if pageParam := c.Query("pageNumber"); pageParam != "" {
		// 	fmt.Sscanf(pageParam, "%d", &page)
		// 	if page < 1 {
		// 		page = 1
		// 	}
		// }

		// limit := 10

		// // Calculate offset
		// from := (page - 1) * limit

		// // Build Elasticsearch query
		// query := map[string]interface{}{
		// 	"query": map[string]interface{}{
		// 		"match": map[string]interface{}{
		// 			"account_number": accountNumber,
		// 		},
		// 	},
		// 	"sort": []map[string]interface{}{
		// 		{
		// 			"timestamp": map[string]interface{}{
		// 				"order": "desc", // Most recent transactions first
		// 			},
		// 		},
		// 	},
		// 	"from": from,
		// 	"size": limit,
		// }

		// // Convert query to JSON
		// var buf bytes.Buffer
		// if err := json.NewEncoder(&buf).Encode(query); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"errorCode": http.StatusInternalServerError,
		// 		"error":     "Failed to build search query",
		// 	})
		// 	return
		// }

		// // Perform the search request
		// res, err := esClient.Search(
		// 	[]string{"bank-transactions-*"},
		// 	&buf,
		// )
		// if err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"errorCode": http.StatusInternalServerError,
		// 		"error":     "Failed to search transaction history: " + err.Error(),
		// 	})
		// 	return
		// }
		// defer res.Body.Close()

		// // Check for Elasticsearch errors
		// if res.IsError() {
		// 	var e map[string]interface{}
		// 	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		// 		c.JSON(http.StatusInternalServerError, gin.H{
		// 			"errorCode": http.StatusInternalServerError,
		// 			"error":     "Failed to parse error response from Elasticsearch",
		// 		})
		// 		return
		// 	}
		// 	// Return the Elasticsearch error
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"errorCode": http.StatusInternalServerError,
		// 		"error":     fmt.Sprintf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"]),
		// 	})
		// 	return
		// }

		// // Parse the response
		// var esResponse struct {
		// 	Hits struct {
		// 		Total struct {
		// 			Value int `json:"value"`
		// 		} `json:"total"`
		// 		Hits []struct {
		// 			Source EsResponseItem `json:"_source"`
		// 		} `json:"hits"`
		// 	} `json:"hits"`
		// }

		// if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		// 	c.JSON(http.StatusInternalServerError, gin.H{
		// 		"errorCode": http.StatusInternalServerError,
		// 		"error":     "Failed to parse search results",
		// 	})
		// 	return
		// }

		// // Extract transactions from the response
		// transactions := make([]TransactionHistoryItem, 0, len(esResponse.Hits.Hits))
		// for _, hit := range esResponse.Hits.Hits {
		// 	transactionHistoryItem := TransactionHistoryItem{
		// 		TransactionID:           hit.Source.TransactionID,
		// 		Amount:                  hit.Source.Amount,
		// 		TransactionType:         hit.Source.TransactionType,
		// 		Status:                  hit.Source.Status,
		// 		Timestamp:               hit.Source.Timestamp,
		// 		BalanceAfterTransaction: hit.Source.BalanceAfterTransaction,
		// 		Description:             hit.Source.Description,
		// 	}

		// 	transactions = append(transactions, transactionHistoryItem)
		// }
		
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
				BalanceAfterTransaction: tx.BalanceAfterTransaction,
			}

			transactions = append(transactions, transactionHistoryItem)
		}

		// mongo changes done

		// Build the response
		response := TransactionHistoryResponse{
			AccountNumber: accountNumber,
			Transactions:  transactions,
			TotalCount:    len(results),
			CurrentPage:   1,
		}

		c.JSON(http.StatusOK, response)
	}
}
