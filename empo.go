package morgana

type Empo interface {
	WithAttributes(details map[string]any) Empo
	GetAttributes() map[string]any
	ToError() error
}

func NewEmpo(msg string) Empo {
	return &empo{msg: msg, details: map[string]any{}}
}

type empo struct {
	msg     string
	details map[string]any
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

func (e *empo) Error() string {
	return e.msg
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
