package errors

const (
	otpNotExists = 40400 + iota + 1
	adminNotExists
	userNotExists
	relationNotExists
	objectNotExists
	businessNotExists
	tokenNotExists
	commentNotExists
	likeRecordNotExists
	favoriteRecordNotExists
	categoryNotExists
	subCategoryNotExists
	propertyNotExists
	productNotExists
	pictureIndexOutOfRange
	eventNotExists
	advertisementNotExists
	ratingNotExists
	reportNotExists
	talkNotExists
	cityNotExists
	postNotExists
	issueNotExists
	usernameNotExists
	appCategoryNotExists
)

// Errors
var (
	ErrOTPNotExists            = New(otpNotExists, "phone/email verification code expired or not exists")(nil)
	ErrAdminNotExists          = New(adminNotExists, "this admin not exists")(nil)
	ErrUserNotExists           = New(userNotExists, "this user not exists")(nil)
	ErrRelationNotExists       = New(relationNotExists, "there is no relation between these users")(nil)
	ErrObjectNotExists         = New(objectNotExists, "ths requested object of user not exists")(nil)
	ErrBusinessNotExists       = New(businessNotExists, "this business not exists")(nil)
	ErrTokenNotExists          = New(tokenNotExists, "requested token not exists")(nil)
	ErrCommentNotExists        = New(commentNotExists, "this comment not exists")(nil)
	ErrLikeRecordNotExists     = New(likeRecordNotExists, "the entity wasn't like by user")(nil)
	ErrFavoriteRecordNotExists = New(favoriteRecordNotExists, "the entity wasn't in user favorites")(nil)
	ErrCategoryNotExists       = New(categoryNotExists, "this category not exists")(nil)
	ErrSubCategoryNotExists    = New(subCategoryNotExists, "this sub category not exists")(nil)
	ErrPropertyNotExists       = New(propertyNotExists, "this property not exists")(nil)
	ErrProductNotExists        = New(productNotExists, "this product not exists")(nil)
	ErrPictureIndexOutOfRange  = New(pictureIndexOutOfRange, "requested picture index is out of range")(nil)
	ErrEventNotExists          = New(eventNotExists, "this event not exists")(nil)
	ErrAdvertisementNotExists  = New(advertisementNotExists, "this advertisement not exists")(nil)
	ErrRatingNotExists         = New(ratingNotExists, "rating not exists")(nil)
	ErrReportNotExists         = New(reportNotExists, "this report not exists")(nil)
	ErrTalkNotExists           = New(talkNotExists, "talk not exists")(nil)
	ErrCityNotExists           = New(cityNotExists, "this city not exists")(nil)
	ErrPostNotExists           = New(postNotExists, "this post not exists")(nil)
	ErrIssueNotExists          = New(issueNotExists, "this issue not exists")(nil)
	ErrUsernameNotExists       = New(usernameNotExists, "this username not exists")(nil)
	ErrAppCategoryNotExists    = New(appCategoryNotExists, "this app category not exists")(nil)
)
