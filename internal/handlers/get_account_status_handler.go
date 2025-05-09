package handlers

import (
	"context"
	"net/http"

	internal "github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-gonic/gin"
)

// AccountStatusResponse represents the response structure for account status
type AccountStatusResponse struct {
	Status        string `json:"status"`
	AccountNumber string `json:"accountNumber,omitempty"`
}

// GetAccountStatusHandler handles requests to check account creation status
func GetAccountStatusHandler(ctx context.Context, database internal.Database, mongoDbClient internal.MongoDBClient) gin.HandlerFunc {
	
	return func(c *gin.Context) {
		accountNumber := c.Param("accountNumber")

		if accountNumber == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"errorCode": http.StatusBadRequest,
				"error":     "Reference ID is required",
			})
			return
		}

		// Construct response
		response := AccountStatusResponse{
			AccountNumber: accountNumber,
		}

		tx, err := database.Begin(ctx)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"errorCode": http.StatusNotFound,
				"error":     "Account creation request not found",
			})
		}
		defer tx.Rollback(ctx) // Will be ignored if transaction is committed

		var status string

		// We will use the FOR UPDATE to avoid race conditions on the row.
		query := `SELECT status FROM accounts WHERE account_number = $1 FOR UPDATE;`

		err = tx.QueryRow(ctx, query, accountNumber).Scan(&status)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"errorCode": http.StatusNotFound,
				"error":     "Account creation request not found",
				"account_number": accountNumber,
			})
		}
		response.Status = status

		c.JSON(http.StatusOK, response)
	}
}
