package algorithms

import (
	"fmt"
	"time"

	"piandfriends/pkg"
)

// Gregory4 computes π using the Gregory-Leibniz series multiplied by 4:
// π = 4 - 4/3 + 4/5 - 4/7 + 4/9 - ...
//
// This is the same series as GregoryLeibniz (π/4 = 1 - 1/3 + 1/5 - ...)
// but written in a form that directly yields π.
//
// This version runs 5 million iterations—enough to see convergence
// in action without waiting hours. For the full stress-test experience,
// see the GregoryLeibniz method (40 billion iterations).
func Gregory4(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  GREGORY–LEIBNIZ (×4 FORM)  ", 50))
	webPrint(pkg.BoxLine("  π = 4 - 4/3 + 4/5 - 4/7 + ...  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  This is the Gregory-Leibniz series multiplied")
	webPrint("  by 4, yielding π directly.")
	webPrint("")
	webPrint("  Why two Gregory-Leibniz methods in this suite?")
	webPrint("")
	webPrint("  1. This ×4 form appears in some textbooks and")
	webPrint("     is a convenient direct expression for π.")
	webPrint("")
	webPrint("  2. The other Gregory-Leibniz method runs 40")
	webPrint("     billion iterations as a hardware stress test,")
	webPrint("     starkly contrasting with faster algorithms")
	webPrint("     like Chudnovsky and Gauss-Legendre.")
	webPrint("")
	webPrint("  This version runs 5 million iterations—enough")
	webPrint("  to watch convergence without waiting hours.")
	webPrint("")

	const totalIters = 5000000

	webPrint(fmt.Sprintf("  Running %s iterations...",
		pkg.FormatIntWithCommas(int64(totalIters))))
	webPrint("")

	start := time.Now()

	// Start with first term: 4
	sum := 4.0
	denom := 3.0
	sign := -1.0

	lastReport := 0
	reportInterval := totalIters / 20 // Report at 5% intervals

	for iter := 1; iter <= totalIters; iter++ {
		select {
		case <-done:
			webPrint("  Calculation stopped.")
			return
		default:
		}

		// Add next term: sign * (4/denom)
		sum += sign * (4.0 / denom)
		denom += 2.0
		sign = -sign

		// Progress reporting
		if iter-lastReport >= reportInterval || iter == totalIters {
			lastReport = iter
			elapsed := time.Since(start)
			pct := float64(iter) / float64(totalIters) * 100

			// Estimate correct digits
			estDigits := int(float64(iter) / 1000000.0 * 1.5)
			if estDigits > 7 {
				estDigits = 7
			}

			webPrint(fmt.Sprintf("  ... %s terms (%.1f%%) ~%d digits %s",
				pkg.FormatIntWithCommas(int64(iter)),
				pct,
				estDigits,
				elapsed.Round(time.Millisecond)))

			if iter%(reportInterval*2) == 0 || iter == totalIters {
				webPrint(fmt.Sprintf("      π ≈ %.10f", sum))
			}
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("  RESULT:")
	webPrint("")
	webPrint(fmt.Sprintf("  π ≈ %.12f", sum))
	webPrint("  Known π = 3.141592653589...")
	webPrint("")

	// Count correct digits
	piStr := fmt.Sprintf("%.12f", sum)
	refStr := "3.141592653589"
	correct := 0
	for i := 0; i < len(piStr) && i < len(refStr); i++ {
		if piStr[i] != refStr[i] {
			break
		}
		correct++
	}
	if correct > 2 {
		correct -= 2
	} else {
		correct = 0
	}

	webPrint(fmt.Sprintf("  Correct decimal digits: %d", correct))
	webPrint(fmt.Sprintf("  Terms computed: %s",
		pkg.FormatIntWithCommas(int64(totalIters))))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  For comparison: the full Gregory-Leibniz stress")
	webPrint("  test runs 40 billion iterations—8,000 times")
	webPrint("  more terms—and yields only 3-4 additional")
	webPrint("  correct digits, taking hours instead of seconds.")
	webPrint("")
	webPrint("  This is the brutal reality of slow convergence.")
	webPrint(pkg.BoxSep(50))
}