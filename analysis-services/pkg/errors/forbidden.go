package errors

const (
	wrongPassword = 40300 + iota + 1
	phoneTokenNotExists
	tokenExpired
	invalidToken
	accessDenied
	noAccess
	userBlocked
	businessNotAvailable
	userNotAvailable
)

// Errors
var (
	ErrWrongPassword        = New(wrongPassword, "you sent wrong password")(nil)
	ErrPhoneTokenNotExists  = New(phoneTokenNotExists, "phone token not exists or expired")(nil)
	ErrTokenExpired         = New(tokenExpired, "the access/refresh token is expired")(nil)
	ErrInvalidToken         = New(invalidToken, "the token is invalid")(nil)
	ErrAccessDenied         = New(accessDenied, "could not access to this method")(nil)
	ErrNoAccess             = New(noAccess, "member have no access to provider")(nil)
	ErrUserBlocked          = New(userBlocked, "this user/business/admin blocked and can't get new access token")(nil)
	ErrBusinessNotAvailable = New(businessNotAvailable, " this business is not available")(nil)
	ErrUserNotAvailable     = New(userNotAvailable, " this user is not available")(nil)
)
