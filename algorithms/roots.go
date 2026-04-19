package algorithms

import (
	"fmt"
	"math"
	"sort"
	"time"

	"piandfriends/pkg"
)

// Roots approximates square and cube roots using integer arithmetic only.
// It builds a table of perfect powers and finds bracketing ratios.
//
// Original concept and algorithm: Richard (Rick) Woolley.
func Roots(webPrint func(string), radical, workpiece int) {
	if radical != 2 && radical != 3 {
		webPrint("  Error: radical must be 2 or 3.")
		return
	}

	webPrint("")
	webPrint(pkg.BoxLine("  ROOTS DEMO — INTEGER ARITHMETIC ONLY  ", 50))
	webPrint(pkg.BoxSep(50))

	if radical == 2 {
		webPrint(fmt.Sprintf("  Finding √%d", workpiece))
	} else {
		webPrint(fmt.Sprintf("  Finding ∛%d", workpiece))
	}
	webPrint("  Method: perfect power ratio bracketing")
	webPrint("  Original algorithm: Richard Woolley")
	webPrint("")

	start := time.Now()

	// Build table of perfect powers
	type Pair struct {
		product int
		root    int
	}
	var table []Pair

	for root := 2; root < 825000; root++ {
		var product int
		if radical == 2 {
			product = root * root
		} else {
			product = root * root * root
		}
		table = append(table, Pair{product, root})
	}

	webPrint(fmt.Sprintf("  Table built: %d perfect powers.", len(table)))
	webPrint("")

	// Search for bracketing ratios
	type Result struct {
		value         float64
		pdiff         float64
		largerPP      int
		smallerPP     int
		rootOfLarger  int
		rootOfSmaller int
	}
	var results []Result

	precision := 500
	if radical == 3 {
		precision = workpiece * 12
		if precision < 600 {
			precision = 600
		}
		if precision > 3000 {
			precision = 3000
		}
	}

	webPrint(fmt.Sprintf("  Precision window: %d", precision))
	webPrint("  Searching...")
	webPrint("")

	for i := 0; i < len(table)-1; i++ {
		smaller := table[i]

		for j := i + 1; j < len(table); j++ {
			larger := table[j]
			target := smaller.product * workpiece

			if larger.product > target {
				dL := larger.product - target
				dS := target - table[j-1].product

				if dL < precision {
					results = append(results, Result{
						value:         float64(larger.root) / float64(smaller.root),
						pdiff:         float64(dL) / float64(larger.product),
						largerPP:      larger.product,
						smallerPP:     smaller.product,
						rootOfLarger:  larger.root,
						rootOfSmaller: smaller.root,
					})
				}
				if dS < precision && j-1 > i {
					results = append(results, Result{
						value:         float64(table[j-1].root) / float64(smaller.root),
						pdiff:         float64(dS) / float64(table[j-1].product),
						largerPP:      table[j-1].product,
						smallerPP:     smaller.product,
						rootOfLarger:  table[j-1].root,
						rootOfSmaller: smaller.root,
					})
				}
				break
			}
		}

		// Progress indicator
		if i%100000 == 0 && i > 0 {
			elapsed := time.Since(start)
			webPrint(fmt.Sprintf("  ... %s iterations, %s",
				pkg.FormatIntWithCommas(int64(i)),
				elapsed.Round(time.Millisecond)))
		}
	}

	if len(results) == 0 {
		webPrint("  No results found. Try a different workpiece.")
		return
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].pdiff < results[j].pdiff
	})

	best := results[0]
	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))
	webPrint(fmt.Sprintf("  Result: %0.9f", best.value))

	if radical == 2 {
		trueVal := math.Sqrt(float64(workpiece))
		webPrint(fmt.Sprintf("  math.Sqrt: %0.9f", trueVal))
		webPrint(fmt.Sprintf("  Difference: %0.9f", math.Abs(best.value-trueVal)))
	} else {
		trueVal := math.Cbrt(float64(workpiece))
		webPrint(fmt.Sprintf("  math.Cbrt: %0.9f", trueVal))
		webPrint(fmt.Sprintf("  Difference: %0.9f", math.Abs(best.value-trueVal)))
	}

	webPrint("")
	webPrint("  The perfect power pair:")
	webPrint(fmt.Sprintf("    %d / %d = %.4f  (target: %d)",
		best.largerPP, best.smallerPP,
		float64(best.largerPP)/float64(best.smallerPP),
		workpiece))
	webPrint(fmt.Sprintf("    Root ratio: %d / %d = %.9f",
		best.rootOfLarger, best.rootOfSmaller, best.value))
	webPrint("")
	webPrint(fmt.Sprintf("  Candidates found: %d", len(results)))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))

	if workpiece == 2 && radical == 3 {
		webPrint("")
		webPrint("  * — The Delian Problem — *")
		webPrint("  The cube root of 2 tormented geometers for 2,000 years.")
		webPrint("  Archimedes would have wept.")
	}

	webPrint(pkg.BoxSep(50))
}