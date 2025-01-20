package morgana_test

import (
	"errors"
	"net/http"
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

	t.Run("WithInternalDetail", func(t *testing.T) {
		m = m.WithInternalDetail("Detail1", "Detail2")
		assert.Equal(t, []any{"Detail1", "Detail2"}, m.GetInternalDetail())
	})

	t.Run("With", func(t *testing.T) {
		m = m.With("WithValue")
		assert.Equal(t, "WithValue", m.GetWith())
	})

	t.Run("WithError", func(t *testing.T) {
		err := errors.New("Test error")
		m = m.WithError(err)
		assert.Contains(t, m.GetInternalDetail(), "Test error")
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
