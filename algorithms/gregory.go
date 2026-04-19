package algorithms

import (
	"fmt"
	"time"

	"piandfriends/pkg"
)

// GregoryLeibniz computes π using the Gregory-Leibniz series:
// π/4 = 1 - 1/3 + 1/5 - 1/7 + 1/9 - ...
// Also converges very slowly.
func GregoryLeibniz(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  GREGORY–LEIBNIZ SERIES  ", 50))
	webPrint(pkg.BoxLine("  π/4 = 1 - 1/3 + 1/5 - 1/7 + ...  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  Running 9 billion iterations...")
	webPrint("")

	start := time.Now()

	denom := 3.0
	sum := 1.0 - 1.0/denom
	var iter int64 = 1

	for iter < 9000000000 {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		iter++
		denom += 2

		if iter%2 == 0 {
			sum += 1.0 / denom
		} else {
			sum -= 1.0 / denom
		}

		pi := 4 * sum

		if iter == 100000000 {
			webPrint(fmt.Sprintf("  100M: π ≈ %.7f", pi))
		}
		if iter == 400000000 {
			webPrint(fmt.Sprintf("  400M: π ≈ %.10f", pi))
		}
		if iter == 1000000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  1B: π ≈ %.10f (%s)", pi, elapsed.Round(time.Millisecond)))
		}
		if iter == 4000000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  4B: π ≈ %.10f (%s)", pi, elapsed.Round(time.Millisecond)))
		}
		if iter == 9000000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  9B: π ≈ %.13f (%s)", pi, elapsed.Round(time.Millisecond)))
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint(fmt.Sprintf("  Final after 9 billion iterations:"))
	webPrint(fmt.Sprintf("  π ≈ 3.1415926535..."))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  James Gregory (1638–1675) and Gottfried Wilhelm")
	webPrint("  Leibniz (1646–1716) independently discovered this")
	webPrint("  series. It is beautiful but impractical.")
	webPrint(pkg.BoxSep(50))
}