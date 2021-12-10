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
	ErrOutOfLimit                     = New(outOfLimit, "you reach the limit of this operation")(nil)
	ErrInvalidParams                  = New(invalidParams, "request params are invalid")
	ErrPhoneAlreadyRegistered         = New(phoneAlreadyRegistered, "this phone already registered")(nil)
	ErrUsernameAlreadyExists          = New(usernameAlreadyExists, "thie username already exists")(nil)
	ErrBusinessRequestClosed          = New(businessRequestClosed, "the requested business request in closed")(nil)
	ErrUserAlreadyExists              = New(userAlreadyExists, "The user with this ID is already exists")(nil)
	ErrUserAlreadyLikeEntity          = New(userAlreadyLikeEntity, "The user already like this item")(nil)
	ErrCommentHasReply                = New(commentHasReply, "The comment currently have a reply")(nil)
	ErrBusinessAlreadyExists          = New(businessAlreadyExists, "The business with this ID is already exists")(nil)
	ErrRelationAlreadyExists          = New(relationAlreadyExists, "these users already have relation")(nil)
	ErrUserAlreadyHaveBusiness        = New(userAlreadyHaveBusiness, "user already have an business registered")(nil)
	ErrCategoryNotAvailable           = New(categoryNotAvailable, "the following category not exists or is deprecated")(nil)
	ErrSubCategoryNotAvailable        = New(subCategoryNotAvailable, "the following sub category not exists or is deprecated")(nil)
	ErrInvalidPropertyDefinition      = New(invalidPropertyDefinition, "the schema of property definition is invalid")(nil)
	ErrInvalidPropertyValue           = New(invalidPropertyValue, "the value of property is invalid")(nil)
	ErrInvalidPropertyRequest         = New(invalidPropertyRequest, "requested property is not available here")(nil)
	ErrEntityAlreadyIsFavorite        = New(entityAlreadyIsFavorite, "requested entity is already in favorites")(nil)
	ErrRatingAlreadyExists            = New(ratingAlreadyExists, "this rating already exists")(nil)
	ErrCityNotAvailable               = New(cityNotAvailable, "the following city not exists or not active")(nil)
	ErrMismatchCategoryAndSubCategory = New(mismatchCategoryAndSubCategory, "category and sub category have a mismatch")(nil)
	ErrRecaptchaFailed                = New(recaptchaFailed, "recaptcha challenge failed")(nil)
	ErrAppCategoryNotAvailable        = New(appCategoryNotAvailable, "the following app category not exists or is deprecated")(nil)
)
