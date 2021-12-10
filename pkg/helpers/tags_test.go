package helpers

import (
	"reflect"
	"testing"
)

func TestGetTags(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want []string
	}{
		{
			"Extract",
			"This #is #a# test. We testing # #123 and #hashtags so #w8 for result #",
			[]string{"is", "hashtags", "w8"},
		},
		{
			"Unicode",
			"#چرا؟",
			[]string{"چرا"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTags(tt.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
