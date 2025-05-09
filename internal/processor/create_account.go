package processor

import (
	"context"
	"fmt"
	"log"
	"time"
)

type CreateAccountProcessor struct {
	ProcessWorker
	Data AccountData
}

// AccountData represents the data needed to create a new account
type AccountData struct {
	AccountNumber     string  `json:"accountNumber"`
	AccountHolderName string  `json:"accountHolderName"`
	InitialDeposit    float64 `json:"initialDeposit"`
	ReferenceID       string  `json:"referenceID"`
}

// CreateAccount creates a new account in the database
func (p *CreateAccountProcessor) CreateAccount(ctx context.Context) error {

	// Validate available balance is not negative
	if p.Data.InitialDeposit < 0 {
		return fmt.Errorf("initial Deposit cannot be negative")
	}

	query := `
		INSERT INTO accounts (
			account_number, 
			account_holder_name, 
			available_balance, 
			branch_code, 
			status,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	_, err := p.PgxConn.Exec(
		ctx,
		query,
		p.Data.AccountNumber,
		p.Data.AccountHolderName,
		p.Data.InitialDeposit,
		"ACTIVE",
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create account: %w", err)
	}


	transactionDoc := TransactionDocument{
		TransactionID:           p.Data.ReferenceID,
		AccountNumber:           p.Data.AccountNumber,
		Amount:                  p.Data.InitialDeposit,
		Type:                    "DEPOSIT",
		Status:                  "COMPLETED",
		Timestamp:               now,
		BalanceAfterTransaction: p.Data.InitialDeposit,
	}

	// // Create index name with date format for better data management
	// indexName := fmt.Sprintf("bank-transactions-%s", now.Format("2006-01-02"))

	// // Index the transaction document
	// transactionDocJSON, err := json.Marshal(transactionDoc)
	// if err != nil {
	// 	panic(fmt.Sprintf("failed to marshal JSON: %v", err))
	// }

	
	// // Insert data to Elasticsearch
	// req := strings.NewReader(string(transactionDocJSON))
	// res, err := p.EsConn.Index(indexName, req)
	// if err != nil {
	// 	// Log the error but don't fail the account creation
	// 	log.Printf("Failed to index transaction in Elasticsearch: %v", err)
	// }
	// if res != nil {
	// 	defer res.Body.Close()
	// }

	//Insert data to MongoDb
	_, err = p.MongoDbConn.Insert(ctx, "transactions", transactionDoc)
	if err != nil {
		log.Printf("Failed to insert transaction in MongoDB: %v", err)
	}


	return nil
}
