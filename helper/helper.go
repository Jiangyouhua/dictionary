package helper

import (
	"strings"
	"unicode"
)

func Camel2Case(name string) string {
	ss := make([]rune, 0)
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				ss = append(ss, '_')
			}
			ss = append(ss, unicode.ToLower(r))
			continue
		}
		ss = append(ss, r)

	}
	return string(ss)
}

func Translate(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\uFEFF", "", -1)
	s = strings.ReplaceAll(s, "'", `\'`)
	s = strings.ReplaceAll(s, `\\'`, `\'`)
	return s
}
