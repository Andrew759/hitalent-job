package service

import (
	"strconv"
	"unicode/utf8"
)

func IsHasCorrectLength(s string, minLength int, maxLength int) bool {
	sLen := utf8.RuneCountInString(s)

	return sLen > minLength && sLen <= maxLength
}

func isInteger(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}
