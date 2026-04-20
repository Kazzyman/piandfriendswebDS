package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

const hexDisplayCap = 2000
const decimalConversionCap = 200

// BBP computes π using the Bailey–Borwein–Plouffe formula.
// It calculates π in hexadecimal and displays the hex result.
// For larger requests, it explains why conversion is impractical,
// then recomputes in base-10 from scratch and displays the decimal result.
func BBP(done chan bool, webPrint func(string), digits int) {
	webPrint(pkg.BoxLine("  BAILEY–BORWEIN–PLOUFFE (BBP)  ", 50))
	webPrint(pkg.BoxLine("  Discovered 1995 via PSLQ algorithm  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")

	// Cap computation at 50,000 digits
	maxDigits := 50000
	if digits > maxDigits {
		webPrint(fmt.Sprintf("  Requested %d digits, capped at %d.", digits, maxDigits))
		digits = maxDigits
	}

	webPrint(fmt.Sprintf("  Computing %d hexadecimal digits...", digits))
	webPrint("")

	start := time.Now()

	// Compute π in hexadecimal
	hexPi := computeBBPHex(done, webPrint, digits, start)
	if hexPi == "" {
		return // stopped by user
	}

	hexElapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("  HEXADECIMAL RESULT:")
	webPrint("")

	// Display hex
	hexDisplay := hexPi
	if len(hexDisplay) > hexDisplayCap {
		hexDisplay = hexDisplay[:hexDisplayCap]
	}
	webPrint(fmt.Sprintf("  π (base-16) = %s", hexDisplay))
	if len(hexPi) > hexDisplayCap {
		webPrint(fmt.Sprintf("  (... and %s more hex digits)",
			pkg.FormatIntWithCommas(int64(len(hexPi)-hexDisplayCap))))
	}

	webPrint("")
	webPrint(fmt.Sprintf("  Hex computation time: %s", hexElapsed.Round(time.Millisecond)))
	webPrint(pkg.BoxSep(50))
	webPrint("")

	// Decide how to handle decimal output
	if digits <= decimalConversionCap {
		// Small request: convert directly
		webPrint("  Converting hex to decimal...")
		convStart := time.Now()
		decimalPi := hexToDecimal(hexPi)
		convElapsed := time.Since(convStart)

		webPrint(fmt.Sprintf("  Conversion time: %s", convElapsed.Round(time.Millisecond)))
		webPrint("")
		webPrint(pkg.BoxSep(50))
		webPrint("  DECIMAL RESULT:")
		webPrint("")
		decimalStr := decimalPi.Text('f', digits+2)
		webPrint(fmt.Sprintf("  π (base-10) = %s", decimalStr))
		webPrint(pkg.BoxSep(50))
	} else {
		// Large request: explain and recompute in base-10
		webPrint("  Converting 50,000 hex digits to decimal is impractical.")
		webPrint("  It would require arbitrary-precision base conversion")
		webPrint("  that big.Float cannot do accurately at this scale.")
		webPrint("")
		webPrint("  Instead, we will recompute π from scratch using the")
		webPrint("  same BBP formula evaluated directly in base-10.")
		webPrint("  This is a separate, independent calculation.")
		webPrint("")
		webPrint(pkg.BoxSep(50))
		webPrint("")
		webPrint(fmt.Sprintf("  Recomputing %d decimal digits...", digits))
		webPrint("")

		decimalStart := time.Now()
		decimalPi := computeBBPDecimal(done, webPrint, digits, decimalStart)
		if decimalPi == nil {
			return // stopped by user
		}
		decimalElapsed := time.Since(decimalStart)

		webPrint("")
		webPrint(pkg.BoxSep(50))
		webPrint("  DECIMAL RESULT:")
		webPrint("")

		displayCap := 3000
		decimalStr := decimalPi.Text('f', digits+2)
		if len(decimalStr) > displayCap+2 {
			webPrint(fmt.Sprintf("  π (base-10) = %s", decimalStr[:displayCap+2]))
			webPrint(fmt.Sprintf("  (... and %s more decimal digits)",
				pkg.FormatIntWithCommas(int64(len(decimalStr)-displayCap-2))))
		} else {
			webPrint(fmt.Sprintf("  π (base-10) = %s", decimalStr))
		}

		webPrint("")
		webPrint(fmt.Sprintf("  Decimal computation time: %s", decimalElapsed.Round(time.Millisecond)))
		webPrint(pkg.BoxSep(50))
	}
	webPrint("")
	webPrint("  The BBP formula famously allows extraction of")
	webPrint("  any individual hexadecimal digit of π without")
	webPrint("  computing the preceding digits.")
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  You may have noticed that the second calculation")
	webPrint("  (base-10) ran considerably faster than the hex")
	webPrint("  calculation. The reason is precision allocation.")
	webPrint("")
	webPrint("  Hex calculation:")
	webPrint("    prec = digits × 4 + 256")
	webPrint(fmt.Sprintf("    For %s hex digits: %s × 4 + 256 = %s bits",
		pkg.FormatIntWithCommas(int64(digits)),
		pkg.FormatIntWithCommas(int64(digits)),
		pkg.FormatIntWithCommas(int64(digits*4+256))))
	webPrint("")
	webPrint("  Decimal calculation:")
	webPrint("    prec = digits × 3.33 + 64")
	webPrint(fmt.Sprintf("    For %s decimal digits: %s × 3.33 + 64 ≈ %s bits",
		pkg.FormatIntWithCommas(int64(digits)),
		pkg.FormatIntWithCommas(int64(digits)),
		pkg.FormatIntWithCommas(int64(float64(digits)*3.33+64))))
	webPrint("")
	webPrint("  The decimal calculation uses ~17% fewer bits.")
	webPrint("  big.Float operations scale with precision—")
	webPrint("  larger numbers take longer to multiply and divide.")
	webPrint("")
	webPrint("  Additionally, the hex calculation extracts each")
	webPrint("  digit individually in a loop after summation.")
	webPrint("  The decimal calculation simply formats the")
	webPrint("  accumulated result.")
	webPrint("")
	webPrint("  The deeper truth: BBP is native to hex. To get")
	webPrint("  N hex digits, you need ~4N bits of precision.")
	webPrint("  To get N decimal digits, you need only ~3.32N")
	webPrint("  bits, because one decimal digit carries less")
	webPrint("  information than one hex digit (log₂10 ≈ 3.32")
	webPrint("  bits vs 4 bits). The decimal computation is")
	webPrint("  solving a slightly 'easier' problem.")
	webPrint(pkg.BoxSep(50))
}

// computeBBPHex computes π in hexadecimal using the BBP formula.
func computeBBPHex(done chan bool, webPrint func(string), digits int, start time.Time) string {
	prec := uint(digits*4 + 256)

	pi := new(big.Float).SetPrec(prec).SetFloat64(0.0)

	one := new(big.Float).SetPrec(prec).SetFloat64(1.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)

	terms := digits + 10
	if terms < 10 {
		terms = 10
	}

	for k := 0; k < terms; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return ""
		default:
		}

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

		if k > 0 && k%1000 == 0 {
			elapsed := time.Since(start)
			pct := float64(k) / float64(terms) * 100
			webPrint(fmt.Sprintf("  ... term %s (%.1f%%) %s",
				pkg.FormatIntWithCommas(int64(k)),
				pct,
				elapsed.Round(time.Millisecond)))
		}
	}

	// Extract hex digits
	hexResult := "3."

	frac := new(big.Float).SetPrec(prec).Sub(pi, new(big.Float).SetPrec(prec).SetFloat64(3.0))

	for i := 0; i < digits; i++ {
		frac.Mul(frac, new(big.Float).SetPrec(prec).SetFloat64(16.0))
		digit, _ := frac.Int(nil)
		d := int(digit.Int64())

		if d < 10 {
			hexResult += string(rune('0' + d))
		} else {
			hexResult += string(rune('A' + d - 10))
		}

		frac.Sub(frac, new(big.Float).SetPrec(prec).SetInt(digit))
	}

	return hexResult
}

