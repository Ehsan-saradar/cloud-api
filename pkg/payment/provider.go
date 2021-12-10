package payment

// Request is the data required to request a new payment session from the provider.
// Optional and platform specific parameters are stored in Meta.
type Request struct {
	Description string
	Amount      int64
	CallbackURL string
	Meta        map[string]string
}

// Session represents the payment session and should be stored in data store or memory cache.
// Unnecessary parameters are stored in Meta.
type Session struct {
	Key        string
	GatewayURL string
	Amount     int64
	Meta       map[string]string
}

// Provider is responsible for requesting payment from the provider and verifying the callback.s
type Provider interface {
	// New requests the provider to create a new payment session.
	New(*Request) (*Session, error)
	// Verify the session after the user returned to callback url and retrive the reference id for payment.
	Verify(*Session) (string, error)
}
