package errors

const (
	unauthorized = 40100 + iota + 1
)

// Errors
var (
	ErrUnauthorized = New(unauthorized, "JWT is not present in request")(nil)
)
