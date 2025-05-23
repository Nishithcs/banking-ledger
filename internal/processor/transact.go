package processor

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (p *TransactionProcessor) transact(ctx context.Context) error {
	// Start a transaction
	tx, err := p.Database.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // Will be ignored if transaction is committed

	// Get current account balance. For update to avoid race conditions
	var currentBalance float64
	query := `SELECT balance FROM accounts 
	WHERE account_number = $1 AND status='ACTIVE' FOR UPDATE;`

	err = tx.QueryRow(ctx, query, p.Data.AccountNumber).Scan(&currentBalance)
	if err != nil {
		return fmt.Errorf("failed to get account balance: %w", err)
	}

	p.Data.AvailableBalance = currentBalance

	// Calculate new balance based on transaction type
	var newBalance float64
	switch p.Data.Type {
	case "DEPOSIT":
		if p.Data.Amount <= 0 {
			return fmt.Errorf("deposit amount must be positive")
		}
		newBalance = currentBalance + p.Data.Amount
	case "WITHDRAWAL":
		if p.Data.Amount <= 0 {
			return fmt.Errorf("withdrawal amount must be positive")
		}
		if currentBalance < p.Data.Amount {
			return fmt.Errorf("insufficient funds")
		}
		newBalance = currentBalance - p.Data.Amount
	default:
		return fmt.Errorf("invalid transaction type: %s", p.Data.Type)
	}

	p.Data.AvailableBalance = newBalance

	// Update account balance
	updateQuery := `UPDATE accounts SET balance = $1 WHERE account_number = $2`
	err = tx.Exec(ctx, updateQuery, newBalance, p.Data.AccountNumber)
	if err != nil {
		return fmt.Errorf("failed to update account balance: %w", err)
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ProcessTransaction handles deposit and withdrawal transactions
func (p *TransactionProcessor) ProcessTransaction(ctx context.Context) error {

	err := p.transact(ctx)
	status := "COMPLETED"
	if err != nil {
		log.Println(err)
		status = "FAILED"
	}
	transactionDoc := TransactionDocument{
		AccountNumber: p.Data.AccountNumber,
		Type:          p.Data.Type,
		Amount:        p.Data.Amount,
		TransactionID: p.Data.TransactionID,
		Timestamp:     time.Now(),
		Status:        status,
		Balance: p.Data.AvailableBalance,
	}

	//Insert data to MongoDb
	_, err = p.MongoDbConn.Insert(ctx, "transactions", transactionDoc)
	if err != nil {
		log.Printf("Failed to insert transaction in MongoDB: %v", err)
	}


	return nil
}
