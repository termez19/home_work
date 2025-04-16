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

	for _, r := range input {
		if !escape {
			switch {
			case unicode.IsDigit(r):
				if pending == "" {
					return "", ErrInvalidString
				}
				repeatTimes, _ := strconv.Atoi(string(r))
				b.WriteString(strings.Repeat(pending, repeatTimes))
				pending = ""
			case r == '\\':
				b.WriteString(pending)
				escape = true
				pending = ""

			default:
				b.WriteString(pending)
				pending = string(r)
			}
		} else {
			if !unicode.IsDigit(r) && r != '\\' {
				return "", ErrInvalidString
			}
			pending = string(r)
			escape = false
		}
	}
	if pending != "" {
		b.WriteString(pending)
	}
	output = b.String()
	return output, nil
}
