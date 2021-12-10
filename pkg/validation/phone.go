package validation

import (
	"regexp"
)

var phoneRegexp = regexp.MustCompile(`^\+[1-9]{1}[0-9]{7,11}$`)

// IsPhone check validation of the phone number.
func IsPhone(str string) bool {
	return phoneRegexp.MatchString(str)
}
