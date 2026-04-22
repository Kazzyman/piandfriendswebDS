package algorithms

import (
	"fmt"
	"math"

	"piandfriends/pkg"
)

// ErdosBorwein computes the Erdős–Borwein constant:
// E = Σ 1/(2^n - 1) for n=1 to ∞
func ErdosBorwein(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  ERDŐS–BORWEIN CONSTANT  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  E = Σ 1/(2^n - 1)  for n=1 to ∞")
	webPrint("  Sum of reciprocals of Mersenne numbers.")
	webPrint("")

	// var sum float64 = 1.0 // n=1 term: 1/(2-1) = 1
	var sum = 1.0 // n=1 term: 1/(2-1) = 1

	webPrint("  Converging...")
	webPrint("")

	for n := 2; n <= 100; n++ {
		select {
		case <-done:
			return
		default:
		}

		sum += 1.0 / (math.Pow(2.0, float64(n)) - 1.0)

		if n == 10 || n == 20 || n == 30 || n == 50 || n == 100 {
			webPrint(fmt.Sprintf("  n=%3d: %0.25f", n, sum))
		}
	}

	webPrint("")
	webPrint(fmt.Sprintf("  Final (n=100): %0.25f", sum))
	webPrint("  Reference: 1.606695152415291763...")
	webPrint("")
	webPrint(pkg.BoxLine("  PAUL ERDŐS (1913–1996)  ", 50))
	webPrint(pkg.BoxLine("  Prolific mathematician, eccentric genius  ", 50))
	webPrint(pkg.BoxLine("  ~1,500 papers, 500+ collaborators  ", 50))
	webPrint(pkg.BoxLine("  The Erdős number: degrees of co-authorship  ", 50))
	webPrint(pkg.BoxSep(50))
}