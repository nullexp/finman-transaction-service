
```markdown
# Transaction Service

## Overview

The Transaction Service is a microservice responsible for managing financial transactions, including creating, updating, retrieving, and deleting transactions. It supports transaction types such as deposit and withdrawal, and provides functionalities like fetching transactions by user ID and with pagination.

## Features

- Create a transaction
- Get a transaction by ID
- Get all transactions
- Update a transaction
- Delete a transaction
- Get transactions by user ID
- Get transactions with pagination
- Get balance by user ID

## Prerequisites

- Docker
- Docker Compose

## Setup

1. Clone the repository:

```sh
git clone https://github.com/yourusername/transaction-service.git
cd transaction-service
```

2. Create a `.env` file in the root directory and set the following environment variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=transaction_db
PORT=8080
```

## Running the Service

To run the service, use Docker Compose:

```sh
docker-compose up
```

This command will build the Docker images and start the containers defined in the `docker-compose.yml` file.

## gRPC Service Definition

The Transaction Service uses gRPC for communication. Below is the service definition for the gRPC API:

```proto
// TransactionService service definition
service TransactionService {
    rpc CreateTransaction(CreateTransactionRequest) returns (CreateTransactionResponse);
    rpc GetTransactionById(GetTransactionByIdRequest) returns (GetTransactionByIdResponse);
    rpc GetTransactionsByUserId(GetTransactionsByUserIdRequest) returns (GetTransactionsByUserIdResponse);
    rpc GetOwnTransactionById(GetOwnTransactionByIdRequest) returns (GetOwnTransactionByIdResponse);
    rpc GetAllTransactions(GetAllTransactionsRequest) returns (GetAllTransactionsResponse);
    rpc UpdateTransaction(UpdateTransactionRequest) returns (UpdateTransactionResponse);
    rpc DeleteTransaction(DeleteTransactionRequest) returns (DeleteTransactionResponse);
    rpc GetTransactionsWithPagination(GetTransactionsWithPaginationRequest) returns (GetTransactionsWithPaginationResponse);
}
```

### Explanation of gRPC Methods

1. **CreateTransaction**: This method takes a `CreateTransactionRequest` and returns a `CreateTransactionResponse`. It is used to create a new transaction.

2. **GetTransactionById**: This method takes a `GetTransactionByIdRequest` and returns a `GetTransactionByIdResponse`. It is used to retrieve a transaction by its ID.

3. **GetTransactionsByUserId**: This method takes a `GetTransactionsByUserIdRequest` and returns a `GetTransactionsByUserIdResponse`. It is used to retrieve all transactions associated with a specific user ID.

4. **GetOwnTransactionById**: This method takes a `GetOwnTransactionByIdRequest` and returns a `GetOwnTransactionByIdResponse`. It is used to retrieve a transaction by its ID, ensuring the transaction belongs to the requesting user.

5. **GetAllTransactions**: This method takes a `GetAllTransactionsRequest` and returns a `GetAllTransactionsResponse`. It is used to retrieve all transactions.

6. **UpdateTransaction**: This method takes an `UpdateTransactionRequest` and returns an `UpdateTransactionResponse`. It is used to update an existing transaction.

7. **DeleteTransaction**: This method takes a `DeleteTransactionRequest` and returns a `DeleteTransactionResponse`. It is used to delete a transaction by its ID.

8. **GetTransactionsWithPagination**: This method takes a `GetTransactionsWithPaginationRequest` and returns a `GetTransactionsWithPaginationResponse`. It is used to retrieve transactions with pagination support.
