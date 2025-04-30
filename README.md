# Kong API Explorer

A tool for exploring and testing Kong API Gateway routes directly from the Kong database.

## Preview

![Preview](https://raw.githubusercontent.com/sahapranta/kong-check/refs/heads/main/assets/demo.gif)

## Features

- List all routes and services in Kong with colorized output
- Check route health and connectivity with custom paths
- Concurrent route checking with configurable timeout
- Support for headers, methods, and hosts from Kong configuration
- Verbose mode for viewing response details

## Installation

1. Clone this repository
2. Install dependencies: `go mod download`
3. Build the application: `go build -o kong-explorer`

## Configuration

Create a `.env` file in the root directory (see `example.env` for reference):

```
# Database configuration
PG_HOST=localhost
PG_PORT=5432
PG_USER=kong
PG_PASSWORD=kong
PG_DATABASE=kong
PG_SSLMODE=disable

# Default settings for checks
DEFAULT_PROTOCOL=http
DEFAULT_HOSTNAME=localhost:8000
```

## Usage

### Listing Routes

```bash
# List all routes with default options
./kong-explorer list

# Show headers for each route
./kong-explorer list --headers
```

### Checking Routes

```bash
# Check all routes
./kong-explorer check --all

# Check specific services
./kong-explorer check --services auth-service,user-service

# Check with a custom path appended to routes
./kong-explorer check --all --path /health

# Check with specific HTTP methods
./kong-explorer check --all --methods GET,POST,PUT

# Verbose output with response details
./kong-explorer check --all -v
```

### Help

```bash
./kong-explorer help
```

## Project Structure

- `cmd/`: Command implementations
- `config/`: Configuration management
- `db/`: Database interaction
- `models/`: Data structures
- `utils/`: Utility functions

## Example Output

### List Command

```
Kong API Routes:

Service: authentication-service
  Route: auth-routes
    Path: /auth
    Methods: GET, POST, PUT

Service: user-service
  Route: user-routes
    Path: /users
    Methods: GET, POST, DELETE

Total: 2 routes across 2 services
```

### Check Command

```
Checking 2 routes (timeout: 5s)...

[GET] authentication-service auth-routes /auth - OK (200) - 145ms
[GET] user-service user-routes /users - OK (200) - 132ms
[GET] user-service user-routes /users - OK (201) - 189ms

Check completed.
```