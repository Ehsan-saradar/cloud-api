package errors

const (
	invalidParams = 40000 + iota + 1
)

// Errors
var (
	ErrInvalidParams = New(invalidParams, "request params are invalid")
)
