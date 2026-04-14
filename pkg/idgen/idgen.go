// pkg/idgen/idgen.go
// Generates short, human-readable, prefixed unique IDs for each domain object.
// Format: <PREFIX>_<timestamp_base36><random_base36>
// Examples: ord_lx4k9mz1, prd_lx4k9mz2, pay_lx4k9mz3
package idgen

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

const base36Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

// domain prefixes
const (
	PrefixOrder   = "ord"
	PrefixProduct = "prd"
	PrefixService = "svc"
	PrefixPayment = "pay"
	PrefixCart    = "crt"
	PrefixCoin    = "cxn"
	PrefixCoupon  = "cpn"
	PrefixInvoice = "inv"
)

// New generates a prefixed unique ID.
// e.g. New(PrefixOrder) → "ord_lx4k9mz2a8f"
func New(prefix string) string {
	ts := toBase36(time.Now().UnixMilli())
	rnd := randomBase36(6)
	return fmt.Sprintf("%s_%s%s", prefix, ts, rnd)
}

// toBase36 encodes an int64 in base-36.
func toBase36(n int64) string {
	if n == 0 {
		return "0"
	}
	var sb strings.Builder
	for n > 0 {
		sb.WriteByte(base36Chars[n%36])
		n /= 36
	}
	// reverse
	s := sb.String()
	runes := []byte(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// randomBase36 generates n random base-36 characters.
func randomBase36(n int) string {
	var sb strings.Builder
	max := big.NewInt(36)
	for i := 0; i < n; i++ {
		idx, _ := rand.Int(rand.Reader, max)
		sb.WriteByte(base36Chars[idx.Int64()])
	}
	return sb.String()
}
