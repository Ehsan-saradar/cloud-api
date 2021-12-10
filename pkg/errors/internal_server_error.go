package errors

const (
	internal = 50000 + iota + 1
)

// Errors
var (
	ErrInternal = New(internal, "server can not do this right")
)
