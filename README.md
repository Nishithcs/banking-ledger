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
  "accountHolderName": "qqqqq",
  "initialDeposit": 1000
}
```

**Curl:**
```bash
curl -X POST http://localhost:8080/accounts \
  -H "Content-Type: application/json" \
  -d '{
        "accountHolderName": "qqqqq",
        "initialDeposit": 1000
      }'
```

**Response:**
```json
{
  "accountNumber": "8225565849",
  "createdAt": "2025-05-09T20:29:16.680558126Z",
  "referenceID": "4af998d4-886e-4d2a-8d70-7198554e5285"
}
```

---

### 2. Make a Transaction  
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

### 3. Get All Transactions for an Account  
**GET** `/accounts/:accountNumber/transactions`

**Example:**
```bash
curl http://localhost:8080/accounts/8225565849/transactions
```

---

### 4. Get Account Status  
**GET** `/account/:accountNumber/status`

**Example:**
```bash
curl http://localhost:8080/account/8225565849/status
```

---

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
