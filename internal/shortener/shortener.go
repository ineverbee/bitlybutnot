package shortener

import (
	"errors"
	"hash/fnv"
	"math"
	"strings"
)

const (
	linkLength  = 10
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	bytesLength = 63
)

//UniqueString func used to generate new unique short link
func UniqueString(n uint32) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)

	for i := 0; n > 0 || i < linkLength; n, i = n/63, i+1 {
		encodedBuilder.WriteByte(letterBytes[(n % 63)])
	}

	return encodedBuilder.String()
}

//HashURL func used to convert long URL into specific number
func HashURL(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

//DecodeString func used to decode shortened link into number
func DecodeString(encoded string) (uint32, error) {
	var number uint32

	for i, symbol := range encoded {
		alphabeticPosition := strings.IndexRune(letterBytes, symbol)

		if alphabeticPosition == -1 {
			return uint32(alphabeticPosition), errors.New("invalid character: " + string(symbol))
		}
		if i == len(encoded)-1 {
			number += uint32(alphabeticPosition) * (uint32(math.Pow(float64(63), float64(i))) - 1)
			continue
		}
		number += uint32(alphabeticPosition) * uint32(math.Pow(float64(63), float64(i)))
	}

	return number, nil
}
