package utils

import "math/rand"

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
var numbers = []rune("1234567890")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandInt(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = numbers[rand.Intn(len(numbers))]
	}
	if b[0] == 0 || b[0] == '0' {
		b[0] = '1'
	}
	return string(b)
}
