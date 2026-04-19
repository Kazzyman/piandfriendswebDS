package algorithms

import (
	"fmt"
	"time"

	"piandfriends/pkg"
)

// CustomSeries is Rick's quick-serve π approximation.
// π = (4/1) - (4/3) + (4/5) - (4/7) + ...
// Gets ~9 digits in seconds.
func CustomSeries(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  CUSTOM SERIES  ", 50))
	webPrint(pkg.BoxLine("  Rick's quick-serve recipe  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  π = 4/1 - 4/3 + 4/5 - 4/7 + 4/9 - 4/11 + ...")
	webPrint("")
	webPrint("  Running 300 million iterations...")
	webPrint("")

	start := time.Now()

	tally := 4.0 / 1.0
	nextOdd := 3.0
	var iter int64 = 1

	for iter < 300000000 {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		iter++
		tally -= 4.0 / nextOdd
		nextOdd += 2
		tally += 4.0 / nextOdd
		nextOdd += 2

		if iter == 10000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  10M: π ≈ %0.6f (%s)", tally, elapsed.Round(time.Millisecond)))
		}
		if iter == 50000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  50M: π ≈ %0.8f (%s)", tally, elapsed.Round(time.Millisecond)))
		}
		if iter == 100000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  100M: π ≈ %0.9f (%s)", tally, elapsed.Round(time.Millisecond)))
		}
		if iter == 200000000 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  200M: π ≈ %0.10f (%s)", tally, elapsed.Round(time.Millisecond)))
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint(fmt.Sprintf("  Final after 300 million iterations:"))
	webPrint(fmt.Sprintf("  π ≈ %0.11f", tally))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  Not the fastest. Not the slowest.")
	webPrint("  Just Rick's quick-serve recipe.")
	webPrint(pkg.BoxSep(50))
}