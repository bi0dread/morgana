
# Morgana - Enhanced Error Management in Go

Morgana is a robust and flexible error management library for Go, designed to enhance error handling and debugging by providing rich metadata, stack traces, and structured error objects.

---

## Features

- Rich metadata support for errors.
- Stack trace generation (single frame and full call stack).
- Custom error codes and types.
- JSON and Safe JSON serialization (with redaction).
- Easy wrapping, chaining, cause tracking, and joining of errors.
- HTTP response writer helper.
- Logger fields extraction for structured logs.
- Panic capture helpers.
- Context trace enrichment.
- gRPC status code mapping helpers.

---

## Installation

```bash
go get -u github.com/bi0dread/morgana
```

---

## Quick Start

### Creating a Basic Error

```go
err := morgana.New("ValidationError").
	WithStatusCode(400).
	WithCustomCode("INVALID_INPUT").
	WithMessage("Invalid input provided for field: %s", "email")
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
wrappedErr := morgana.Wrap(originalErr, morgana.New("NetworkError").ToError())
fmt.Println(morgana.GetStringDetail(wrappedErr))
```

### Stack Errors and Cause

```go
cause := fmt.Errorf("disk full")
err := morgana.New("IOError").WithMessage("failed writing file").WithCause(cause)
fmt.Println(err.Cause()) // root cause: disk full

other := fmt.Errorf("low-level error")
err = err.WithError(other) // pushes a stack Morgana from other
fmt.Println(err.GetMorganaStackErrors())
```

### Full Stack Capture

```go
err := morgana.New("FileError").
	WithMessage("File not found: %s", "config.yaml").
	WithFullStack(2, 64)
fmt.Println(err.GetStackFrames())
```

### Safe JSON and Redaction

```go
err := morgana.New("SecretError").
	WithAddMetaDataKey("token", "super-secret").
	WithRedactedKey("token")
fmt.Println(err.ToJsonSafe()) // token value redacted
```

### HTTP Writer Helper

```go
func handler(w http.ResponseWriter, r *http.Request) {
	m := morgana.New("Unauthorized").WithStatusCode(http.StatusUnauthorized).WithMessage("auth required")
	m.WriteHTTP(w, true) // write safe JSON
}
```

### Panic Capture

```go
err := morgana.Recover(func(){
	panic("boom")
})
if err != nil {
	fmt.Println(morgana.GetStringDetail(err))
}
```

### Join Multiple Errors

```go
err := morgana.Join(
	fmt.Errorf("first"),
	morgana.New("Second").WithMessage("second").ToError(),
)
fmt.Println(morgana.GetStringDetail(err))
```

### gRPC Helpers (code mapping)

```go
m := morgana.New("NotFound").WithStatusCode(http.StatusNotFound).WithMessage("missing")
code := m.ToGRPCCode() // 5 (NotFound)
_ = code

m.FromGRPCCode(3) // InvalidArgument -> 400
fmt.Println(m.GetStatusCode())
```

---

## API Notes

- `WithTrace(ctx)` pulls common correlation IDs from context.
- `WithID/ GetID` provides an error correlation ID.
- `WithFieldError/ GetFieldErrors` helps shape validation errors (HTTP 422 style).
- `ToFields()` returns structured fields for logging.
- `Empo` implements `Unwrap()` and can carry a `cause` for standard error traversal.

---


