package base62

import (
	"math/big"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func ToBase62(data []byte) string {
	num := new(big.Int).SetBytes(data)
	if num.Cmp(big.NewInt(0)) == 0 {
		return string(base62Chars[0])
	}

	base := big.NewInt(62)
	result := ""

	for num.Cmp(big.NewInt(0)) > 0 && len(result) < 9 {
		remainder := new(big.Int)
		num.DivMod(num, base, remainder)
		result = string(base62Chars[remainder.Int64()]) + result
	}

	if len(result) < 8 {
		result = "0" + result
	}
	return result
}
