package utils

import (
	"fmt"
	"math/rand"
)

func RandVerificationCode() string {
	num := 100_000 + rand.Intn(899_999)
	return fmt.Sprintf("%d", num)
}
