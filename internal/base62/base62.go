package base62

import "strings"

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ToBase62(numToHash int64) string {
	base62Builder := strings.Builder{}
	base62Builder.Grow(6)

	for numToHash > 0 {
		reminder := numToHash % 62

		base62Builder.WriteByte(base62Chars[reminder])

		numToHash /= 62
	}

	return base62Builder.String()
}
