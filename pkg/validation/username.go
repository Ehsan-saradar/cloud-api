package validation

import (
	"regexp"
)

var usernameRegexp = regexp.MustCompile(`^[a-z][a-z0-9_-]{1,13}[a-z0-9]$`)

// IsUsername check validation of username str.
func IsUsername(str string) bool {
	return usernameRegexp.MatchString(str)
}
