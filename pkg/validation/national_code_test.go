package validation

import (
	"testing"
)

func TestIsNationalCode(t *testing.T) {
	tests := []struct {
		name string
		str  string
		want bool
	}{
		{
			"InvalidNationalCode",
			"7731689951",
			false,
		},
		{
			"ValidNationalCode",
			"0925211788",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNationalCode(tt.str); got != tt.want {
				t.Errorf("IsNationalCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
