package handlers

import (
	"time"
)

// AccountRequest represents the fields required for account creation
type AccountRequest struct {
	Name              string  `json:"name"           binding:"required"`
	InitialDeposit    float64 `json:"initialDeposit" binding:"required"`
	ReferenceID       string  `json:"referenceID"`                        
	AccountNumber	  string  `json:"accountNumber"`
}

// AccountResponse is the
type AccountResponse struct {
	ReferenceID     string    `json:"referenceID"`
	AccountNumber   string    `json:"AccountNumber"`
	CreatedAt       time.Time `json:"createdAt"`
}

// AccountStatusResponse represents the response structure for account status
type AccountStatusResponse struct {
	Status        string `json:"status"`
	AccountNumber string `json:"accountNumber,omitempty"`
}

// TransactionResponse represents the response structure sent back to clients
// after a successful transaction request
type TransactionResponse struct {
	TransactionID string    `json:"transactionID"`
	CreatedAt     time.Time `json:"createdAt"`
}

// TransactionHistoryItem represents a single transaction in the history
type TransactionHistoryItem struct {
	TransactionID           string    `json:"id"`
	Amount                  float64   `json:"amount"`
	TransactionType         string    `json:"type"`
	Status                  string    `json:"status"`
	Timestamp               time.Time `json:"timestamp"`
	Balance                 float64   `json:"balance"`
	Description             string    `json:"description,omitempty"`
}

// TransactionHistoryResponse represents the response structure for transaction history
type TransactionHistoryResponse struct {
	AccountNumber string                   `json:"accountNumber"`
	Transactions  []TransactionHistoryItem `json:"transactions"`
	TotalCount    int                      `json:"totalCount"`
}
