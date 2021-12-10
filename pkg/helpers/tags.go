package helpers

import (
	"strings"
	"unicode"
)

func GetTags(s string) []string {
	fields := strings.FieldsFunc(s, func(r rune) bool {
		return !(unicode.IsLetter(r) || unicode.IsDigit(r) || r == '#' || r == '_')
	})

	res := []string{}
	for _, v := range fields {
		if strings.Count(v, "#") != 1 {
			continue
		}

		runes := []rune(v)
		idx := strings.IndexRune(v, '#')
		if idx < len(runes)-1 && !unicode.IsNumber(runes[idx+1]) {
			tag := string(runes[idx+1:])
			res = append(res, tag)
		}
	}
	return res
}
