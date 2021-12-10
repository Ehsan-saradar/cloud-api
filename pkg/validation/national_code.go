package validation

import (
	"github.com/asaskevich/govalidator"
)

// IsNationalCode check validation of str in IRAN national code format.
func IsNationalCode(str string) bool {
	if !govalidator.IsNumeric(str) || len(str) != 10 {
		return false
	}

	control := int(str[9] - 48)
	sum := 0
	for i := 8; i >= 0; i-- {
		sum += (10 - i) * (int(str[i] - 48))
	}

	gainedControl := sum % 11
	if (gainedControl < 2 && gainedControl == control) || 11-gainedControl == control {
		return true
	}

	return false
}
