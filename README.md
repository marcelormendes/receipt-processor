# Receipt Processor

A Go microservice that processes receipt data and calculates reward points based on specific rules. Built with modern Go practices including strong typing, context support, and structured logging.

## Description

This application provides a RESTful API for processing receipts and calculating points based on predefined rules. It uses in-memory storage for simplicity, so all data is lost when the service restarts. The application follows idiomatic Go patterns and best practices.

## Features

- RESTful API for receipt processing and point calculation
- Strong type validation for dates, times, and prices
- Context support for request cancellation and timeouts
- Structured logging with log levels
- Thread-safe storage implementation
- Graceful shutdown
- Environment-based configuration
- Health check endpoint
- Comprehensive test coverage

## Prerequisites

- Docker (for containerized deployment)
- Go 1.21+ (optional, for local development)

## Building and Running

### Using Docker (Recommended)

Build the Docker image:

```bash
docker build -t receipt-processor .
```

Run the container:

```bash
docker run -p 8080:8080 receipt-processor
```

With custom environment variables:

```bash
docker run -p 8080:8080 \
  -e PORT=8080 \
  -e LOG_LEVEL=DEBUG \
  -e GIN_MODE=release \
  receipt-processor
```

The service will be available at http://localhost:8080

### Local Development

If you prefer to run the application locally:

```bash
# Build the application
go build -o receipt-processor

# Run the application
./receipt-processor

# Run with environment variables
PORT=3000 LOG_LEVEL=DEBUG ./receipt-processor
```

## Configuration

The application can be configured via environment variables:

| Variable  | Description                              | Default |
| --------- | ---------------------------------------- | ------- |
| PORT      | Port to run the server on                | 8080    |
| LOG_LEVEL | Logging level (DEBUG, INFO, WARN, ERROR) | INFO    |
| GIN_MODE  | Gin mode (debug, release, test)          | debug   |

## API Endpoints

The service provides the following endpoints:

### 1. Process a Receipt

```
POST /receipts/process
```

Processes a receipt and returns a unique ID for later reference.

**Request Body:**

```json
{
  "retailer": "String",
  "purchaseDate": "YYYY-MM-DD",
  "purchaseTime": "HH:MM",
  "items": [
    {
      "shortDescription": "String",
      "price": "String (decimal)"
    }
  ],
  "total": "String (decimal)"
}
```

**Response:**

```json
{
  "id": "UUID string"
}
```

**Status Codes:**

- `200 OK`: Receipt processed successfully
- `400 Bad Request`: Invalid receipt data
- `500 Internal Server Error`: Processing error

### 2. Get Points for a Receipt

```
GET /receipts/{id}/points
```

Returns the points calculated for a previously processed receipt.

**Response:**

```json
{
  "points": integer
}
```

**Status Codes:**

- `200 OK`: Points retrieved successfully
- `404 Not Found`: Receipt ID not found
- `500 Internal Server Error`: Processing error

### 3. Health Check

```
GET /health
```

Returns the health status of the service.

**Response:**

```json
{
  "status": "ok"
}
```

**Status Codes:**

- `200 OK`: Service is healthy

## Example Usage

### Example 1: Target Receipt

Process the receipt:

```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {
        "shortDescription": "Mountain Dew 12PK",
        "price": "6.49"
      },
      {
        "shortDescription": "Emils Cheese Pizza",
        "price": "12.25"
      },
      {
        "shortDescription": "Knorr Creamy Chicken",
        "price": "1.26"
      },
      {
        "shortDescription": "Doritos Nacho Cheese",
        "price": "3.35"
      },
      {
        "shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
        "price": "12.00"
      }
    ],
    "total": "35.35"
  }'
```

Sample response:

```json
{
  "id": "a9f2e68d-f6fc-4b87-9122-56ff11f06981"
}
```

Get the points (use the ID from the previous response):

```bash
curl http://localhost:8080/receipts/a9f2e68d-f6fc-4b87-9122-56ff11f06981/points
```

Sample response:

```json
{
  "points": 28
}
```

### Example 2: M&M Corner Market Receipt

Process the receipt:

```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "M&M Corner Market",
    "purchaseDate": "2022-03-20",
    "purchaseTime": "14:33",
    "items": [
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      },
      {
        "shortDescription": "Gatorade",
        "price": "2.25"
      }
    ],
    "total": "9.00"
  }'
```

Get the points (use the ID from the response):

```bash
curl http://localhost:8080/receipts/{id}/points
```

## Point Calculation Rules

1. One point for every alphanumeric character in the retailer name
2. 50 points if the total is a round dollar amount with no cents
3. 25 points if the total is a multiple of 0.25
4. 5 points for every two items on the receipt
5. If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned
6. 6 points if the purchase day is odd
7. 10 points if the purchase time is between 14:00 and 16:00 (exclusive)

## Testing

To run the tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with code coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...
```

## Code Structure

The service is organized into the following packages:

- `api`: HTTP handlers and routing
- `models`: Data structures and validation
- `services`: Business logic including point calculation
- `storage`: Data persistence (in-memory implementation)
- `tests`: End-to-end tests

## Design Decisions

- **Custom Types**: Used for dates, times, and prices to ensure strong validation and type safety
- **Interface-Based Design**: Storage is defined via interfaces for better testability and extensibility
- **Context Support**: All operations support context for cancellation and timeouts
- **Structured Logging**: Using Go's standard `log/slog` package for structured, leveled logging
- **Thread Safety**: All shared state is protected with appropriate synchronization

## Future Improvements

- Database persistence (e.g., PostgreSQL, MongoDB)
- Authentication and authorization
- Rate limiting
- API versioning
- Metrics and monitoring
- Distributed tracing

## License

MIT
