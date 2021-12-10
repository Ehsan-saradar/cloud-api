package validation

import (
	"github.com/asaskevich/govalidator"
)

const (
	pictureMaxSize = (1 << 20) * 1
)

func init() {
	govalidator.TagMap["phone"] = IsPhone
	govalidator.TagMap["adminScope"] = IsAdminScope
	govalidator.TagMap["username"] = IsUsername
	govalidator.TagMap["nationalCode"] = IsNationalCode
	govalidator.TagMap["reportReferenceType"] = IsReportReferenceType
}

func isImageContentType(typ string) bool {
	return typ == "image/png" || typ == "image/jpeg" || typ == "image/webp"
}
