package app

import (
	"errors"
	"fmt"
)

// Application error codes.
//
// NOTE: These are meant to be generic and they map well to HTTP error codes.
// Different applications can have very different error code requirements so
// these should be expanded as needed (or introduce subcodes).
const (
	ERR_CONFLICT       = "conflict"
	ERR_INTERNAL       = "internal"
	ERR_INVALID        = "invalid"
	ERR_NOTFOUND       = "not_found"
	ERR_NOTIMPLEMENTED = "not_implemented"
	ERR_UNAUTHORIZED   = "unauthorized"
)

// Error is an application-specific error. Application errors can be
// unwrapped by the caller to extract out the code & message.
//
// Any non-application error (such as a disk error) should be reported as an
// EINTERNAL error and the human user should only see "Internal error" as the
// message. These low-level internal error details should only be logged and
// reported to the operator of the application (not the end user).
type Error struct {
	// Machine-readable error code.
	Code string

	// Human-readable error message.
	Message string
}

// Error implements the error interface. Not used by the application otherwise.
func (e *Error) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// ErrorCode unwraps an application error and returns its code.
// Non-application errors always return EINTERNAL.
func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	}
	if errors.As(err, &e) {
		return e.Code
	}
	return ERR_INTERNAL
}

// ErrorMessage unwraps an application error and returns its message.
// Non-application errors always return "Internal error".
func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	}
	if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error"
}

// Errorf is a helper function to return an Error with a given code and formatted message
func Errorf(code string, format string, args ...interface{}) error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}
