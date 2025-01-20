package morgana

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"runtime"
	"strings"
)

const (
	morgana_key_data = "morgana_key_data"
)

type Morgana interface {
	WithStatusCode(statusCode int) Morgana
	WithCustomCode(customCode string) Morgana
	WithType(ref string) Morgana
	WithMessage(msg string, args ...string) Morgana
	WithInternalDetail(detail ...any) Morgana
	WithStackTrace(skip int) Morgana
	With(value string) Morgana
	WithError(err error) Morgana
	WithAddMetaDataKey(key string, value any) Morgana

	GetStatusCode() int
	GetCustomCode() string
	GetType() string
	GetMessage() string
	GetInternalDetail() []any
	GetWith() string
	GetMetaDataKey(key string) any
	GetMetaData() map[string]any
	GetMorganaStackErrors() []Morgana

	Clone(stackLevel int) Morgana

	ToError() error

	String() string
	ToJson() string
}
type morgana struct {
	Type               string
	WithValue          string
	Msg                string
	StatusCode         int
	CustomCode         string
	StackTrace         string
	InternalDetail     []any
	morganaStackErrors []Morgana
	MetaData           map[string]any
}

func (m *morgana) GetMorganaStackErrors() []Morgana {
	return m.morganaStackErrors
}

func (m *morgana) WithAddMetaDataKey(key string, value any) Morgana {
	if m.MetaData == nil {
		m.MetaData = make(map[string]any)
	}
	m.MetaData[key] = value
	return m
}

func (m *morgana) GetMetaDataKey(key string) any {
	if m.MetaData == nil {
		m.MetaData = make(map[string]any)
		return ""
	}

	if val, ok := m.MetaData[key]; ok {
		return val
	}

	return ""
}

func (m *morgana) GetMetaData() map[string]any {
	return m.MetaData
}

