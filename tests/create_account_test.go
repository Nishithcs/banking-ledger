package handlers_test

import (
	"context"
	"testing"

	"github.com/Nishithcs/banking-ledger/internal/processor"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_NegativeDeposit(t *testing.T) {
	// Create a CreateAccountProcessor with negative initial deposit
	accountData := processor.AccountData{
		AccountNumber:  "123456789",
		Name:           "John Doe",
		InitialDeposit: -100.0,
		ReferenceID:    "ref123",
	}

	processor := processor.CreateAccountProcessor{
		Data: accountData,
	}

	// Call CreateAccount
	err := processor.CreateAccount(context.Background())

	// Assert that an error occurred
	assert.Error(t, err)
	assert.Equal(t, "initial Deposit cannot be negative", err.Error())
}
