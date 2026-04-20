package algorithms

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// MonteCarlo estimates π using random sampling.
// Darts thrown at a unit square, counting those inside the quarter circle.
// Inefficient but beautiful—order emerging from chaos.
//
// This is Rick's second-favorite method.
func MonteCarlo(webPrint func(string), gridSize int) {
	webPrint(pkg.BoxLine("  MONTE CARLO METHOD  ", 50))
	webPrint(pkg.BoxLine("  Rick's Second Favorite  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  Imagine throwing darts randomly at a square")
	webPrint("  board with a quarter-circle inscribed inside.")
	webPrint("  The ratio of darts landing inside the circle")
	webPrint("  to the total thrown converges to π/4.")
	webPrint("")
	webPrint("  Area of quarter-circle: π/4")
	webPrint("  Area of square: 1")
	webPrint("  Ratio = π/4  →  π = 4 × (inside / total)")
	webPrint("")
	webPrint("  Instead of darts, we sample points on a grid.")
	webPrint("  For each point (x, y) in the unit square,")
	webPrint("  if x² + y² ≤ 1, the point is inside.")
	webPrint("")

	// Cap grid size for web demo
	maxGrid := 50000
	if gridSize > maxGrid {
		webPrint(fmt.Sprintf("  Grid size capped at %s for web demo.",
			pkg.FormatIntWithCommas(int64(maxGrid))))
		gridSize = maxGrid
	}
	if gridSize < 100 {
		gridSize = 100
	}

	totalPoints := gridSize * gridSize
	webPrint(fmt.Sprintf("  Grid size: %s × %s = %s points",
		pkg.FormatIntWithCommas(int64(gridSize)),
		pkg.FormatIntWithCommas(int64(gridSize)),
		pkg.FormatIntWithCommas(int64(totalPoints))))
	webPrint("")

	if totalPoints > 100000000 {
		webPrint("  That's a lot of darts. This may take a while...")
		webPrint("")
	}

	start := time.Now()

	inside := big.NewInt(0)

	// Step size for grid sampling
	step := 1.0 / float64(gridSize)
	halfStep := step / 2.0

	reportInterval := gridSize / 10
	if reportInterval < 1 {
		reportInterval = 1
	}

	for i := 0; i < gridSize; i++ {
		x := float64(i)*step + halfStep

		for j := 0; j < gridSize; j++ {
			y := float64(j)*step + halfStep

			if x*x+y*y <= 1.0 {
				inside.Add(inside, big.NewInt(1))
			}
		}

		// Progress reporting
		if i > 0 && i%reportInterval == 0 {
			elapsed := time.Since(start)
			pct := float64(i) / float64(gridSize) * 100
			webPrint(fmt.Sprintf("  ... %.0f%% complete (%s)",
				pct, elapsed.Round(time.Millisecond)))
		}
	}

	// π ≈ 4 × (inside / total)
	insideFloat := new(big.Float).SetInt(inside)
	totalFloat := new(big.Float).SetInt(big.NewInt(int64(totalPoints)))
	ratio := new(big.Float).Quo(insideFloat, totalFloat)
	pi := new(big.Float).Mul(ratio, big.NewFloat(4.0))

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("  RESULT:")
	webPrint("")

	piStr := pi.Text('f', 15)
	webPrint(fmt.Sprintf("  Estimated π: %s", piStr))

	piFloat, _ := pi.Float64()
	webPrint(fmt.Sprintf("  math.Pi:     %0.15f", math.Pi))
	webPrint(fmt.Sprintf("  Difference:  %0.15f", math.Abs(piFloat-math.Pi)))

	// Count correct digits
	correct := 0
	ref := "3.141592653589793"
	piCheck := pi.Text('f', 15)
	for i := 0; i < len(ref) && i < len(piCheck); i++ {
		if piCheck[i] != ref[i] {
			break
		}
		if i >= 2 {
			correct++
		}
	}

	webPrint(fmt.Sprintf("  Correct digits: %d", correct))
	webPrint(fmt.Sprintf("  Points sampled: %s",
		pkg.FormatIntWithCommas(int64(totalPoints))))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  The Monte Carlo method was pioneered during")
	webPrint("  the Manhattan Project in the 1940s by")
	webPrint("  Stanislaw Ulam and John von Neumann.")
	webPrint("")
	webPrint("  Ulam conceived the idea while recovering from")
	webPrint("  illness and playing solitaire. He wondered")
	webPrint("  what the probability of winning a particular")
	webPrint("  solitaire layout was, and realized he could")
	webPrint("  simply simulate many random games.")
	webPrint("")
	webPrint("  They named it after the Monte Carlo casino")
	webPrint("  in Monaco—because at its heart, it is a game")
	webPrint("  of chance.")
	webPrint("")
	webPrint("  Why is it Rick's second favorite?")
	webPrint("  Because no equations, no series, no clever")
	webPrint("  algebra. Just pure randomness, applied")
	webPrint("  repeatedly, converging on one of the most")
	webPrint("  profound constants in mathematics.")
	webPrint("")
	webPrint("  It is inefficient. It is beautiful.")
	webPrint(pkg.BoxSep(50))
}