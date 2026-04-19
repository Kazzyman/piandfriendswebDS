package algorithms

import (
	"fmt"
	"time"

	"piandfriends/pkg"
)

// Wallis computes π using John Wallis's infinite product (1655):
// π/2 = (2/1)*(2/3) * (4/3)*(4/5) * (6/5)*(6/7) * ...
// Convergence is notoriously slow—a good hardware speed test.
func Wallis(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  JOHN WALLIS INFINITE PRODUCT  ", 50))
	webPrint(pkg.BoxLine("  c. 1655  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  π/2 = (2/1)·(2/3) · (4/3)·(4/5) · (6/5)·(6/7) ...")
	webPrint("")
	webPrint("  Running 40 billion iterations...")
	webPrint("  This is a hardware speed test.")
	webPrint("")

	start := time.Now()

	numer := 2.0
	denom1 := 1.0
	denom2 := 3.0
	product := (numer / denom1) * (numer / denom2)

	var iter int64

	// Phase 1: 1 billion iterations
	for iter < 1000000000 {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		iter++
		numer += 2
		denom1 += 2
		denom2 += 2
		product *= (numer / denom1) * (numer / denom2)

		if iter == 10000 {
			pi := product * 2
			webPrint(fmt.Sprintf("  10,000: π ≈ %.6f", pi))
		}
		if iter == 50000 {
			pi := product * 2
			webPrint(fmt.Sprintf("  50,000: π ≈ %.7f", pi))
		}
		if iter == 500000 {
			pi := product * 2
			webPrint(fmt.Sprintf("  500,000: π ≈ %.8f", pi))
		}
		if iter == 2000000 {
			pi := product * 2
			webPrint(fmt.Sprintf("  2,000,000: π ≈ %.9f", pi))
		}
		if iter == 40000000 {
			pi := product * 2
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  40,000,000: π ≈ %.10f (%s)",
				pi, elapsed.Round(time.Millisecond)))
		}
		if iter == 400000000 {
			pi := product * 2
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  400,000,000: π ≈ %.10f (%s)",
				pi, elapsed.Round(time.Millisecond)))
		}
		if iter == 1000000000 {
			pi := product * 2
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  1,000,000,000: π ≈ %.11f (%s)",
				pi, elapsed.Round(time.Millisecond)))
			webPrint("  ... continuing to 40 billion ...")
		}
	}

	// Phase 2: up to 40 billion
	for iter < 40000000000 {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		iter++
		numer += 2
		denom1 += 2
		denom2 += 2
		product *= (numer / denom1) * (numer / denom2)

		if iter%10000000000 == 0 {
			pi := product * 2
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  %s: π ≈ %.12f (%s)",
				pkg.FormatIntWithCommas(iter),
				pi,
				elapsed.Round(time.Millisecond)))
		}
	}

	pi := product * 2
	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint(fmt.Sprintf("  Final after 40 billion iterations:"))
	webPrint(fmt.Sprintf("  π ≈ %.12f", pi))
	webPrint(fmt.Sprintf("  math.Pi = %.12f", 3.141592653589793))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  Wallis's product converges agonizingly slowly.")
	webPrint("  After 40 billion iterations, we have ~10 digits.")
	webPrint("  But in 1655, this was revolutionary.")
	webPrint(pkg.BoxSep(50))
}