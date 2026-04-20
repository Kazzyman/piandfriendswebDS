package pkg

import (
	"math/big"
)

// CrossVerify runs BBP decimal computation to check the primary result.
// Returns true if the digits match up to checkDigits, false otherwise.
// This is how real π calculations are validated—not against a hardcoded string.
func CrossVerify(primary *big.Float, checkDigits int) bool {
	if checkDigits <= 0 {
		checkDigits = 100
	}

	// Compute BBP decimal to the same precision
	secondary := ComputeBBPToDigits(checkDigits)

	// Compare digit by digit
	primaryStr := primary.Text('f', checkDigits+2)
	secondaryStr := secondary.Text('f', checkDigits+2)

	minLen := len(primaryStr)
	if len(secondaryStr) < minLen {
		minLen = len(secondaryStr)
	}

	for i := 0; i < minLen; i++ {
		if primaryStr[i] != secondaryStr[i] {
			return false
		}
	}
	return true
}

// ComputeBBPToDigits calculates π using the BBP formula evaluated in base-10.
// This is the same algorithm used in the BBP demo's decimal recomputation.
// It serves as an independent cross-check for other algorithms.
func ComputeBBPToDigits(digits int) *big.Float {
	// Precision in bits: ~3.32 bits per decimal digit, plus margin
	prec := uint(float64(digits)*3.33 + 64)

	pi := new(big.Float).SetPrec(prec).SetFloat64(0.0)

	one := new(big.Float).SetPrec(prec).SetFloat64(1.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)

	// Terms needed: roughly digits * 0.83
	terms := int(float64(digits) * 0.83)
	if terms < 10 {
		terms = 10
	}

	for k := 0; k < terms; k++ {
		// 16^k
		pow16 := new(big.Int).Exp(big.NewInt(16), big.NewInt(int64(k)), nil)
		pow16f := new(big.Float).SetPrec(prec).SetInt(pow16)

		k8 := float64(k * 8)

		t1 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 1)
		t2 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 4)
		t3 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 5)
		t4 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 6)

		term1 := new(big.Float).SetPrec(prec).Quo(four, t1)
		term2 := new(big.Float).SetPrec(prec).Quo(two, t2)
		term3 := new(big.Float).SetPrec(prec).Quo(one, t3)
		term4 := new(big.Float).SetPrec(prec).Quo(one, t4)

		inner := new(big.Float).SetPrec(prec)
		inner.Sub(term1, term2)
		inner.Sub(inner, term3)
		inner.Sub(inner, term4)

		inner.Quo(inner, pow16f)

		pi.Add(pi, inner)
	}

	return pi
}

// VerifyAndReport checks the computed π against BBP decimal and returns a formatted
// result string suitable for display.
func VerifyAndReport(computed *big.Float, digits int, algorithmName string) string {
	if digits <= 0 {
		digits = 100
	}
	// TODO I made this adjustment, because I know reasonable! Pun fortunate! 
	/*
	if digits > 500 {
		digits = 500 // Cap verification to keep response times reasonable
	}
*/
	if CrossVerify(computed, digits) {
		return "COLOR:green:  ✓ Cross-verified against BBP formula — correct to " + itoa(digits) + " digits"
	}
	return "COLOR:red:  ✗ Cross-check failed — possible implementation error"
}

// itoa is a tiny helper to avoid importing strconv in this file.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	tmp := n
	for tmp > 0 {
		digits = append([]byte{byte('0' + tmp%10)}, digits...)
		tmp /= 10
	}
	return string(digits)
}