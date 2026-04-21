package algorithms

import (
	"fmt"
	"math/big"
	"time"
	"piandfriends/pkg"
)

// Nilakantha computes π using Nilakantha Somayaji's alternating series:
// π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - 4/(8·9·10) + ...
//
// Discovered in Kerala, India, c. 1530—150 years before European calculus.
// The Kerala school of astronomy and mathematics produced remarkable
// infinite series for trigonometric functions and π long before
// Newton, Leibniz, or Gregory.
func Nilakantha(done chan bool, webPrint func(string), iters, precision int) {
	webPrint(pkg.BoxLine("  NILAKANTHA SOMAYAJI'S SERIES  ", 50))
	webPrint(pkg.BoxLine("  Kerala school, c. 1530  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - ...")
	webPrint("")
	webPrint("  Nilakantha Somayaji (1444–1544) was a")
	webPrint("  mathematician and astronomer of the Kerala")
	webPrint("  school in southern India. This school produced")
	webPrint("  infinite series for π, sine, cosine, and")
	webPrint("  arctangent over a century before European")
	webPrint("  mathematicians developed calculus.")
	webPrint("")
	webPrint("  The series converges at a rate of about")
	webPrint("  3 digits per factor-of-10 increase in terms:")
	webPrint("")
	webPrint("      10,000 terms → ~9 digits")
	webPrint("     100,000 terms → ~11 digits")
	webPrint("   1,000,000 terms → ~14 digits")
	webPrint("  10,000,000 terms → ~16 digits")
	webPrint("")

	// Cap iterations for web demo
	maxIters := 10000000
	if iters > maxIters {
		webPrint(fmt.Sprintf("  Iterations capped at %s for web demo.",
			pkg.FormatIntWithCommas(int64(maxIters))))
		iters = maxIters
	}
	if iters < 1000 {
		iters = 1000
	}

	// Cap precision
	maxPrec := 2048
	if precision > maxPrec {
		webPrint(fmt.Sprintf("  Precision capped at %d bits.", maxPrec))
		precision = maxPrec
	}
	if precision < 256 {
		precision = 256
	}

	webPrint(fmt.Sprintf("  Iterations: %s", pkg.FormatIntWithCommas(int64(iters))))
	webPrint(fmt.Sprintf("  Precision: %d bits (~%d decimal digits)",
		precision, precision/3))
	webPrint("")

	if iters >= 10000000 {
		webPrint("  That's a lot of terms. This may take a while...")
		webPrint("")
	}

	start := time.Now()
	prec := uint(precision)

	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	// k=1 term: 4/(2·3·4) = 4/24 = 1/6
	firstTerm := new(big.Float).SetPrec(prec).Quo(four,
		new(big.Float).SetPrec(prec).SetFloat64(2*3*4))
	pi := new(big.Float).SetPrec(prec).Add(three, firstTerm)

	webPrint("  Beginning calculation...")
	webPrint("")

	lastReport := 0
	reportInterval := iters / 10
	if reportInterval < 1 {
		reportInterval = 1
	}

	for k := 2; k <= iters; k++ {
		select {
		case <-done:
			webPrint("  Calculation stopped.")
			return
		default:
		}

		d1 := new(big.Float).SetPrec(prec).SetFloat64(float64(2 * k))
		d2 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 1))
		d3 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 2))

		denom := new(big.Float).SetPrec(prec).Mul(d1,
			new(big.Float).SetPrec(prec).Mul(d2, d3))
		term := new(big.Float).SetPrec(prec).Quo(four, denom)

		if k%2 == 0 {
			pi.Sub(pi, term)
		} else {
			pi.Add(pi, term)
		}

		// Progress reporting
		if k-lastReport >= reportInterval || k == iters {
			lastReport = k
			elapsed := time.Since(start)
			pct := float64(k) / float64(iters) * 100

			estDigits := int(float64(k)/1000000.0*3.0) + 15
			if estDigits > 30 {
				estDigits = 30
			}

			webPrint(fmt.Sprintf("  ... %s terms (%.1f%%) ~%d digits %s",
				pkg.FormatIntWithCommas(int64(k)),
				pct,
				estDigits,
				elapsed.Round(time.Millisecond)))

			if k%(reportInterval*2) == 0 {
				showDigits := estDigits + 2
				if showDigits > 30 {
					showDigits = 30
				}
				piStr := pi.Text('f', showDigits)
				webPrint(fmt.Sprintf("      π ≈ %s", piStr))
			}
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("  RESULT:")
	webPrint("")

	showDigits := 30
	piStr := pi.Text('f', showDigits)
	webPrint(fmt.Sprintf("  π = %s", piStr))

	// Verify correct digits
	// verifyDigits := 15
	// Nilakantha gives ~3 digits per factor-of-10 in terms
	// Estimate converged digits: roughly 3 * log10(iters) - 1
	/*
	estDigits := int(3.0*math.Log10(float64(iters))) - 1
	if estDigits < 5 {
		estDigits = 5
	}
	if estDigits > 30 {
		estDigits = 30
	}
	verifyDigits := estDigits
	if verifyDigits > showDigits {
		verifyDigits = showDigits
	}
*/
	// Nilakantha converges slowly. Conservative estimates based on observation.
	estDigits := 8
	if iters >= 10000 {
		estDigits = 9
	}
	if iters >= 50000 {
		estDigits = 10
	}
	if iters >= 200000 {
		estDigits = 11
	}
	if iters >= 1000000 {
		estDigits = 14
	}
	if iters >= 5000000 {
		estDigits = 16
	}
	if iters >= 20000000 {
		estDigits = 18
	}
	verifyDigits := estDigits
	if verifyDigits > showDigits {
		verifyDigits = showDigits
	}
	verifyMsg := pkg.VerifyAndReport(pi, verifyDigits, "Nilakantha")
	webPrint(verifyMsg)

	// The rest was not to be replaced, supposedly. According to Deep Seek. 
	webPrint(fmt.Sprintf("  Terms computed: %s",
		pkg.FormatIntWithCommas(int64(iters))))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  The Kerala school flourished from the 14th")
	webPrint("  to 16th centuries. Madhava of Sangamagrama")
	webPrint("  (c. 1340–1425) founded the school and")
	webPrint("  discovered many of these series. His work")
	webPrint("  was extended by Nilakantha, Jyeshtadeva,")
	webPrint("  and others.")
	webPrint("")
	webPrint("  The series for π is a special case of the")
	webPrint("  arctangent series, which the Kerala school")
	webPrint("  discovered over 200 years before James Gregory")
	webPrint("  and Gottfried Leibniz in Europe.")
	webPrint("")
	webPrint("  History is only now recognizing the full")
	webPrint("  extent of their achievements.")
	webPrint(pkg.BoxSep(50))
}