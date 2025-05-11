# Banking Ledger API

A simple banking ledger system built with Go, PostgreSQL, MongoDB, and RabbitMQ. It allows you to create accounts, make transactions, and query account and transaction details through a clean RESTful API.

---

## ⚙️ Setup Instructions
cd te 
### 1. Clone the Repository
```bash
git clone https://github.com/Nishithcs/banking-ledger.git
cd banking-ledger
```

### 2. Start the Services Using Docker Compose
```bash
docker compose up -d
```

This will spin up the following services:
- API Service
- Account Creator Worker
- Transaction Processor Worker
- PostgreSQL
- MongoDB
- RabbitMQ

### 3. Check Logs
```bash
docker compose logs -f
```

---

## API Endpoints

Base URL: `http://localhost:8080`

### 1. Create a Bank Account  
**POST** `/accounts`

**Request:**
```json
{
  "name": "Virat Kohli",
  "initialDeposit": 1000
}
```

**Curl:**
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
        "name": "Virat Kohli",
        "initialDeposit": 1000
      }'
```

**Response:**
```json
{
  "accountNumber": "9080900634",
  "createdAt": "2025-05-09T20:29:16.680558126Z",
  "referenceID": "4af998d4-886e-4d2a-8d70-7198554e5285"
}
```

---

### 2. Get Account Status  
**GET** `/account/:accountNumber/status`

**Example:**
```bash
curl http://localhost:8080/account/9080900634/status
```

---

### 3. Make a Transaction  
**POST** `/transactions`

**Request:**
```json
// Deposit
{
  "accountNumber": "8225565849",
  "amount": 400,
  "type": "DEPOSIT"
}

// Withdraw
{
  "accountNumber": "8225565849",
  "amount": 200,
  "type": "WITHDRAWAL"
}
```

**Curl:**
```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d '{
        "accountNumber": "8225565849",
        "amount": 400,
        "type": "DEPOSIT"
      }'
```

**Response:**
```json
{
  "createdAt": "2025-05-10T08:30:22.563Z",
  "transactionId": "3fa85f64-5717-4562-b3fc-2c963f66afa6"
}
```

---

### 4. Get Transactions status
**GET** `/transaction/:transactionId/status`

**Example:**
```bash
curl http://localhost:8080/transaction/3fa85f64-5717-4562-b3fc-2c963f66afa6/status
```

**Response:**
```json
{
  "id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "amount": 400,
  "type": "DEPOSIT",
  "status": "COMPLETED",
  "timestamp": "2025-05-10T20:44:17.184Z",
  "balance": 1400
}
```

---

### 5. Get All Transactions for an Account  
**GET** `/accounts/:accountNumber/transactions`

**Example:**
```bash
curl http://localhost:8080/accounts/8225565849/transactions
```

**Response:**
```json
{
  "accountNumber": "8225565849",
  "transactions": [
    {
      "id": "4af998d4-886e-4d2a-8d70-7198554e5285",
      "amount": 1000,
      "type": "DEPOSIT",
      "status": "COMPLETED",
      "timestamp": "2025-05-09T20:29:16.681Z",
      "balance": 1000
    },
    {
      "id": "551fd5f7-4e83-45c4-86b3-23f76e830a71",
      "amount": 400,
      "type": "DEPOSIT",
      "status": "COMPLETED",
      "timestamp": "2025-05-09T20:29:44.366Z",
      "balance": 1400
    },
    {
      "id": "c45765ed-29c7-42ba-a2b3-9f159ae0bcdd",
      "amount": 400,
      "type": "DEPOSIT",
      "status": "FAILED",
      "timestamp": "2025-05-10T19:24:46.574Z",
      "balance": 0
    }
  ],
  "totalCount": 3
}

```

## System Architecture

The system follows a modular, decoupled architecture comprising the following components:

- **API Service**: Exposes RESTful endpoints for Frontend application interactions and request handling.

- **Account Creator Worker**: Subscribes to RabbitMQ queues to process account creation events asynchronously.

- **Transaction Processor Worker**: Handles transaction requests in the background via message queues.

- **PostgreSQL**: Acts as the primary data store for persistent account-related data.

- **MongoDB**: Stores transaction history, optimized for high-performance read operations.

- **RabbitMQ**: Serves as the message broker enabling asynchronous communication between services.

---
