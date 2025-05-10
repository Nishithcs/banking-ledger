# 🏦 Banking Ledger API

A simple banking ledger system built with Go, PostgreSQL, MongoDB, and RabbitMQ. It allows you to create accounts, make transactions, and query account and transaction details through a clean RESTful API.

---

## 📌 API Endpoints

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
      }' | jq
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
curl http://localhost:8080/account/9080900634/status | jq
```

---

### 3. Make a Transaction  
**POST** `/transactions`

**Request:**
```json
{
  "accountNumber": "8225565849",
  "amount": 400,
  "type": "DEPOSIT"
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
curl http://localhost:8080/transaction/3fa85f64-5717-4562-b3fc-2c963f66afa6/status | jq
```

---

### 5. Get All Transactions for an Account  
**GET** `/accounts/:accountNumber/transactions`

**Example:**
```bash
curl http://localhost:8080/accounts/8225565849/transactions | jq
```

## ⚙️ Setup Instructions

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

## 🧱 System Architecture

The system is composed of multiple decoupled components:

- **API Service**: RESTful API to handle client interactions.
- **Account Creator Worker**: Listens to account creation events via RabbitMQ.
- **Transaction Processor Worker**: Handles transaction requests asynchronously.
- **PostgreSQL**: Stores persistent account data.
- **MongoDB**: Stores transaction history optimized for reads.
- **RabbitMQ**: Message broker for asynchronous processing.

---
