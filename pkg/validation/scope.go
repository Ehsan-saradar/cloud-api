package validation

import (
	"api.cloud.io/pkg/security/auth/scopes"
)

// IsAdminScope check validation of str scope.
func IsAdminScope(str string) bool {
	return scopes.Admin == str
}
