# FinMan User Service

The FinMan User Service is a microservice responsible for managing user information in the FinMan system. It includes functionalities for creating, retrieving, updating, and deleting users, as well as handling user authentication and pagination.

## Table of Contents
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Environment Variables](#environment-variables)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)

## Features
- Create a new user
- Retrieve user by ID
- Retrieve all users
- Update user details
- Delete a user
- Get user by username and password
- Get users with pagination
- Handle user authentication
- Database migrations using [golang-migrate](https://github.com/golang-migrate/migrate)
- JWT authentication
- Password hashing with bcrypt

## Requirements
- Go 1.18+
- PostgreSQL
- Docker (for containerized deployment)

## Installation
1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/finman-transaction-service.git
    cd finman-transaction-service
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Set up the environment variables in a `.env` file (see [Environment Variables](#environment-variables) section).

4. Run the database migrations:
    ```bash
    migrate -path ./migrations -database postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable up
    ```

## Usage
### Running the Service
1. Build the service:
    ```bash
    go build -o user-service .
    ```

2. Run the service:
    ```bash
    ./user-service
    ```

### Using Docker
1. Build the Docker image:
    ```bash
    docker build -t finman-transaction-service .
    ```

2. Run the Docker container:
    ```bash
   docker-compose up --build 
    ```

## Environment Variables
Create a `.env` file in the root directory with the following content:
  ```bash
  DB_HOST=localhost
  DB_PORT=5432
  DB_USER=postgres
  DB_PASSWORD=postgres
  DB_NAME=finman-transaction
  PORT=8082
  IP=0.0.0.0
      ```