package errors

const (
	payloadTooLarge = 41300 + iota + 1
)

// Errors
var (
	ErrPayloadTooLarge = New(payloadTooLarge, "the payload is very large to process")(nil)
)
