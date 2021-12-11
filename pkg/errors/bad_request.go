package errors

const (
	outOfLimit = 40000 + iota + 1
	invalidParams
	phoneAlreadyRegistered
	usernameAlreadyExists
	businessRequestClosed
	userAlreadyExists
	userAlreadyLikeEntity
	commentHasReply
	businessAlreadyExists
	relationAlreadyExists
	userAlreadyHaveBusiness
	categoryNotAvailable
	subCategoryNotAvailable
	invalidPropertyDefinition
	invalidPropertyValue
	invalidPropertyRequest
	entityAlreadyIsFavorite
	ratingAlreadyExists
	cityNotAvailable
	mismatchCategoryAndSubCategory
	recaptchaFailed
	appCategoryNotAvailable
)

// Errors
var (
	ErrInvalidParams = New(invalidParams, "request params are invalid")
)
