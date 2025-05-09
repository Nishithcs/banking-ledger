package processor

import (
	"time"

	internal "github.com/Nishithcs/banking-ledger/pkg"
)

type ProcessWorker struct {
	Database     internal.Database
	MongoDbConn internal.MongoDBClient
}

// Log the account creation transaction to Elasticsearch
type TransactionDocument struct {
	TransactionID           string    `json:"transaction_id"  bson:"transaction_id"`
	AccountNumber           string    `json:"account_number"  bson:"account_number"`
	Amount                  float64   `json:"amount"  bson:"amount"`
	Type                    string    `json:"type"  bson:"type"`
	Status                  string    `json:"status"  bson:"status"`
	Timestamp               time.Time `json:"timestamp"  bson:"timestamp"`
	Balance                 float64   `json:"balance"  bson:"balance"`
}