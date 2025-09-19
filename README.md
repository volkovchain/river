# River - Blockchain Salary Payment Automation

River is a Go-based application that automates salary payments on the Ethereum blockchain. It simplifies the process of distributing recurring payments to employees using their blockchain wallets, with built-in support for ERC-20 tokens like USDC.

## Table of Contents

- [Features](#features)
- [How It Works](#how-it-works)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
  - [Adding Employees](#adding-employees)
  - [Processing Payments](#processing-payments)
  - [Repayment](#repayment)
- [Database Schema](#database-schema)
- [Development](#development)
  - [Building](#building)
  - [Testing](#testing)
- [Architecture](#architecture)
- [Security](#security)

## Features

- **Automated Salary Payments**: Automatically distribute salaries to employees via blockchain transactions
- **ERC-20 Token Support**: Built-in support for tokens like USDC with 6 decimal precision
- **Batch Processing**: Process multiple salary payments in a single execution
- **Retry Mechanism**: Handle failed transactions with repayment functionality
- **SQLite Database**: Local storage for employee data and payment history
- **Configuration Management**: Flexible configuration via environment variables and .env files
- **CLI Interface**: Simple command-line interface for easy automation

## How It Works

1. Employee data is stored in a local SQLite database
2. When executed, River creates a new salary record and processes payments to all employees
3. For each employee, River constructs an ERC-20 token transfer transaction
4. Transactions are signed with configured private keys and broadcasted to the Ethereum network
5. Payment status is tracked and stored in the database

## Prerequisites

- Go 1.19 or higher
- Access to an Ethereum node (Infura, Alchemy, or local node)
- ERC-20 tokens (e.g., USDC) in the wallet corresponding to the private key
- SQLite3

## Installation

```bash
# Clone the repository
git clone <repository-url>
cd river

# Install dependencies
go mod download

# Build the application
go build -o river
```

## Configuration

River can be configured using environment variables or .env files. Create a `.env` file in the project root:

```env
# Ethereum node URL (required)
NODE=https://mainnet.infura.io/v3/YOUR_PROJECT_ID

# Private keys for signing transactions (required)
# Multiple keys can be specified separated by commas
PRIVATE_KEYS=YOUR_PRIVATE_KEY

# Database path (optional, defaults to ./main.db)
DATABASE_PATH=./main.db
```

Alternatively, you can use a `main.env` file with the same format.

## Usage

### Adding Employees

Employees are stored in the `employers` table of the SQLite database. Add employees using SQL:

```sql
INSERT INTO employers (name, addr, amount_salary) VALUES ('John Doe', '0xWalletAddress', 1000000);
```

**Note on Amounts**: The `amount_salary` field represents the smallest unit of the token. For USDC (6 decimals), to send 1 USDC, you would specify 1000000 (1 * 10^6).

### Processing Payments

To process salary payments:

```bash
./river
```

The application will prompt for confirmation before executing payments.

### Repayment

To retry failed payments or process payments that were interrupted:

```bash
./river repay
```

## Database Schema

River uses a SQLite database with the following tables:

- `employers`: Employee information (name, wallet address, salary amount)
- `salaries`: Salary records with status tracking
- `payments`: Individual payment records with transaction details

## Development

### Building

```bash
# Build the application
go build -o river

# Build with specific flags
go build -ldflags "-s -w" -o river
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Architecture

River follows a clean architecture with the following components:

- **CLI Layer** (`cmd/`): Command-line interface using Cobra
- **Business Logic** (`internal/service/`): Core salary and payment processing logic
- **Blockchain Integration** (`internal/client/ethereum/`): Ethereum client implementation
- **Data Access** (`db/`): Database repositories for employees, salaries, and payments
- **Configuration** (`internal/config/`): Configuration management using Viper
- **Entities** (`internal/entity/`): Domain models

## Security

- **Private Keys**: Never commit private keys to version control. Use environment variables or secure vaults.
- **Transaction Verification**: Always verify transaction details before confirming payments.
- **Network Security**: Use secure connections to Ethereum nodes (HTTPS/WSS).
- **Access Control**: Restrict access to the application and database files.