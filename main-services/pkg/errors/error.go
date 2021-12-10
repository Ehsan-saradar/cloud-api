package errors

import (
	"github.com/semrush/zenrpc"
)

// ErrorClass is the base func of error.
type ErrorClass func(error, ...interface{}) *zenrpc.Error

// New create a new ErrorClass.
func New(code int, message string) ErrorClass {
	return func(err error, data ...interface{}) *zenrpc.Error {
		zenrpcError := &zenrpc.Error{
			Code:    code,
			Message: message,
			Err:     err,
		}

		if len(data) > 0 {
			zenrpcError.Data = data[0]
		}

		return zenrpcError
	}
}
