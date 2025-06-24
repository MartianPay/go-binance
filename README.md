# Go-Binance SDK

A Go SDK for Binance Capital Account APIs, providing easy access to deposit, withdrawal, and account management endpoints.

## Installation

```bash
go get github.com/MartianPay/go-binance
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/MartianPay/go-binance"
    "github.com/MartianPay/go-binance/models"
)

func main() {
    client := binance.NewClient("your-api-key", "your-secret-key")
    
    // Enable fast withdraw for instant internal transfers
    err := client.Account.EnableFastWithdrawSwitch(0)
    if err != nil {
        log.Fatal(err)
    }
    
    // Get deposit address
    addr, err := client.Deposit.GetDepositAddress(models.DepositAddressRequest{
        Coin: "BTC",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("BTC Address: %s\n", addr.Address)
}
```

## Features

### Deposit Operations
- Get deposit address
- Get deposit history

### Withdrawal Operations  
- Submit withdrawal request
- Get withdrawal history
- Get withdrawal quota

### Account Management
- Get all coins information
- Get account information
- Universal asset transfer
- Get user assets
- Enable fast withdraw switch (for instant internal transfers)

## Authentication

The SDK uses HMAC SHA256 authentication. You need to provide your API key and secret key when creating the client.

## Examples

See the example files in the `examples/` directory for usage examples.

### Running the Examples

There are two example programs available:

#### 1. Account Information Example

This example demonstrates how to get account information and enable fast withdraw:

```bash
# Navigate to the project directory
cd /path/to/go-binance

# Set your API credentials (replace with your actual keys)
export BINANCE_API_KEY="your-api-key"
export BINANCE_SECRET_KEY="your-secret-key"

# Run the account information example
go run examples/account_info.go
```

#### 2. Withdrawal Example

This example demonstrates deposit and withdrawal operations:

```bash
# Navigate to the project directory
cd /path/to/go-binance

# Set your API credentials (replace with your actual keys)
export BINANCE_API_KEY="your-api-key"
export BINANCE_SECRET_KEY="your-secret-key"

# Run the withdrawal example
go run examples/withdrawal.go

# Or if you want to use testnet
export BINANCE_TESTNET=true
go run examples/withdrawal.go
```

The withdrawal example demonstrates:
- Getting deposit addresses for multiple networks
- Fetching deposit history
- Checking withdrawal quotas
- Submitting withdrawals (commented out by default for safety)
- Retrieving withdrawal history

## API Documentation

For detailed API documentation, visit: https://developers.binance.com/docs/wallet