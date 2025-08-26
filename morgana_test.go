package morgana_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bi0dread/morgana"
	"github.com/stretchr/testify/assert"
)

func TestMorgana(t *testing.T) {
	m := morgana.New("TestType")

	t.Run("WithStatusCode", func(t *testing.T) {
		m = m.WithStatusCode(http.StatusOK)
		assert.Equal(t, http.StatusOK, m.GetStatusCode())
	})

	t.Run("WithCustomCode", func(t *testing.T) {
		m = m.WithCustomCode("CUSTOM_CODE")
		assert.Equal(t, "CUSTOM_CODE", m.GetCustomCode())
	})

	t.Run("WithType", func(t *testing.T) {
		m = m.WithType("NewType")
		assert.Equal(t, "NewType", m.GetType())
	})

	t.Run("WithMessage", func(t *testing.T) {
		m = m.WithMessage("Test message")
		assert.Equal(t, "Test message", m.GetMessage())
	})

	t.Run("With", func(t *testing.T) {
		m = m.With("WithValue")
		assert.Equal(t, "WithValue", m.GetWith())
	})

	t.Run("WithError", func(t *testing.T) {
		err := errors.New("Test error")
		m = m.WithError(err)
		stack := m.GetMorganaStackErrors()
		assert.NotEmpty(t, stack)
	})

	t.Run("WithAddMetaDataKey", func(t *testing.T) {
		m = m.WithAddMetaDataKey("key", "value")
		assert.Equal(t, "value", m.GetMetaDataKey("key"))
	})

	t.Run("GetMetaData", func(t *testing.T) {
		metaData := m.GetMetaData()
		assert.NotNil(t, metaData)
		assert.Equal(t, "value", metaData["key"])
	})

	t.Run("Clone", func(t *testing.T) {
		clone := m.Clone(1)
		assert.Equal(t, m.GetType(), clone.GetType())
		assert.Equal(t, m.GetCustomCode(), clone.GetCustomCode())
	})

	t.Run("ToError", func(t *testing.T) {
		err := m.ToError()
		assert.NotNil(t, err)
	})

	t.Run("String", func(t *testing.T) {
		str := m.String()
		assert.NotEmpty(t, str)
	})

	t.Run("ToJson", func(t *testing.T) {
		jsonStr := m.ToJson()
		assert.NotEmpty(t, jsonStr)
	})

	// New features
	t.Run("FullStackCapture", func(t *testing.T) {
		mm := morgana.New("Stack").WithFullStack(2, 8)
		frames := mm.GetStackFrames()
		assert.NotEmpty(t, frames)
	})

	t.Run("SafeJSONRedaction", func(t *testing.T) {
		mm := morgana.New("Secret").WithAddMetaDataKey("token", "abc").WithRedactedKey("token")
		out := mm.ToJsonSafe()
		assert.Contains(t, out, "[REDACTED]")
	})

	t.Run("HTTPWriter", func(t *testing.T) {
		rec := httptest.NewRecorder()
		mm := morgana.New("Unauthorized").WithStatusCode(http.StatusUnauthorized).WithMessage("nope")
		mm.WriteHTTP(rec, true)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "nope")
	})

	t.Run("PanicRecover", func(t *testing.T) {
		err := morgana.Recover(func() { panic("boom") })
		assert.NotNil(t, err)
	})

	t.Run("JoinHelper", func(t *testing.T) {
		err := morgana.Join(errors.New("first"), morgana.New("Second").WithMessage("second").ToError())
		assert.NotNil(t, err)
	})

	t.Run("GRPCCodeMapping", func(t *testing.T) {
		mm := morgana.New("NotFound").WithStatusCode(http.StatusNotFound)
		code := mm.ToGRPCCode()
		assert.Equal(t, 5, code) // NotFound
		mm = mm.FromGRPCCode(3)  // InvalidArgument
		assert.Equal(t, http.StatusBadRequest, mm.GetStatusCode())
	})

	t.Run("WithTraceAndFields", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "trace_id", "tid-123")
		mm := morgana.New("Trace").WithTrace(ctx)
		fields := mm.ToFields()
		assert.Equal(t, "tid-123", mm.GetMetaData()["trace_id"])
		assert.NotEmpty(t, fields["id"]) // has correlation id
	})

	t.Run("WithIDAndFieldErrors", func(t *testing.T) {
		mm := morgana.New("Val").WithID("id-1").WithFieldError("name", "required", "missing")
		assert.Equal(t, "id-1", mm.GetID())
		fes := mm.GetFieldErrors()
		assert.Equal(t, 1, len(fes))
		assert.Equal(t, "name", fes[0].Field)
	})

	t.Run("CauseAndUnwrap", func(t *testing.T) {
		cause := errors.New("root cause")
		mm := morgana.New("Cause").WithCause(cause)
		err := mm.ToError()
		assert.NotNil(t, err)
	})
}

func TestWrap(t *testing.T) {
	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")

	t.Run("Wrap both errors", func(t *testing.T) {
		wrappedErr := morgana.Wrap(err1, err2)
		assert.NotNil(t, wrappedErr)
	})

	t.Run("Wrap first error nil", func(t *testing.T) {
		wrappedErr := morgana.Wrap(nil, err2)
		assert.Equal(t, err2, wrappedErr)
	})

	t.Run("Wrap second error nil", func(t *testing.T) {
		wrappedErr := morgana.Wrap(err1, nil)
		assert.Equal(t, err1, wrappedErr)
	})
}

func TestGetStringDetail(t *testing.T) {
	err := errors.New("Test error")
	detail := morgana.GetStringDetail(err)
	assert.Equal(t, "Test error", detail)
}
