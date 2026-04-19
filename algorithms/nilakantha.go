package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// Nilakantha computes π using Nilakantha Somayaji's alternating series:
// π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - ...
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

	start := time.Now()

	prec := uint(precision)

	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	d1 := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	d2 := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	d3 := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	// Initial term: 3 + 4/(2·3·4)
	firstTerm := new(big.Float).SetPrec(prec).Quo(four,
		new(big.Float).Mul(d1, new(big.Float).Mul(d2, d3)))
	sum := new(big.Float).SetPrec(prec).Add(three, firstTerm)

	for k := 1; k < iters; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		d1.Add(d1, two)
		d2.Add(d2, two)
		d3.Add(d3, two)

		term := new(big.Float).SetPrec(prec).Quo(four,
			new(big.Float).Mul(d1, new(big.Float).Mul(d2, d3)))

		if k%2 == 0 {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}

		// Progress
		if k > 0 && k%1000000 == 0 {
			webPrint(fmt.Sprintf("  ... %s iterations",
				pkg.FormatIntWithCommas(int64(k))))
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	showDigits := 30
	piStr := sum.Text('f', showDigits+2)
	webPrint(fmt.Sprintf("  π = %s", piStr))

	// Cross-verify
	verifyMsg := pkg.VerifyAndReport(sum, showDigits, "Nilakantha")
	webPrint(verifyMsg)

	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  Nilakantha Somayaji (1444–1544) was a mathematician")
	webPrint("  and astronomer of the Kerala school. His series")
	webPrint("  for π predates European calculus by 150 years.")
	webPrint(pkg.BoxSep(50))
}