func (m *morgana) String() string {

	var builder strings.Builder

	if len(m.Type) != 0 {
		builder.WriteString(fmt.Sprintf("Type: %s", m.Type))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.WithValue) != 0 {
		builder.WriteString(fmt.Sprintf("With: %v", m.WithValue))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.Msg) != 0 {
		builder.WriteString(fmt.Sprintf("Msg: %s", m.Msg))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.CustomCode) != 0 {
		builder.WriteString(fmt.Sprintf("CustomCode: %s", m.CustomCode))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	builder.WriteString(fmt.Sprintf("StatusCode: %d", m.StatusCode))
	builder.WriteString(fmt.Sprintf("%v", " , "))

	if len(m.StackTrace) != 0 {
		builder.WriteString(fmt.Sprintf("StackTrace: %s", m.StackTrace))
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	if m.InternalDetail != nil && len(m.InternalDetail) != 0 {
		builder.WriteString(fmt.Sprintf("InternalDetail: %+#v\n", m.InternalDetail))
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	if m.morganaStackErrors != nil && len(m.morganaStackErrors) != 0 {
		builder.WriteString(fmt.Sprintf("MorganaStackErrors: \n"))
		builder.WriteString(fmt.Sprintf("%v", m.morganaStackErrors))

	}

	if m.MetaData != nil && len(m.MetaData) != 0 {
		builder.WriteString(fmt.Sprintf("MetaData: %+#v\n", m.MetaData))
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	return builder.String()

}

func (m *morgana) ToJson() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(bytes)

}

func (m *morgana) stringSimple() string {
	var builder strings.Builder

	if len(m.Type) != 0 {
		builder.WriteString(fmt.Sprintf("Type: %s", m.Type))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.WithValue) != 0 {
		builder.WriteString(fmt.Sprintf("With: %s", m.WithValue))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.Msg) != 0 {
		builder.WriteString(fmt.Sprintf("Msg: %s", m.Msg))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if len(m.CustomCode) != 0 {
		builder.WriteString(fmt.Sprintf("CustomCode: %s", m.CustomCode))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	if m.MetaData != nil && len(m.MetaData) != 0 {
		builder.WriteString(fmt.Sprintf("MetaData: %v\n", m.MetaData))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	builder.WriteString(fmt.Sprintf("StatusCode: %d", m.StatusCode))
	builder.WriteString(fmt.Sprintf("%v", " , "))

	return builder.String()
}

func (m *morgana) ToError() error {

	return NewEmpo(m.stringSimple()).WithAttributes(map[string]any{morgana_key_data: m}).ToError()
}

func FromError(err error) Morgana {
	if err == nil {
		return nil
	}

	empoData := GetEmpo(err)
	if empoData != nil {

		if val, ok := empoData.GetAttributes()[morgana_key_data]; ok {

			if morgana, ok := val.(Morgana); ok {
				return morgana
			}
		}
	}

	return New("GENERAL").WithStatusCode(http.StatusNotImplemented).WithInternalDetail(err)
}

func GetMorgana(err error) Morgana {
	if err == nil {
		return nil
	}

	empoData := GetEmpo(err)
	if empoData != nil {

		if val, ok := empoData.GetAttributes()[morgana_key_data]; ok {

			if morgana, ok := val.(Morgana); ok {
				return morgana
			}
		}
	}

	return nil
}

func (m *morgana) Is(err error) bool {
	morganaError := FromError(err)
	return (morganaError.GetCustomCode() == m.GetCustomCode()) && (m.GetWith() == morganaError.GetWith()) && (m.GetType() == morganaError.GetType())
}

func New(typeValue string) Morgana {
	mor := &morgana{Type: typeValue, morganaStackErrors: make([]Morgana, 0), MetaData: make(map[string]any), InternalDetail: make([]any, 0)}
	mor.WithStackTrace(3)
	return mor
}

func (m *morgana) Clone(stackLevel int) Morgana {

	e := New(m.Type).WithCustomCode(m.CustomCode).WithInternalDetail(m.InternalDetail).
		WithStatusCode(m.StatusCode).WithMessage(m.Msg).With(m.WithValue).WithStackTrace(stackLevel)

	for _, stackError := range m.morganaStackErrors {
		e = e.WithError(stackError.ToError())
	}

	return e
}

func (m *morgana) WithStackTrace(skip int) Morgana {

	pc := make([]uintptr, 15)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	m.StackTrace = fmt.Sprintf("%s:%d %s\n", frame.File, frame.Line, frame.Function)

	return m

}

func (m *morgana) WithInternalDetail(detail ...any) Morgana {

	if detail != nil || len(detail) != 0 {
		m.InternalDetail = append(m.InternalDetail, detail...)
	}
	return m
}

func (m *morgana) WithMessage(msg string, args ...string) Morgana {
	r := strings.NewReplacer(args...)
	m.Msg = r.Replace(msg)
	return m
}

func (m *morgana) GetMessage() string {

	return m.Msg
}

func (m *morgana) WithStatusCode(statusCode int) Morgana {

	m.StatusCode = statusCode
	return m
}

func (m *morgana) WithCustomCode(customCode string) Morgana {
	m.CustomCode = customCode
	return m
}

func (m *morgana) WithType(ref string) Morgana {
	m.Type = ref
	return m
}
func (m *morgana) With(value string) Morgana {

	m.WithValue = value
	return m
}

func (m *morgana) WithError(err error) Morgana {

	if err == nil {
		return m
	}

	empoData := GetEmpo(err)
	if empoData != nil {

		if val, ok := empoData.GetAttributes()[morgana_key_data]; ok {

			if mor, ok := val.(Morgana); ok {

				m.morganaStackErrors = append(m.morganaStackErrors, mor)
			}
		} else {
			m.InternalDetail = append(m.InternalDetail, err.Error())
		}
	} else {
		m.InternalDetail = append(m.InternalDetail, err.Error())
	}

	return m

}

func (m *morgana) GetWithStackTrace() string {

	return m.StackTrace
}

func (m *morgana) GetStatusCode() int {

	return m.StatusCode
}

func (m *morgana) GetCustomCode() string {
	return m.CustomCode
}

func (m *morgana) GetInternalDetail() []any {

	return m.InternalDetail
}

func (m *morgana) GetType() string {

	return m.Type
}
func (m *morgana) GetWith() string {
	return m.WithValue
}

func Wrap(err error, err2 error) error {
	if err == nil && err2 == nil {
		return nil
	}

	err1Morgana := GetMorgana(err)
	err2Morgana := GetMorgana(err2)

	if err1Morgana == nil && err2Morgana == nil {
		if err != nil && err2 != nil {

			return errors.Wrap(err2, err.Error())
		}

		if err != nil && err2 == nil {
			return err
		}

		if err == nil && err2 != nil {
			return err2
		}

	}

	if err1Morgana != nil && err2Morgana == nil {
		err1Morgana.WithError(err2)
		return err1Morgana.ToError()
	}
	if err1Morgana == nil && err2Morgana != nil {
		err2Morgana.WithError(err)
		return err2Morgana.ToError()
	}
	if err1Morgana != nil && err2Morgana != nil {
		err1Morgana.WithError(err2)
		return err1Morgana.ToError()
	}

	return nil
}

func GetStringDetail(err error) string {

	if err == nil {
		return ""
	}

	morg := GetMorgana(err)
	if morg != nil {
		return morg.String()
	}

	return err.Error()
}
