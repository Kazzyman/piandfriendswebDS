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
	webPrint(fmt.Sprintf("  Grid size: %d × %d = %s points",
		gridSize, gridSize,
		pkg.FormatIntWithCommas(int64(gridSize*gridSize))))
	webPrint("")

	if gridSize > 119999 {
		webPrint("  That grid size makes me puke!")
		webPrint("  Please choose something smaller than 120,000.")
		return
	}

	start := time.Now()

	inside := big.NewInt(0)
	total := big.NewInt(int64(gridSize * gridSize))

	// Step size for grid sampling
	step := 1.0 / float64(gridSize)
	halfStep := step / 2.0

	for i := 0; i < gridSize; i++ {
		x := float64(i)*step + halfStep

		for j := 0; j < gridSize; j++ {
			y := float64(j)*step + halfStep

			if x*x+y*y <= 1.0 {
				inside.Add(inside, big.NewInt(1))
			}
		}

		// Progress every 10% of rows
		if gridSize > 100 && i%(gridSize/10) == 0 && i > 0 {
			pct := float64(i) / float64(gridSize) * 100
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  ... %.0f%% complete (%s)",
				pct, elapsed.Round(time.Millisecond)))
		}
	}

	// π ≈ 4 × (inside / total)
	insideFloat := new(big.Float).SetInt(inside)
	totalFloat := new(big.Float).SetInt(total)
	ratio := new(big.Float).Quo(insideFloat, totalFloat)
	pi := new(big.Float).Mul(ratio, big.NewFloat(4.0))

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	piStr := pi.Text('f', 15)
	webPrint(fmt.Sprintf("  Estimated π: %s", piStr))

	// Show float64 version for comparison
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
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Points sampled: %s",
		pkg.FormatIntWithCommas(int64(gridSize*gridSize))))
	webPrint("")
	webPrint("  The Monte Carlo method was pioneered during")
	webPrint("  the Manhattan Project by Ulam and von Neumann.")
	webPrint("  It is beautiful precisely because it is inefficient.")
	webPrint(pkg.BoxSep(50))
}