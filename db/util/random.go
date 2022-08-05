package util

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min * rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomUsername() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "HKT", "GBP"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomLocation() string {
	return RandomString(6)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomContactNumber() sql.NullString {
	return sql.NullString{String: RandomString(0)}
}

func RandomAddress() sql.NullString {
	return sql.NullString{String: RandomString(0)}
}

func RandomExchangeRate() float64 {
	return rand.Float64() * float64(RandomInt(1, 13))
}

func RandomFloat64() float64 {
	return rand.Float64() * float64(RandomInt(1, 1000))
}

func RoundFloatAtDecimal(num float64, decimal int) (float64, error) {
	if num <= 0 && num > 10 {
		return 0.0, fmt.Errorf("rounding support 1 to 10 decimal only")
	}
	decimalLevel := math.Pow(10, float64(decimal))
	return math.Round(num*decimalLevel) / decimalLevel, nil
}
