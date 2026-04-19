package pkg

import (
	"math/big"
)

// CrossVerify runs a fast secondary algorithm (BBP) to check the primary result.
// Returns true if the digits match up to checkDigits, false otherwise.
// This is how real π calculations are validated—not against a hardcoded string.
func CrossVerify(primary *big.Float, checkDigits int) bool {
	if checkDigits <= 0 {
		checkDigits = 100
	}

	// Compute BBP to the same precision
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

// ComputeBBPToDigits calculates π using the BBP formula to the specified
// number of decimal digits. Used as a cross-check for other algorithms.
func ComputeBBPToDigits(digits int) *big.Float {
	// Precision in bits: ~3.32 bits per decimal digit, plus margin
	prec := uint(float64(digits)*3.33 + 64)

	pi := new(big.Float).SetPrec(prec)

	// BBP series: π = Σ (1/16^k) * (4/(8k+1) - 2/(8k+4) - 1/(8k+5) - 1/(8k+6))
	// The number of terms needed for 'digits' decimal places is roughly digits * log2(10) / 4 ≈ digits * 0.83
	terms := int(float64(digits) * 0.83)
	if terms < 10 {
		terms = 10
	}

	one := new(big.Float).SetPrec(prec).SetFloat64(1.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)

	for k := 0; k < terms; k++ {
		kf := new(big.Float).SetPrec(prec).SetFloat64(float64(k))

		// 16^k
		pow16 := new(big.Int).Exp(big.NewInt(16), big.NewInt(int64(k)), nil)
		pow16f := new(big.Float).SetPrec(prec).SetInt(pow16)

		// 8k+1, 8k+4, 8k+5, 8k+6
		t1 := new(big.Float).SetPrec(prec).Add(new(big.Float).Mul(big.NewFloat(8), kf), one)
		t2 := new(big.Float).SetPrec(prec).Add(new(big.Float).Mul(big.NewFloat(8), kf), four)
		t3 := new(big.Float).SetPrec(prec).Add(new(big.Float).Mul(big.NewFloat(8), kf), new(big.Float).SetFloat64(5.0))
		t4 := new(big.Float).SetPrec(prec).Add(new(big.Float).Mul(big.NewFloat(8), kf), new(big.Float).SetFloat64(6.0))

		// 4/(8k+1)
		term1 := new(big.Float).SetPrec(prec).Quo(four, t1)

		// 2/(8k+4)
		term2 := new(big.Float).SetPrec(prec).Quo(two, t2)

		// 1/(8k+5)
		term3 := new(big.Float).SetPrec(prec).Quo(one, t3)

		// 1/(8k+6)
		term4 := new(big.Float).SetPrec(prec).Quo(one, t4)

		// Sum inside parentheses
		inner := new(big.Float).SetPrec(prec)
		inner.Sub(term1, term2)
		inner.Sub(inner, term3)
		inner.Sub(inner, term4)

		// Divide by 16^k
		inner.Quo(inner, pow16f)

		// Add to total
		pi.Add(pi, inner)
	}

	return pi
}

// VerifyAndReport checks the computed π against BBP and returns a formatted
// result string suitable for display.
func VerifyAndReport(computed *big.Float, digits int, algorithmName string) string {
	if digits <= 0 {
		digits = 100
	}

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