
# Morgana - Enhanced Error Management in Go

Morgana is a robust and flexible error management library for Go, designed to enhance error handling and debugging by providing rich metadata, stack traces, and structured error objects.

---

## Features

- Rich metadata support for errors.
- Stack trace generation.
- Custom error codes and types.
- JSON serialization of errors.
- Easy wrapping and chaining of errors.
- Flexible interface for creating and handling errors.

---

## Installation

```bash
go get -u github.com/bi0dread/morgana
```

---

## Quick Start

### Creating a Basic Error

```go
package main

import (
	"fmt"
	"github.com/bi0dread/morgana"
)

func main() {
	err := morgana.New("ValidationError").
		WithStatusCode(400).
		WithCustomCode("INVALID_INPUT").
		WithMessage("Invalid input provided for field: %s", "email")

	fmt.Println(err.String())
}
```

### Adding Internal Details

```go
err := morgana.New("DatabaseError").
	WithInternalDetail("Failed to connect to the database", "retry in 5 seconds")

fmt.Println(err.String())
```

### Attaching Metadata

```go
err := morgana.New("NetworkError").
	WithAddMetaDataKey("endpoint", "/api/v1/resource").
	WithAddMetaDataKey("method", "GET")

fmt.Println(err.GetMetaData())
```

### Wrapping Errors

```go
originalErr := fmt.Errorf("connection timeout")
wrappedErr := morgana.Wrap(originalErr, morgana.New("NetworkError"))

fmt.Println(morgana.GetStringDetail(wrappedErr))
```

### Cloning Errors

```go
err := morgana.New("FileError").
	WithMessage("File not found: %s", "config.yaml")

clonedErr := err.Clone(2)
fmt.Println(clonedErr.String())
```

### Converting to JSON

```go
err := morgana.New("AuthorizationError").
	WithMessage("User not authorized for this action")

fmt.Println(err.ToJson())
```

### Handling Errors with Morgana

```go
func handleError(err error) {
	morgErr := morgana.FromError(err)
	if morgErr != nil {
		fmt.Println("Custom Code:", morgErr.GetCustomCode())
		fmt.Println("Message:", morgErr.GetMessage())
	} else {
		fmt.Println("Standard Error:", err)
	}
}
```

---

## API Reference

### Morgana Interface

#### Creation

- `New(typeValue string) Morgana`: Creates a new Morgana instance with the specified type.

#### Chainable Methods

- `WithStatusCode(statusCode int) Morgana`: Sets the HTTP status code.
- `WithCustomCode(customCode string) Morgana`: Adds a custom error code.
- `WithType(ref string) Morgana`: Sets the error type.
- `WithMessage(msg string, args ...string) Morgana`: Adds a formatted error message.
- `WithInternalDetail(detail ...any) Morgana`: Adds internal details.
- `WithStackTrace(skip int) Morgana`: Attaches a stack trace.
- `With(value string) Morgana`: Adds an additional string value.
- `WithError(err error) Morgana`: Wraps another error.
- `WithAddMetaDataKey(key string, value any) Morgana`: Adds metadata.

#### Retrieval Methods

- `GetStatusCode() int`: Retrieves the HTTP status code.
- `GetCustomCode() string`: Gets the custom error code.
- `GetType() string`: Returns the error type.
- `GetMessage() string`: Retrieves the error message.
- `GetInternalDetail() []any`: Retrieves internal details.
- `GetWith() string`: Gets the additional string value.
- `GetMetaDataKey(key string) any`: Retrieves a specific metadata value.
- `GetMetaData() map[string]any`: Retrieves all metadata.
- `GetMorganaStackErrors() []Morgana`: Gets wrapped errors.

#### Utility Methods

- `Clone(stackLevel int) Morgana`: Clones the error with a new stack trace level.
- `ToError() error`: Converts the Morgana instance to a standard error.
- `String() string`: Converts the error to a string.
- `ToJson() string`: Serializes the error to JSON.

---

## Best Practices

1. Use `WithInternalDetail` for debugging information.
2. Attach meaningful `CustomCode` values for better error categorization.
3. Always add stack traces (`WithStackTrace`) for production debugging.
4. Wrap errors to preserve the original context and chain multiple errors.
5. Use `GetMorgana` to extract enriched error information from standard errors.

---