// computeBBPDecimal computes π in decimal using the BBP formula evaluated in base-10.
func computeBBPDecimal(done chan bool, webPrint func(string), digits int, start time.Time) *big.Float {
	prec := uint(float64(digits)*3.33 + 64)

	pi := new(big.Float).SetPrec(prec).SetFloat64(0.0)

	one := new(big.Float).SetPrec(prec).SetFloat64(1.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)

	terms := int(float64(digits) * 0.83)
	if terms < 10 {
		terms = 10
	}

	for k := 0; k < terms; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return nil
		default:
		}

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

		if k > 0 && k%1000 == 0 {
			elapsed := time.Since(start)
			pct := float64(k) / float64(terms) * 100
			webPrint(fmt.Sprintf("  ... term %s (%.1f%%) %s",
				pkg.FormatIntWithCommas(int64(k)),
				pct,
				elapsed.Round(time.Millisecond)))
		}
	}

	return pi
}

// hexToDecimal converts a hexadecimal string (like "3.243F6A...") to a big.Float.
// Only used for small digit counts where conversion is accurate.
func hexToDecimal(hexStr string) *big.Float {
	result := new(big.Float).SetFloat64(3.0)

	if len(hexStr) > 2 && hexStr[1] == '.' {
		fracDigits := hexStr[2:]
		pow16 := new(big.Float).SetFloat64(1.0 / 16.0)

		for _, ch := range fracDigits {
			var val float64
			if ch >= '0' && ch <= '9' {
				val = float64(ch - '0')
			} else if ch >= 'A' && ch <= 'F' {
				val = float64(ch - 'A' + 10)
			} else if ch >= 'a' && ch <= 'f' {
				val = float64(ch - 'a' + 10)
			} else {
				continue
			}

			term := new(big.Float).Mul(pow16, new(big.Float).SetFloat64(val))
			result.Add(result, term)
			pow16.Mul(pow16, new(big.Float).SetFloat64(1.0/16.0))
		}
	}

	return result
}