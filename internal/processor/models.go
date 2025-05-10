package processor

import (
	"time"

	"github.com/Nishithcs/banking-ledger/pkg"
)

type ProcessWorker struct {
	Database     pkg.Database
	MongoDbConn pkg.MongoDBClient
}

// Log the account creation transaction to Mongo DB
type TransactionDocument struct {
	TransactionID           string    `json:"transaction_id"  bson:"transaction_id"`
	AccountNumber           string    `json:"account_number"  bson:"account_number"`
	Amount                  float64   `json:"amount"  bson:"amount"`
	Type                    string    `json:"type"  bson:"type"`
	Status                  string    `json:"status"  bson:"status"`
	Timestamp               time.Time `json:"timestamp"  bson:"timestamp"`
	Balance                 float64   `json:"balance"  bson:"balance"`
}


type CreateAccountProcessor struct {
	ProcessWorker
	Data AccountData
}

// AccountData represents the data needed to create a new account
type AccountData struct {
	AccountNumber     string  `json:"accountNumber"`
	Name              string  `json:"name"`
	InitialDeposit    float64 `json:"initialDeposit"`
	ReferenceID       string  `json:"referenceID"`
}

type TransactionProcessor struct {
	ProcessWorker
	Data TransactionData
}

// TransactionData represents the data needed for a transaction
type TransactionData struct {
	AccountNumber     string  `json:"accountNumber"`
	Amount            float64 `json:"amount"`
	AvailableBalance  float64 `json:"availableBalance"`
	Type              string  `json:"type"`
	TransactionID     string  `json:"transactionId"`
}