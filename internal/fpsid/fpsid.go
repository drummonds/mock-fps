package fpsid

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// NumericReference returns a 6-digit zero-padded numeric string.
func NumericReference() string {
	return randDigits(6)
}

// EndToEndReference returns "FPS" + YYYYMMDD + 6-digit numeric (17 chars).
func EndToEndReference() string {
	date := time.Now().UTC().Format("20060102")
	return "FPS" + date + randDigits(6)
}

// SchemeTransactionID returns a 26-digit numeric string.
func SchemeTransactionID() string {
	return randDigits(26)
}

// ProcessingDate returns today's date as YYYY-MM-DD.
func ProcessingDate() string {
	return time.Now().UTC().Format(time.DateOnly)
}

func randDigits(n int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), nil)
	v, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand failed: %v", err))
	}
	return fmt.Sprintf("%0*s", n, v.String())
}
