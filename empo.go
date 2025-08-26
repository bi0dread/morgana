package morgana

type Empo interface {
	WithAttributes(details map[string]any) Empo
	GetAttributes() map[string]any
	ToError() error
	// New: carry a cause and support standard unwrapping
	WithCause(err error) Empo
}

func NewEmpo(msg string) Empo {
	return &empo{msg: msg, details: map[string]any{}}
}

type empo struct {
	msg     string
	details map[string]any
	cause   error
}

func (e *empo) GetAttributes() map[string]any {

	return e.details
}

func (e *empo) ToError() error {
	return e
}

func (e *empo) WithAttributes(details map[string]any) Empo {

	e.details = details
	return e
}

func (e *empo) WithCause(err error) Empo {
	e.cause = err
	return e
}

func (e *empo) Error() string {
	return e.msg
}

// Unwrap enables errors.Is/As to traverse the cause chain
func (e *empo) Unwrap() error {
	return e.cause
}

func GetEmpo(err error) Empo {
	if err == nil {
		return nil
	}

	if e, ok := err.(Empo); ok {
		return e
	}

	return nil
}
