package morgana

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"errors"

	perror "github.com/pkg/errors"
)

const (
	morgana_key_data = "morgana_key_data"
)

type Morgana interface {
	WithStatusCode(statusCode int) Morgana
	WithCustomCode(customCode string) Morgana
	WithType(ref string) Morgana
	WithMessage(msg string, args ...string) Morgana
	WithStackTrace(skip int) Morgana
	With(value string) Morgana
	WithError(err error) Morgana
	WithAddMetaDataKey(key string, value any) Morgana
	WithCause(err error) Morgana
	Cause() error

	GetStatusCode() int
	GetCustomCode() string
	GetType() string
	GetMessage() string
	GetWith() string
	GetMetaDataKey(key string) any
	GetMetaData() map[string]any
	GetMorganaStackErrors() []Morgana

	Clone(stackLevel int) Morgana

	ToError() error

	String() string
	ToJson() string

	// New functionality
	WithFullStack(skip int, maxFrames int) Morgana
	GetStackFrames() []StackFrame
	WithAddMetaData(map[string]any) Morgana
	HasMetaDataKey(key string) bool
	WithRedactedKey(key string) Morgana
	ToJsonSafe() string
	StringSafe() string
	WriteHTTP(w http.ResponseWriter, safe bool)
	ToFields() map[string]any
	WithTrace(ctx context.Context) Morgana
	WithID(id string) Morgana
	GetID() string
	WithFieldError(field string, code string, msg string) Morgana
	GetFieldErrors() []FieldError

	// gRPC helpers
	ToGRPCCode() int
	FromGRPCCode(code int) Morgana
}

type StackFrame struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code,omitempty"`
	Msg   string `json:"msg"`
}

