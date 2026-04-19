package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

const gaussMaxIters = 12

// GaussLegendre computes π using the Gauss-Legendre algorithm.
// Quadratic convergence: correct digits double each iteration.
func GaussLegendre(webPrint func(string), iters int) {
	if iters > gaussMaxIters {
		webPrint(fmt.Sprintf("  Maximum iterations is %d for web demo.", gaussMaxIters))
		return
	}
	if iters < 1 {
		iters = 1
	}

	precBits := uint(1<<uint(iters)) * 4
	if precBits < 64 {
		precBits = 64
	}

	webPrint(pkg.BoxLine("  GAUSS-LEGENDRE ALGORITHM  ", 50))
	webPrint(pkg.BoxLine("  Quadratic convergence  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint(fmt.Sprintf("  Iterations: %d", iters))
	webPrint(fmt.Sprintf("  Precision: %d bits (~%d decimal digits)", precBits, precBits/4))
	webPrint("")

	start := time.Now()

	two := new(big.Float).SetPrec(precBits).SetFloat64(2.0)
	four := new(big.Float).SetPrec(precBits).SetFloat64(4.0)

	a := new(big.Float).SetPrec(precBits).SetFloat64(1.0)

	b := new(big.Float).SetPrec(precBits).Sqrt(two)
	b.Quo(new(big.Float).SetPrec(precBits).SetFloat64(1.0), b)

	t := new(big.Float).SetPrec(precBits).SetFloat64(0.25)
	p := new(big.Float).SetPrec(precBits).SetFloat64(1.0)

	for i := 1; i <= iters; i++ {
		aNext := new(big.Float).SetPrec(precBits).Add(a, b)
		aNext.Quo(aNext, two)

		bNext := new(big.Float).SetPrec(precBits).Mul(a, b)
		bNext.Sqrt(bNext)

		diff := new(big.Float).SetPrec(precBits).Sub(a, aNext)
		diff.Mul(diff, diff)
		diff.Mul(p, diff)

		tNext := new(big.Float).SetPrec(precBits).Sub(t, diff)
		pNext := new(big.Float).SetPrec(precBits).Mul(two, p)

		a, b, t, p = aNext, bNext, tNext, pNext

		sumAB := new(big.Float).SetPrec(precBits).Add(a, b)
		pi := new(big.Float).SetPrec(precBits).Mul(sumAB, sumAB)
		pi.Quo(pi, new(big.Float).SetPrec(precBits).Mul(four, t))

		// Show progress
		if i == 1 || i == 2 || i == 4 || i == 8 || i == iters {
			digits := i * 2
			if digits > 100 {
				digits = 100
			}
			piStr := pi.Text('f', digits+2)
			webPrint(fmt.Sprintf("  Iteration %2d: %s", i, piStr))
		}
	}

	sumAB := new(big.Float).SetPrec(precBits).Add(a, b)
	pi := new(big.Float).SetPrec(precBits).Mul(sumAB, sumAB)
	pi.Quo(pi, new(big.Float).SetPrec(precBits).Mul(four, t))

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	// Determine how many digits to show and verify
	digitsToShow := 100
	if iters >= 8 {
		digitsToShow = 200
	}
	if iters >= 10 {
		digitsToShow = 500
	}
	if iters >= 12 {
		digitsToShow = 1000
	}

	piStr := pi.Text('f', digitsToShow+2)

	// Cross-verify with BBP
	verifyMsg := pkg.VerifyAndReport(pi, digitsToShow, "Gauss-Legendre")

	webPrint(fmt.Sprintf("  Final after %d iterations:", iters))
	webPrint(fmt.Sprintf("  π = %s", piStr))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint(verifyMsg)
	webPrint(pkg.BoxSep(50))
}