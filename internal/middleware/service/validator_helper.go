package service

import (
	"unicode/utf8"
)

func IsHasCorrectLength(s string, minLength int, maxLength int) bool {
	sLen := utf8.RuneCountInString(s)

	return sLen >= minLength && sLen <= maxLength
}