type morgana struct {
	Type               string
	WithValue          string
	Msg                string
	StatusCode         int
	CustomCode         string
	StackTrace         string
	morganaStackErrors []Morgana
	MetaData           map[string]any
	// New fields
	StackFrames  []StackFrame
	redactedKeys map[string]struct{}
	ID           string
	FieldErrors  []FieldError
	cause        error
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

func (m *morgana) WithAddMetaData(md map[string]any) Morgana {
	if md == nil {
		return m
	}
	if m.MetaData == nil {
		m.MetaData = make(map[string]any)
	}
	for k, v := range md {
		m.MetaData[k] = v
	}
	return m
}

func (m *morgana) GetMetaDataKey(key string) any {
	if m.MetaData == nil {
		return ""
	}

	if val, ok := m.MetaData[key]; ok {
		return val
	}

	return ""
}

func (m *morgana) HasMetaDataKey(key string) bool {
	if m.MetaData == nil {
		return false
	}
	_, ok := m.MetaData[key]
	return ok
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

	if len(m.StackFrames) != 0 {
		builder.WriteString("StackFrames: ")
		for i, f := range m.StackFrames {
			builder.WriteString(fmt.Sprintf("[%d] %s:%d %s ", i, f.File, f.Line, f.Function))
		}
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	if len(m.morganaStackErrors) != 0 {
		builder.WriteString("MorganaStackErrors: \n")
		builder.WriteString(fmt.Sprintf("%v", m.morganaStackErrors))

	}

	if len(m.MetaData) != 0 {
		builder.WriteString(fmt.Sprintf("MetaData: %+#v\n", m.MetaData))
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	if len(m.FieldErrors) != 0 {
		builder.WriteString(fmt.Sprintf("FieldErrors: %+#v\n", m.FieldErrors))
		builder.WriteString(fmt.Sprintf("%v\n", " , -----------------------------------------------------------"))
	}

	if len(m.ID) != 0 {
		builder.WriteString(fmt.Sprintf("ID: %s", m.ID))
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

	if len(m.MetaData) != 0 {
		builder.WriteString(fmt.Sprintf("MetaData: %v\n", m.MetaData))
		builder.WriteString(fmt.Sprintf("%v", " , "))
	}

	builder.WriteString(fmt.Sprintf("StatusCode: %d", m.StatusCode))
	builder.WriteString(fmt.Sprintf("%v", " , "))

	return builder.String()
}

func (m *morgana) ToError() error {

	emp := NewEmpo(m.stringSimple()).WithAttributes(map[string]any{morgana_key_data: m})
	if mCause := m.Cause(); mCause != nil {
		emp = emp.WithCause(mCause)
	}
	return emp.ToError()
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

	return New("GENERAL").WithStatusCode(http.StatusNotImplemented).WithMessage(err.Error()).WithCause(err)
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
	mor := &morgana{Type: typeValue, morganaStackErrors: make([]Morgana, 0), MetaData: make(map[string]any), StackFrames: make([]StackFrame, 0), redactedKeys: make(map[string]struct{}), FieldErrors: make([]FieldError, 0)}
	//mor.WithStackTrace(3)
	//mor.WithFullStack(3, 32)
	mor.ensureID()
	return mor
}

func (m *morgana) Clone(stackLevel int) Morgana {

	e := New(m.Type).WithCustomCode(m.CustomCode).
		WithStatusCode(m.StatusCode).WithMessage(m.Msg).With(m.WithValue).WithStackTrace(stackLevel)

	// deep-copy metadata
	if len(m.MetaData) != 0 {
		md := make(map[string]any, len(m.MetaData))
		for k, v := range m.MetaData {
			md[k] = v
		}
		e = e.WithAddMetaData(md)
	}

	// copy redacted keys
	if mm, ok := e.(*morgana); ok {
		mm.StackFrames = append(mm.StackFrames[:0:0], m.StackFrames...)
		mm.FieldErrors = append(mm.FieldErrors[:0:0], m.FieldErrors...)
		for k := range m.redactedKeys {
			if mm.redactedKeys == nil {
				mm.redactedKeys = make(map[string]struct{})
			}
			mm.redactedKeys[k] = struct{}{}
		}
		mm.ID = m.ID
		mm.cause = m.cause
	}

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

func (m *morgana) WithFullStack(skip int, maxFrames int) Morgana {
	if maxFrames <= 0 {
		maxFrames = 32
	}
	pc := make([]uintptr, maxFrames)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	m.StackFrames = m.StackFrames[:0]
	for {
		frame, more := frames.Next()
		m.StackFrames = append(m.StackFrames, StackFrame{File: frame.File, Line: frame.Line, Function: frame.Function})
		if !more {
			break
		}
	}
	return m
}

func (m *morgana) GetStackFrames() []StackFrame {
	return m.StackFrames
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

	// If the error carries a Morgana, attach it to the stack
	if mor := GetMorgana(err); mor != nil {
		m.morganaStackErrors = append(m.morganaStackErrors, mor)
		return m
	}

	// Otherwise, traverse unwrap chain and convert each into a Morgana stack error
	for e := err; e != nil; e = errors.Unwrap(e) {
		if mor := GetMorgana(e); mor != nil {
			m.morganaStackErrors = append(m.morganaStackErrors, mor)
			continue
		}
		// Create a lightweight Morgana for this error without touching parent InternalDetail
		child := New("GENERAL").WithMessage(e.Error())
		m.morganaStackErrors = append(m.morganaStackErrors, child)
		// Optionally record cause for chain traversal
		m.WithCause(e)
	}

	return m

}

func (m *morgana) WithCause(err error) Morgana {
	if err == nil {
		return m
	}
	if m.cause != nil {
		return m
	}
	root := err
	for {
		u := errors.Unwrap(root)
		if u == nil {
			break
		}
		root = u
	}
	m.cause = root
	return m
}

func (m *morgana) Cause() error {
	if m.cause != nil {
		return m.cause
	}
	if len(m.morganaStackErrors) != 0 {
		last := m.morganaStackErrors[len(m.morganaStackErrors)-1]
		if c := last.Cause(); c != nil {
			return c
		}
		return last.ToError()
	}
	return nil
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

			return perror.Wrap(err2, err.Error())
		}

		if err != nil && err2 == nil {
			return err
		}

		if err == nil && err2 != nil {
			return err2
		}

	}

	if err1Morgana != nil && err2Morgana == nil {
		err1Morgana.WithError(err2).WithCause(err)
		return err1Morgana.ToError()
	}
	if err1Morgana == nil && err2Morgana != nil {
		err2Morgana.WithError(err).WithCause(err2)
		return err2Morgana.ToError()
	}
	if err1Morgana != nil && err2Morgana != nil {
		err1Morgana.WithError(err2).WithCause(err)
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

// -------- New helper functionality --------

func (m *morgana) WithRedactedKey(key string) Morgana {
	if m.redactedKeys == nil {
		m.redactedKeys = make(map[string]struct{})
	}
	m.redactedKeys[key] = struct{}{}
	return m
}

func (m *morgana) redactMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}
	out := make(map[string]any, len(input))
	for k, v := range input {
		if _, ok := m.redactedKeys[k]; ok {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}

func (m *morgana) ToJsonSafe() string {
	type safe struct {
		Type        string         `json:"type,omitempty"`
		With        string         `json:"with,omitempty"`
		Msg         string         `json:"msg,omitempty"`
		StatusCode  int            `json:"statusCode,omitempty"`
		CustomCode  string         `json:"customCode,omitempty"`
		StackTrace  string         `json:"stackTrace,omitempty"`
		StackFrames []StackFrame   `json:"stackFrames,omitempty"`
		MetaData    map[string]any `json:"metaData,omitempty"`
		FieldErrors []FieldError   `json:"fieldErrors,omitempty"`
		ID          string         `json:"id,omitempty"`
	}
	s := safe{
		Type:        m.Type,
		With:        m.WithValue,
		Msg:         m.Msg,
		StatusCode:  m.StatusCode,
		CustomCode:  m.CustomCode,
		StackTrace:  m.StackTrace,
		StackFrames: m.StackFrames,
		MetaData:    m.redactMap(m.MetaData),
		FieldErrors: m.FieldErrors,
		ID:          m.ID,
	}
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)
}

func (m *morgana) StringSafe() string {
	// Use safe JSON for a concise representation
	return m.ToJsonSafe()
}

func (m *morgana) WriteHTTP(w http.ResponseWriter, safe bool) {
	if w == nil {
		return
	}
	if m.StatusCode == 0 {
		m.StatusCode = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(m.StatusCode)
	if safe {
		_, _ = w.Write([]byte(m.ToJsonSafe()))
		return
	}
	_, _ = w.Write([]byte(m.ToJson()))
}

func (m *morgana) ToFields() map[string]any {
	fields := map[string]any{
		"type":       m.Type,
		"with":       m.WithValue,
		"msg":        m.Msg,
		"statusCode": m.StatusCode,
		"customCode": m.CustomCode,
		"id":         m.ID,
	}
	if len(m.StackTrace) != 0 {
		fields["stackTrace"] = m.StackTrace
	}
	if len(m.StackFrames) != 0 {
		fields["stackFrames"] = m.StackFrames
	}
	if len(m.FieldErrors) != 0 {
		fields["fieldErrors"] = m.FieldErrors
	}
	if len(m.MetaData) != 0 {
		fields["metaData"] = m.redactMap(m.MetaData)
	}
	return fields
}

func (m *morgana) WithTrace(ctx context.Context) Morgana {
	if ctx == nil {
		return m
	}
	try := func(key any) {
		if v := ctx.Value(key); v != nil {
			m.WithAddMetaDataKey(fmt.Sprintf("%v", key), v)
		}
	}
	// common keys
	try("trace_id")
	try("traceId")
	try("request_id")
	try("requestId")
	try("correlation_id")
	try("correlationId")
	return m
}

func (m *morgana) WithID(id string) Morgana {
	m.ID = id
	return m
}

func (m *morgana) GetID() string {
	return m.ID
}

func (m *morgana) ensureID() {
	if m.ID != "" {
		return
	}
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return
	}
	m.ID = hex.EncodeToString(b)
}

func (m *morgana) WithFieldError(field string, code string, msg string) Morgana {
	m.FieldErrors = append(m.FieldErrors, FieldError{Field: field, Code: code, Msg: msg})
	return m
}

func (m *morgana) GetFieldErrors() []FieldError {
	return m.FieldErrors
}

// gRPC helpers (code mapping only, no external dependency)
func (m *morgana) ToGRPCCode() int {
	// numeric values follow google.golang.org/grpc/codes mapping
	// 0 OK, 1 Canceled, 2 Unknown, 3 InvalidArgument, ...
	switch m.StatusCode {
	case http.StatusOK:
		return 0 // OK
	case http.StatusBadRequest:
		return 3 // InvalidArgument
	case http.StatusUnauthorized:
		return 16 // Unauthenticated
	case http.StatusForbidden:
		return 7 // PermissionDenied
	case http.StatusNotFound:
		return 5 // NotFound
	case http.StatusConflict:
		return 6 // AlreadyExists
	case http.StatusTooManyRequests:
		return 8 // ResourceExhausted
	case http.StatusNotImplemented:
		return 12 // Unimplemented
	case http.StatusServiceUnavailable:
		return 14 // Unavailable
	case http.StatusGatewayTimeout:
		return 4 // DeadlineExceeded
	default:
		if m.StatusCode >= 500 {
			return 13 // Internal
		}
		return 2 // Unknown
	}
}

func (m *morgana) FromGRPCCode(code int) Morgana {
	var httpCode int
	switch code {
	case 0: // OK
		httpCode = http.StatusOK
	case 3: // InvalidArgument
		httpCode = http.StatusBadRequest
	case 16: // Unauthenticated
		httpCode = http.StatusUnauthorized
	case 7: // PermissionDenied
		httpCode = http.StatusForbidden
	case 5: // NotFound
		httpCode = http.StatusNotFound
	case 6: // AlreadyExists
		httpCode = http.StatusConflict
	case 8: // ResourceExhausted
		httpCode = http.StatusTooManyRequests
	case 12: // Unimplemented
		httpCode = http.StatusNotImplemented
	case 14: // Unavailable
		httpCode = http.StatusServiceUnavailable
	case 4: // DeadlineExceeded
		httpCode = http.StatusGatewayTimeout
	default:
		httpCode = http.StatusInternalServerError
	}
	m.StatusCode = httpCode
	return m
}

// Implement fmt.Formatter for pretty printing with %+v
func (m *morgana) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprint(s, m.String())
			return
		}
		fallthrough
	case 's':
		_, _ = fmt.Fprint(s, m.stringSimple())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", m.stringSimple())
	}
}

// Panic helpers
func FromPanic(p any) Morgana {
	if p == nil {
		return nil
	}
	m := New("PANIC").WithMessage(fmt.Sprintf("panic: %v", p)).WithStatusCode(http.StatusInternalServerError)
	m.WithFullStack(4, 64)
	return m
}

func Recover(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = FromPanic(r).ToError()
		}
	}()
	if fn != nil {
		fn()
	}
	return nil
}

// Join joins multiple errors, preferring Morgana if present
func Join(errs ...error) error {
	var firstM Morgana
	for _, e := range errs {
		if e == nil {
			continue
		}
		if m := GetMorgana(e); m != nil {
			if firstM == nil {
				firstM = m.Clone(4)
			}
			firstM.WithError(e)
		}
	}
	if firstM != nil {
		return firstM.ToError()
	}
	// Fallback to local join
	return joinErrors(errs...)
}

func joinErrors(errs ...error) error {
	var nonNil []error
	for _, e := range errs {
		if e != nil {
			nonNil = append(nonNil, e)
		}
	}
	if len(nonNil) == 0 {
		return nil
	}
	if len(nonNil) == 1 {
		return nonNil[0]
	}
	root := nonNil[0]
	for i := 1; i < len(nonNil); i++ {
		root = perror.Wrap(nonNil[i], root.Error())
	}
	return root
}
