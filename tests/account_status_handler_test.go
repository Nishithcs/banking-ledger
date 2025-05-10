package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"fmt"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	"github.com/Nishithcs/banking-ledger/pkg"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDatabase := &MockDatabase{
		BeginFunc: func(ctx context.Context) (pkg.Transaction, error) {
			return &MockTransaction{
				QueryRowFunc: func(ctx context.Context, query string, args ...any) pkg.RowResult {
					return &MockRowResult{
						ScanFunc: func(dest ...any) error {
							*dest[0].(*string) = "active"
							return nil
						},
					}
				},
				RollbackFunc: func(ctx context.Context) error {
					return nil
				},
			}, nil
		},
	}

	router := gin.Default()
	handler := handlers.GetAccountStatus(context.Background(), mockDatabase, nil)
	router.GET("/account-status/:accountNumber", handler)

	req, _ := http.NewRequest(http.MethodGet, "/account-status/12345", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response handlers.AccountStatusResponse
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "active", response.Status)
	assert.Equal(t, "12345", response.AccountNumber)
}

func TestGetAccountStatus_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockDatabase := &MockDatabase{
		BeginFunc: func(ctx context.Context) (pkg.Transaction, error) {
			return &MockTransaction{
				QueryRowFunc: func(ctx context.Context, query string, args ...any) pkg.RowResult {
					return &MockRowResult{
						ScanFunc: func(dest ...any) error {
							return fmt.Errorf("no rows in result set")
						},
					}
				},
				RollbackFunc: func(ctx context.Context) error {
					return nil
				},
			}, nil
		},
	}

	router := gin.Default()
	handler := handlers.GetAccountStatus(context.Background(), mockDatabase, nil)
	router.GET("/account-status/:accountNumber", handler)

	req, _ := http.NewRequest(http.MethodGet, "/account-status/12345", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}
