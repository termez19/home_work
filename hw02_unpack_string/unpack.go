package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var output string
	var pending string
	escape := false
	previousRune := rune(0)
	b := strings.Builder{}

	runes := []rune(input)

	for i, r := range runes {
		if !isLegitRune(r) {
			return "", ErrInvalidString
		}
		if i == 0 && (!unicode.IsLetter(r) && r != '\\') {
			return "", ErrInvalidString
		}
		if escape {
			switch {
			case unicode.IsDigit(r) || r == '\\':
				pending = string(r)
				previousRune = r
				escape = false
			default:
				return "", ErrInvalidString

			}
		} else {
			switch {
			case unicode.IsDigit(r):
				if pending == "" {
					return "", ErrInvalidString
				}
				repeat, _ := strconv.Atoi(string(r))
				b.WriteString(strings.Repeat(pending, repeat))
				pending = ""

			case unicode.IsLetter(r):
				if pending == "" {
					pending = string(r)
					previousRune = (r)
				} else {
					b.WriteString(pending)
					previousRune = r
					pending = string(r)
				}
			case r == '\\':
				if previousRune == previousRune {
				}
				escape = true
				previousRune = r
				pending = ""
			}
		}
	}
	if pending != "" {
		b.WriteString(pending)
	}
	output = b.String()
	return output, nil
}

func isLegitRune(r rune) bool {
	if r == '\\' || unicode.IsLetter(r) || unicode.IsDigit(r) {
		return true
	}
	return false
}
