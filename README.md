# Fintech Ledger Microservice

A high-performance, double-entry bookkeeping ledger microservice built with Go. This project demonstrates a production-ready approach to handling financial transactions with strict ACID compliance, concurrency safety, and clean architecture.

## 🚀 Key Features

- **Double-Entry Bookkeeping**: Ensures that every transaction has a corresponding debit and credit, maintaining a balanced ledger at all times.
- **ACID Compliance**: Leverages PostgreSQL transactions to guarantee data integrity.
- **Concurrency Safety**: Uses `SELECT FOR UPDATE` to prevent race conditions during concurrent balance updates.
- **Clean Architecture**: Structured into distinct layers (API, Service, Repository) for maintainability and testability.
- **Type-Safe SQL**: Utilizes `sqlc` to generate type-safe Go code from SQL queries, eliminating runtime SQL errors.
- **Containerized**: Fully Dockerized application with Docker Compose for easy orchestration.

## 🛠️ Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **Web Framework**: [Echo v4](https://echo.labstack.com/) - High performance, extensible, minimalist Go web framework.
- **Database**: [PostgreSQL](https://www.postgresql.org/) - The World's Most Advanced Open Source Relational Database.
- **Database Driver**: [lib/pq](https://github.com/lib/pq)
- **SQL Generator**: [sqlc](https://sqlc.dev/) - Compile SQL to type-safe Go.
- **Infrastructure**: Docker, Docker Compose.

## 📂 Architecture

The project follows a **Clean Architecture** pattern to separate concerns:

- `cmd/server`: Entry point of the application.
- `internal/api`: HTTP transport layer (handlers, routing) using Echo.
- `internal/service`: Business logic layer.
- `internal/repository`: Data access layer (interacting with the database).
- `internal/db`: Generated Go code from SQL queries (via `sqlc`).
- `migrations`: Database schema migrations.

## 🏁 Getting Started

### Prerequisites

- [Go 1.20+](https://go.dev/dl/)
- [Docker](https://www.docker.com/products/docker-desktop) & Docker Compose

### Running with Docker (Recommended)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/yourusername/ledger-microservice.git
    cd ledger-microservice
    ```

2.  **Start the services:**
    ```bash
    docker-compose up --build
    ```
    This will start both the PostgreSQL database and the Go API server. The API will be available at `http://localhost:8080`.

### Running Locally

1.  **Start PostgreSQL:**
    Ensure you have a PostgreSQL instance running. You can use the docker-compose file to just start the db:
    ```bash
    docker-compose up -d postgres
    ```

2.  **Run Migrations:**
    (Optional) If you have `golang-migrate` installed, run migrations. Otherwise, the docker setup handles initialization via `schema.sql`.

3.  **Run the Application:**
    ```bash
    export DB_SOURCE="postgresql://root:secret@localhost:5432/ledger?sslmode=disable"
    go run cmd/server/main.go
    ```

## 🔌 API Reference

### Create Account
Create a new account with an initial balance.

- **URL**: `/accounts`
- **Method**: `POST`
- **Body**:
    ```json
    {
        "owner": "John Doe",
        "currency": "USD"
    }
    ```

### Get Account
Retrieve account details.

- **URL**: `/accounts/:id`
- **Method**: `GET`

### Create Transfer
Execute a money transfer between two accounts.

- **URL**: `/transactions`
- **Method**: `POST`
- **Body**:
    ```json
    {
        "from_account_id": 1,
        "to_account_id": 2,
        "amount": 100,
        "currency": "USD"
    }
    ```



## 📄 License

This project is licensed under the MIT License.
