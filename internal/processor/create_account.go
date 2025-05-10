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
	Name              string  `json:"name"`
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
			name, 
			balance, 
			status,
			created_at
		) VALUES ($1, $2, $3, $4, $5)
	`

	now := time.Now()
	err := p.Database.Exec(
		ctx,
		query,
		p.Data.AccountNumber,
		p.Data.Name,
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
		Balance: p.Data.InitialDeposit,
	}

	//Insert data to MongoDb
	_, err = p.MongoDbConn.Insert(ctx, "transactions", transactionDoc)
	if err != nil {
		log.Printf("Failed to insert transaction in MongoDB: %v", err)
	}


	return nil
}
