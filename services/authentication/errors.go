package authentication

import "fmt"

type Type int

const (
	ErrTypeAuthorizationMalformed Type = iota
	ErrTypeServerError
)

type Error struct {
	Cause error
	Type  Type
}

func newAuthenticationError(t Type, cause error) error {
	return &Error{
		Type:  t,
		Cause: cause,
	}
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("Type %d: %s", e.Type, e.Cause.Error())
	}
	return fmt.Sprintf("Type %d: no cause given", e.Type)
}
