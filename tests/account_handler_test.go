package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nishithcs/banking-ledger/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccount_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockQueue := &MockMessageQueue{}
	router := gin.Default()
	handler := handlers.CreateAccount(context.Background(), mockQueue, "test-queue")
	router.POST("/create-account", handler)

	body := map[string]interface{}{
		"name":           "John Smith",
		"initialDeposit": 1000.0,
	}
	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/create-account", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusAccepted, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["referenceID"])
	assert.NotEmpty(t, response["createdAt"])
	assert.NotEmpty(t, response["accountNumber"])
}

func TestCreateAccount_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockQueue := &MockMessageQueue{}
	router := gin.Default()
	handler := handlers.CreateAccount(context.Background(), mockQueue, "test-queue")
	router.POST("/create-account", handler)

	// Missing required fields
	body := map[string]interface{}{
		"name": "Jake Smith",
	}
	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/create-account", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// func TestCreateAccount_QueuePublishError(t *testing.T) {
//     gin.SetMode(gin.TestMode)

//     mockQueue := &MockMessageQueue{ShouldFail: true}
//     router := gin.Default()
//     handler := handlers.CreateAccount(context.Background(), mockQueue, "test-queue")
//     router.POST("/create-account", handler)

//     body := map[string]interface{}{
//         "name":           "Alice Doe",
//         "initialDeposit": 750.0,
//     }
//     bodyJSON, _ := json.Marshal(body)
//     req, _ := http.NewRequest(http.MethodPost, "/create-account", bytes.NewBuffer(bodyJSON))
//     req.Header.Set("Content-Type", "application/json")

//     resp := httptest.NewRecorder()
//     router.ServeHTTP(resp, req)

//     assert.Equal(t, http.StatusInternalServerError, resp.Code)

//     var response map[string]interface{}
//     err := json.Unmarshal(resp.Body.Bytes(), &response)
//     assert.NoError(t, err)
//     assert.Equal(t, "mock publish error", response["error"])
// }

func TestCreateAccount_PublishFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockQueue := &MockMessageQueue{ShouldFail: true}
	router := gin.Default()
	handler := handlers.CreateAccount(context.Background(), mockQueue, "test-queue")
	router.POST("/create-account", handler)

	body := map[string]interface{}{
		"name":           "Jake smith",
		"initialDeposit": 500.0,
	}
	bodyJSON, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "/create-account", bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)

	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(http.StatusInternalServerError), response["errorCode"])
}
