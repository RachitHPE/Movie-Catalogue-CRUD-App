package error

import (
	"fmt"
	"reflect"
)

const stackSize = 4096

// A Tag represents an error identifier of any type.
type Tag interface{}

// A Gerror is a tagged error with a stack trace embedded in the Error() string.
type Gerror interface {
	// Returns the tag used to create this error.
	Tag() Tag

	// Returns the concrete type of the tag used to create this error.
	TagType() reflect.Type

	// Returns the string form of this error,
	// which includes the tag value, the tag type, the error message, and a stack trace.
	Error() string

	// Test the tag used to create this error for equality with a given tag.
	// Returns `true` if and only if the two are equal.
	EqualTag(Tag) bool

	// Message
	Message() string

	// Cause
	Cause() error
}

// New Returns an error containing the given tag and message and the current stack trace.
func New(tag Tag, message string) *GeneralError {
	return &GeneralError{
		tag:     tag,
		typ:     reflect.TypeOf(tag),
		cause:   nil,
		message: message,
	}
}

// Newf Returns an error containing the given tag and format string and the current stack trace.
// The given inserts are applied to the format string to produce an error message.
func Newf(tag Tag, format string, insert ...interface{}) Gerror {
	return New(tag, fmt.Sprintf(format, insert...))
}

// NewFromError Return an error containing the given tag, the cause of the error, and the current stack trace.
func NewFromError(tag Tag, cause error) Gerror {
	if cause != nil {
		return &GeneralError{
			tag:     tag,
			typ:     reflect.TypeOf(tag),
			cause:   cause,
			message: cause.Error(),
		}
	}

	return nil
}

type GeneralError struct {
	tag     Tag
	typ     reflect.Type
	cause   error
	message string
}

func (e *GeneralError) Error() string {
	return e.message
}

func (e *GeneralError) Tag() Tag {
	return e.tag
}

func (e *GeneralError) TagType() reflect.Type {
	return e.typ
}

func (e *GeneralError) EqualTag(tag Tag) bool {
	return e.typ == reflect.TypeOf(tag) && e.tag == tag
}

func (e *GeneralError) Message() string {
	return e.message
}

func (e *GeneralError) Cause() error {
	return e.cause
}

// ErrorCode service error constants.
type ErrorCode string

func (e ErrorCode) String() string {
	return string(e)
}

// GetErrorType ...
func GetErrorType(err error) ErrorCode {
	gErr, ok := err.(Gerror)
	if ok {
		if errCode, ok2 := gErr.Tag().(ErrorCode); ok2 {
			return errCode
		}
	}

	return InternalError
}