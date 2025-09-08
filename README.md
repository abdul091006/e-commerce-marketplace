# Wallet Service Documentation

## Overview

The Wallet Service is a Golang-based microservice designed to manage digital wallets for an e-commerce marketplace. This service provides RESTful APIs for wallet creation, balance management, and transaction processing, with seamless integration capabilities for Frappe applications.

## Architecture

- **Language**: Go 1.22
- **Framework**: Fiber (HTTP web framework)
- **Database**: PostgreSQL 16
- **Containerization**: Docker & Docker Compose
- **Integration**: Frappe Framework compatible

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1/wallets
```

### 1. Create Wallet

Creates a new wallet for a user.

**Endpoint**: `POST /wallets`

**Request Body**: None 

**Response**:
```json
{
  "success": true,
  "message": "Wallet created successfully",
  "data": {
    "wallet_user_id": "5b3e5331-4dc1-49b8-85be-411c4885c0bd",
    "balances": {
    "Coins": 0,
    "Exp": 0
  },
  "created_at": "2025-09-08T09:32:17.849080675Z",
  "updated_at": "2025-09-08T09:32:17.849080675Z"
  },
  "timestamp": "2025-09-08T09:32:17.852732012Z"
}
```

**Status Codes**:
- `201 Created`: Wallet created successfully
- `409 Conflict`: Wallet already exists for user
- `500 Internal Server Error`: Unexpected error

### 2. Get Wallet

Retrieves wallet information by user ID.

**Endpoint**: `GET /wallets/{id}`

**Parameters**:
- `id` (path): Wallet User ID associated with the wallet

**Response**:
```json
{
  "success": true,
  "message": "Wallet retrieved successfully",
  "data": {
    "wallet_user_id": "5b3e5331-4dc1-49b8-85be-411c4885c0bd",
    "balances": {
    "Coins": 0,
    "Exp": 0
  },
  "created_at": "2025-09-08T09:32:17.849080675Z",
  "updated_at": "2025-09-08T09:32:17.849080675Z"
  },
  "timestamp": "2025-09-08T09:32:17.852732012Z"
}
```

**Status Codes**:
- `200 OK`: Wallet retrieved successfully
- `400 Bad Request`: Missing or invalid user ID
- `404 Not Found`: Wallet not found
- `500 Internal Server Error`: Unexpected error

### 3. Add Balance

Adds funds to a wallet.

**Endpoint**: `POST /wallets/{id}/add`

**Parameters**:
- `id` (path): Wallet User ID associated with the wallet

**Request Body**:
```json
{
  "type": "Coins",
  "amount": "500.5"
}
```

**Response**:
```json
{
 "success": true,
  "message": "Balance added successfully",
  "data": {
    "wallet_user_id": "5b3e5331-4dc1-49b8-85be-411c4885c0bd",
    "balances": {
    "Coins": 500.5,
    "Exp": 0
  },
  "created_at": "2025-09-08T09:32:17.849080675Z",
  "updated_at": "2025-09-08T09:32:17.849080675Z"
  },
  "timestamp": "2025-09-08T09:32:17.852732012Z"
}
```

**Status Codes**:
- `200 OK`: Balance added successfully
- `400 Bad Request`: Invalid request body or amount
- `404 Not Found`: Wallet not found
- `500 Internal Server Error`: Unexpected error

### 4. Deduct Balance

Deducts funds from a wallet.

**Endpoint**: `POST /wallets/{id}/deduct`

**Parameters**:
- `id` (path): Wallet User ID associated with the wallet

**Request Body**:
```json
{
  "type": "Coins",
  "amount": "300.5"
}
```

**Response**:
```json
{
 "success": true,
  "message": "Balance deducted successfully",
  "data": {
    "wallet_user_id": "5b3e5331-4dc1-49b8-85be-411c4885c0bd",
    "balances": {
    "Coins": 200,
    "Exp": 0
  },
  "created_at": "2025-09-08T09:32:17.849080675Z",
  "updated_at": "2025-09-08T09:32:17.849080675Z"
  },
  "timestamp": "2025-09-08T09:32:17.852732012Z"
}
```

**Status Codes**:
- `200 OK`: Balance deducted successfully
- `400 Bad Request`: Invalid request body, amount, or insufficient balance
- `404 Not Found`: Wallet not found
- `500 Internal Server Error`: Unexpected error

## Error Handling

The service implements comprehensive error handling with specific error codes:

### Error Response Format
```json
{
  "success": false,
  "message": "Error description",
  "error": "Additional error details (if applicable)"
}
```

### Error Codes
- `CodeWalletNotFound`: Wallet doesn't exist
- `CodeWalletExists`: Wallet already exists for user
- `CodeInsufficientBalance`: Not enough funds for deduction
- `CodeInvalidAmount`: Invalid amount specified
- `CodeInvalidBalanceType`: Invalid transaction type
- `CodeValidationError`: Request validation failed

## Deployment

### Prerequisites
- Docker and Docker Compose installed
- PostgreSQL 16 (handled by Docker)

### Environment Variables

Create a `.env` file in your project root:

```env
# Database Configuration
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgresql
DB_NAME=ecommerce_marketplace
DB_SSLMODE=disable

# Frappe Integration
FRAPPE_URL=http://ecommerce.local:8000
FRAPPE_API_KEY=a420e4791cb29de
FRAPPE_API_SECRET=55822b4d4ed62f8
```

### Service Configuration

The Docker Compose configuration includes:

- **App Service**: 
  - Port: 8080 (mapped to host)
  - Base image: golang:1.22-alpine

- **Database Service**:
  - Port: 5433 (mapped to host)
  - PostgreSQL 16

## Development

### Local Development Setup

1. **Install Go dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run with Docker**:
   ```bash
   docker-compose up
   ```
