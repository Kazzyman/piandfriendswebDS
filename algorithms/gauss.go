package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

const gaussMaxIters = 16

// GaussLegendre computes π using the Gauss-Legendre algorithm.
// Developed by Carl Friedrich Gauss, refined by Adrien-Marie Legendre.
// Quadratic convergence: correct digits double each iteration.
func GaussLegendre(webPrint func(string), iters int) {
	webPrint(pkg.BoxLine("  GAUSS-LEGENDRE ALGORITHM  ", 50))
	webPrint(pkg.BoxLine("  C.F. Gauss & A.M. Legendre  ", 50))
	webPrint(pkg.BoxLine("  c. 1800  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  The Gauss-Legendre algorithm maintains four")
	webPrint("  values that evolve with each iteration:")
	webPrint("")
	webPrint("    a = arithmetic mean (starts at 1)")
	webPrint("    b = geometric mean  (starts at 1/√2)")
	webPrint("    t = correction term (starts at 1/4)")
	webPrint("    p = power of 2      (starts at 1)")
	webPrint("")
	webPrint("  Each iteration:")
	webPrint("    a_next = (a + b) / 2")
	webPrint("    b_next = √(a × b)")
	webPrint("    t_next = t - p × (a - a_next)²")
	webPrint("    p_next = 2 × p")
	webPrint("")
	webPrint("  Then: π ≈ (a + b)² / (4 × t)")
	webPrint("")
	webPrint("  Convergence is quadratic—the number of")
	webPrint("  correct digits doubles with each iteration.")
	webPrint("  This is a consequence of the arithmetic-")
	webPrint("  geometric mean (AGM) iteration, which")
	webPrint("  converges quadratically. The error term t")
	webPrint("  shrinks at the same quadratic rate.")
	webPrint("")

	if iters > gaussMaxIters {
		webPrint(fmt.Sprintf("  Capped at %d iterations for web demo.", gaussMaxIters))
		webPrint("  Beyond this, runtime and memory grow")
		webPrint("  exponentially (literally: O(2ⁿ)).")
		iters = gaussMaxIters
	}
	if iters < 1 {
		iters = 1
	}

	webPrint(fmt.Sprintf("  Iterations requested: %d", iters))
	webPrint("")
	webPrint("  Expected digits by iteration:")
	webPrint("    Iteration  1 →      1 digit")
	webPrint("    Iteration  5 →     16 digits")
	webPrint("    Iteration 10 →    512 digits")
	webPrint("    Iteration 12 →  2,048 digits")
	webPrint("    Iteration 14 →  8,192 digits")
	webPrint("    Iteration 16 → 32,768 digits")
	webPrint("")
	webPrint("  Each iteration doubles the digit count.")
	webPrint("  Each iteration also doubles the required")
	webPrint("  precision in bits, quadrupling the work")
	webPrint("  for big.Float multiplication.")
	webPrint("")

	start := time.Now()

	// Precision: each iteration doubles correct digits
	// We allocate enough bits for the final iteration
	precBits := uint(1<<uint(iters)) * 4
	if precBits < 64 {
		precBits = 64
	}

	webPrint(fmt.Sprintf("  Precision: %d bits (~%d decimal digits)",
		precBits, precBits/4))
	webPrint("")
	webPrint("  Beginning calculation...")
	webPrint("")

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
		showAt := []int{1, 2, 4, 8, 12, 16}
		shouldShow := false
		for _, v := range showAt {
			if i == v {
				shouldShow = true
				break
			}
		}
		if shouldShow || i == iters {
			digitsToShow := 1 << uint(i)
			if digitsToShow > 200 {
				digitsToShow = 200
			}
			piStr := pi.Text('f', digitsToShow+2)
			gap := new(big.Float).SetPrec(precBits).Sub(a, b)
			webPrint(fmt.Sprintf("  Iteration %2d: %s", i, piStr))
			webPrint(fmt.Sprintf("             gap (a-b) = %s", gap.Text('e', 4)))
		}
	}

	sumAB := new(big.Float).SetPrec(precBits).Add(a, b)
	pi := new(big.Float).SetPrec(precBits).Mul(sumAB, sumAB)
	pi.Quo(pi, new(big.Float).SetPrec(precBits).Mul(four, t))

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("  RESULT:")
	webPrint("")

	digitsToShow := 3000
	piStr := pi.Text('f', digitsToShow+2)
	webPrint(fmt.Sprintf("  π = %s", piStr))

	// If more than 3000 digits were produced, note it
	expectedDigits := 1 << uint(iters)
	if expectedDigits > 3000 {
		webPrint(fmt.Sprintf("  (... and %s more correct digits)",
			pkg.FormatIntWithCommas(int64(expectedDigits-3000))))
	}

	webPrint("")
	webPrint(fmt.Sprintf("  Iterations: %d", iters))
	webPrint(fmt.Sprintf("  Expected correct digits: %s",
		pkg.FormatIntWithCommas(int64(expectedDigits))))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  Carl Friedrich Gauss developed the arithmetic-")
	webPrint("  geometric mean (AGM) around 1800. Adrien-Marie")
	webPrint("  Legendre refined it into this algorithm for π.")
	webPrint("")
	webPrint("  Quadratic convergence means each iteration")
	webPrint("  doubles the number of correct digits. The AGM")
	webPrint("  iteration has the property that a and b")
	webPrint("  converge to the same limit quadratically:")
	webPrint("  the number of matching digits doubles each step.")
	webPrint("")
	webPrint("  Gauss never published this algorithm. It was")
	webPrint("  found in his notebooks after his death.")
	webPrint(pkg.BoxSep(50))
}