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
	b := strings.Builder{}

	runes := []rune(input)

	for i, r := range runes {
		if i == 0 && (unicode.IsDigit(r) && r != '\\') {
			return "", ErrInvalidString
		}
		if escape {
			switch {
			case unicode.IsDigit(r) || r == '\\':
				if pending != "" {
					b.WriteString(pending)
				}
				pending = string(r)
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
			case r == '\\':
				b.WriteString(pending)
				escape = true
				pending = ""

			default:
				b.WriteString(pending)
				pending = string(r)
			}
		}
	}
	if pending != "" {
		b.WriteString(pending)
	}
	output = b.String()
	return output, nil
}
