package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// Nilakantha computes π using Nilakantha Somayaji's alternating series:
// π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - 4/(8·9·10) + ...
// Discovered in Kerala, India, c. 1530—150 years before Newton.
func Nilakantha(done chan bool, webPrint func(string), iters, precision int) {
	webPrint(pkg.BoxLine("  NILAKANTHA SOMAYAJI'S SERIES  ", 50))
	webPrint(pkg.BoxLine("  Kerala school, c. 1530  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint(fmt.Sprintf("  Iterations: %s", pkg.FormatIntWithCommas(int64(iters))))
	webPrint(fmt.Sprintf("  Precision: %d bits", precision))
	webPrint("")

	if iters > 1000000000 {
		webPrint("  Cannot exceed 1 billion iterations.")
		return
	}
	if iters < 1 {
		iters = 1
	}

	start := time.Now()
	prec := uint(precision)

	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	// k=1 term: 4/(2·3·4) = 4/24 = 1/6
	firstTerm := new(big.Float).SetPrec(prec).Quo(four, new(big.Float).SetPrec(prec).SetFloat64(2*3*4))
	pi := new(big.Float).SetPrec(prec).Add(three, firstTerm)

	// Loop for k=2 to iters
	for k := 2; k <= iters; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		// Calculate denominators: (2k), (2k+1), (2k+2)
		d1 := new(big.Float).SetPrec(prec).SetFloat64(float64(2 * k))
		d2 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 1))
		d3 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 2))

		// Denominator product using big.Float
		denom := new(big.Float).SetPrec(prec).Mul(d1, new(big.Float).SetPrec(prec).Mul(d2, d3))

		// Term = 4 / denom
		term := new(big.Float).SetPrec(prec).Quo(four, denom)

		// k=2 subtract, k=3 add, k=4 subtract, etc.
		if k%2 == 0 {
			pi.Sub(pi, term)
		} else {
			pi.Add(pi, term)
		}

		// Progress reporting
		if k%100000 == 0 {
			elapsed := time.Since(start)
			pct := float64(k) / float64(iters) * 100
			webPrint(fmt.Sprintf("  ... %s terms (%.1f%%) %s",
				pkg.FormatIntWithCommas(int64(k)),
				pct,
				elapsed.Round(time.Millisecond)))
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	showDigits := 32
	if showDigits > precision/3 {
		showDigits = precision / 3
	}
	piStr := pi.Text('f', showDigits)
	webPrint(fmt.Sprintf("  π = %s", piStr))

	// Cross-verify with BBP
	// Nilakantha converges slowly: error ≈ 1/(2*N³)
	// For 1M iterations, expect ~16-18 correct digits.
	// Verify conservatively at 15 digits.
	verifyDigits := 15
	if verifyDigits > showDigits {
		verifyDigits = showDigits
	}
	if showDigits < 20 {
		verifyDigits = showDigits
	}
	verifyMsg := pkg.VerifyAndReport(pi, verifyDigits, "Nilakantha")
	webPrint(verifyMsg)

	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  Nilakantha Somayaji (1444–1544) was a mathematician")
	webPrint("  and astronomer of the Kerala school. His series")
	webPrint("  for π predates European calculus by 150 years.")
	webPrint(pkg.BoxSep(50))
}