package tools

import (
	"math/rand"
	"time"
)

func Randpw(length int) string {
	if length <= 8 {
		length = 8
	}
	// LowerLetters is the list of lowercase letters.
	lowerLetters := "abcdefghijklmnopqrstuvwxyz"
	// UpperLetters is the list of uppercase letters.
	upperLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Digits is the list of permitted digits.
	digits := "0123456789"
	// Symbols is the list of symbols.
	symbols := "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"

	all := lowerLetters + upperLetters + digits + symbols

	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, length)
	buf[0] = lowerLetters[rand.Intn(len(lowerLetters))]
	buf[1] = upperLetters[rand.Intn(len(upperLetters))]
	buf[2] = digits[rand.Intn(len(digits))]
	buf[3] = symbols[rand.Intn(len(symbols))]
	for i := 4; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf) // E.g. "3i[g0|)z"
}
