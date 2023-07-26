package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

//goland:noinspection ALL
func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func GenerateRandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return GenerateRandomString(7)
}

func RandomMoney() int64 {
	return RandomInt(0, 100000)
}

func RandomCurrency() string {
	currencies := [4]string{"USD", "EUR", "KZT", "RUB"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
