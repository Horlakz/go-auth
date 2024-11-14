package helper

import (
	"math/rand"
	"time"
)

func GenerateRandomDigits(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = byte(rand.Intn(10) + 48)
	}

	return string(code)
}
