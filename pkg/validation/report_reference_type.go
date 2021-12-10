package validation

// IsReportReferenceType check validation of str report reference type.
func IsReportReferenceType(str string) bool {
	switch str {
	case "Business", "Comment", "Event", "Product", "UserBusinessTalk", "User":
		return true
	}

	return false
}
