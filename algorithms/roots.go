// Package algorithms implements various mathematical algorithms for
// approximating π, calculating roots, and exploring related constants.
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
// The search stops early when the gap between consecutive perfect powers
// exceeds the precision window, making it fast and efficient.
//
// Original concept and algorithm: Richard (Rick) Woolley.
func Roots(webPrint func(string), radical, workpiece int) {
	webPrint(pkg.BoxLine("  ROOTS VIA PERFECT POWERS  ", 50))
	webPrint(pkg.BoxLine("  Original Algorithm by Rick Woolley  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")

	if radical == 2 {
		webPrint(fmt.Sprintf("  Finding √%d", workpiece))
	} else {
		webPrint(fmt.Sprintf("  Finding ∛%d", workpiece))
	}
	webPrint("  Method: integer arithmetic only.")
	webPrint("  ─────────────────────────────────────────────────")
	webPrint("")

	start := time.Now()

	// Build table of perfect powers
	type Pair struct {
		product int
		root    int
	}
	var table []Pair

	if radical == 2 {
		webPrint("  Building table of perfect squares...")
	} else {
		webPrint("  Building table of perfect cubes...")
	}

	for root := 2; root < 825000; root++ {
		var product int
		if radical == 2 {
			product = root * root
		} else {
			product = root * root * root
		}
		table = append(table, Pair{product, root})
	}

	if radical == 2 {
		webPrint(fmt.Sprintf("  Table built: %s perfect squares.",
			pkg.FormatIntWithCommas(int64(len(table)))))
	} else {
		webPrint(fmt.Sprintf("  Table built: %s perfect cubes.",
			pkg.FormatIntWithCommas(int64(len(table)))))
	}

	// Show first few entries
	webPrint("")
	webPrint("  First few entries in the table:")
	for k := 0; k < 5 && k < len(table); k++ {
		if radical == 2 {
			webPrint(fmt.Sprintf("    root=%d  →  %d^2 = %d",
				table[k].root, table[k].root, table[k].product))
		} else {
			webPrint(fmt.Sprintf("    root=%d  →  %d^3 = %d",
				table[k].root, table[k].root, table[k].product))
		}
	}
	webPrint("    ...")
	webPrint("")

	webPrint("  The method: for each perfect power in the table,")
	webPrint("  search forward for another whose ratio brackets")
	webPrint(fmt.Sprintf("  our target number %d.", workpiece))
	webPrint("  When largerPP / smallerPP ≈ workpiece,")
	webPrint("  then rootOfLarger / rootOfSmaller ≈ answer.")
	webPrint("")

	// Show example with first two entries
	if len(table) >= 2 {
		a := table[0]
		b := table[1]
		webPrint("  Example with first two entries:")
		webPrint(fmt.Sprintf("    entry[0]: root=%d  product=%d", a.root, a.product))
		webPrint(fmt.Sprintf("    entry[1]: root=%d  product=%d", b.root, b.product))
		webPrint(fmt.Sprintf("    ratio = %d / %d = %.4f  (target: %d)",
			b.product, a.product,
			float64(b.product)/float64(a.product),
			workpiece))
		webPrint(fmt.Sprintf("    root ratio = %d / %d = %.6f",
			b.root, a.root,
			float64(b.root)/float64(a.root)))
		webPrint("    That ratio is our first rough approximation.")
		webPrint("    We keep searching for ratios ever closer to the target.")
	}
	webPrint("")

	// Calculate precision window
	precision := calculatePrecision(radical, workpiece)
	webPrint(fmt.Sprintf("  Precision window set to %d", precision))
	webPrint("")
	webPrint("  Starting search...")
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

	firstHitShown := false
	midpointShown := false
	midpoint := len(table) / 2
	threeQuarterMark := (len(table) * 3) / 4

	var trueVal float64
	if radical == 2 {
		trueVal = math.Sqrt(float64(workpiece))
	} else {
		trueVal = math.Cbrt(float64(workpiece))
	}

	// Helper function for the search
	runSearch := func() (perfectMatch bool) {
		results = nil
		firstHitShown = false
		midpointShown = false

		for i := 0; i < len(table)-1; i += 2 {
			smaller := table[i]
			target := smaller.product * workpiece

			prevCount := len(results)

			// Search forward for bracketing pair
			for j := i + 1; j < len(table); j++ {
				larger := table[j]

				if larger.product > target {
					dL := larger.product - target
					dS := target - table[j-1].product

					// Check for perfect match
					if dL == 0 || dS == 0 {
						webPrint("")
						webPrint("  ── Perfect Result ───────────────────────────────")
						if radical == 2 {
							webPrint(fmt.Sprintf("  %d is a perfect square.", workpiece))
							webPrint(fmt.Sprintf("  Its square root is exactly %0.0f", trueVal))
						} else {
							webPrint(fmt.Sprintf("  %d is a perfect cube.", workpiece))
							webPrint(fmt.Sprintf("  Its cube root is exactly %0.0f", trueVal))
						}
						webPrint(fmt.Sprintf("  Completed in: %s", time.Since(start).Round(time.Millisecond)))
						webPrint("  ─────────────────────────────────────────────────")
						return true
					}

					// Record candidates within precision window
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

					// Show first candidate
					if !firstHitShown && len(results) > prevCount {
						firstHitShown = true
						best := results[len(results)-1]
						webPrint(fmt.Sprintf("  First candidate found at table index %d:", i))
						webPrint(fmt.Sprintf("    smaller PP : %s  (root %d)",
							pkg.FormatIntWithCommas(int64(smaller.product)), smaller.root))
						webPrint(fmt.Sprintf("    larger PP  : product near %s",
							pkg.FormatIntWithCommas(int64(target))))
						webPrint(fmt.Sprintf("    root ratio = %0.9f", best.value))
						if radical == 2 {
							webPrint(fmt.Sprintf("    math.Sqrt  = %0.9f", trueVal))
						} else {
							webPrint(fmt.Sprintf("    math.Cbrt  = %0.9f", trueVal))
						}
						webPrint(fmt.Sprintf("    difference = %0.9f", math.Abs(best.value-trueVal)))
						webPrint("")
						webPrint("  Continuing search for a better approximation...")
						webPrint("")
					}

					break
				}
			}

			// Progress every 80,000 iterations
			if i%80000 == 0 && i > 0 {
				elapsed := time.Since(start)
				webPrint(fmt.Sprintf("  %s iterations completed...  elapsed: %s",
					pkg.FormatIntWithCommas(int64(i)),
					elapsed.Round(time.Millisecond)))
				if len(results) > 0 {
					sort.Slice(results, func(a, b int) bool {
						return results[a].pdiff < results[b].pdiff
					})
					best := results[0]
					webPrint(fmt.Sprintf("  Best so far: %0.9f  (candidates: %d)",
						best.value, len(results)))
				}
			}

			// Midpoint summary
			if !midpointShown && i >= midpoint {
				midpointShown = true
				mid := table[i]
				webPrint("")
				webPrint(fmt.Sprintf("  ── Midpoint check (index %s of %s) ──",
					pkg.FormatIntWithCommas(int64(i)),
					pkg.FormatIntWithCommas(int64(len(table)))))
				webPrint(fmt.Sprintf("  Now testing ratios involving root=%d", mid.root))
				webPrint(fmt.Sprintf("  whose perfect power is %s",
					pkg.FormatIntWithCommas(int64(mid.product))))
				if len(results) > 0 {
					sort.Slice(results, func(a, b int) bool {
						return results[a].pdiff < results[b].pdiff
					})
					best := results[0]
					webPrint(fmt.Sprintf("  Best so far: %0.9f", best.value))
					if radical == 2 {
						webPrint(fmt.Sprintf("  True value (math.Sqrt): %0.9f", trueVal))
					} else {
						webPrint(fmt.Sprintf("  True value (math.Cbrt): %0.9f", trueVal))
					}
				}
				webPrint("")
			}

			// Three-quarter mark
			if midpointShown && i == threeQuarterMark {
				tq := table[i]
				webPrint("")
				webPrint(fmt.Sprintf("  ── Three-quarter mark (index %s) ──",
					pkg.FormatIntWithCommas(int64(i))))
				webPrint(fmt.Sprintf("  Now testing root=%d (perfect power: %s)",
					tq.root, pkg.FormatIntWithCommas(int64(tq.product))))
				if len(results) > 0 {
					sort.Slice(results, func(a, b int) bool {
						return results[a].pdiff < results[b].pdiff
					})
					best := results[0]
					webPrint(fmt.Sprintf("  Best so far: %0.9f  candidates: %d",
						best.value, len(results)))
				}
				webPrint("")
			}

			// Early exit: gap exceeds precision window
			if i+2 < len(table) {
				gap := table[i+2].product - table[i].product
				if gap > precision && len(results) > 0 {
					webPrint("")
					webPrint(fmt.Sprintf("  Gap between consecutive perfect powers (%s)",
						pkg.FormatIntWithCommas(int64(gap))))
					webPrint(fmt.Sprintf("  now exceeds precision window (%d).", precision))
					webPrint("  No better result is mathematically possible.")
					webPrint("  Searching 2000 more iterations just to be sure...")
					for extra := i + 2; extra < i+2002 && extra < len(table)-1; extra += 2 {
						sm := table[extra]
						tg := sm.product * workpiece
						for j := extra + 1; j < len(table); j++ {
							lg := table[j]
							if lg.product > tg {
								dL := lg.product - tg
								if dL < precision {
									results = append(results, Result{
										value:         float64(lg.root) / float64(sm.root),
										pdiff:         float64(dL) / float64(lg.product),
										largerPP:      lg.product,
										smallerPP:     sm.product,
										rootOfLarger:  lg.root,
										rootOfSmaller: sm.root,
									})
								}
								break
							}
						}
					}
					webPrint("  Done. Stopping search.")
					return false
				}
			}
		}
		return false
	}

	// First attempt
	if runSearch() {
		return
	}

	// Auto-widen and retry if nothing found
	for len(results) == 0 {
		if precision >= 500000 {
			webPrint("")
			webPrint("  Could not find a result even at maximum precision window.")
			webPrint("  This workpiece may be too large for the current table size.")
			return
		}
		precision *= 2
		webPrint("")
		webPrint(fmt.Sprintf("  No results found. Widening precision window to %d and retrying...  (elapsed: %s)",
			precision, time.Since(start).Round(time.Millisecond)))
		webPrint("")
		if runSearch() {
			return
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
	webPrint("  ── Result ───────────────────────────────────────")
	if radical == 2 {
		webPrint(fmt.Sprintf("  Square Root of %d", workpiece))
		webPrint(fmt.Sprintf("  Our result    : %0.9f", best.value))
		webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Sqrt)", trueVal))
		webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best.value-trueVal)))

		frac := best.value - math.Floor(best.value)
		if frac < 0.01 || frac > 0.99 {
			webPrint("")
			webPrint("  Note: result is very close to a whole number.")
			webPrint(fmt.Sprintf("  %d sits just below %d² = %d,",
				workpiece,
				int(math.Round(best.value)),
				int(math.Round(best.value))*int(math.Round(best.value))))
			webPrint("  which is the best rational approximation available.")
			webPrint("  This is a known characteristic of the method for")
			webPrint("  numbers that sit just below a perfect square.")
		}
	} else {
		webPrint(fmt.Sprintf("  Cube Root of %d", workpiece))
		webPrint(fmt.Sprintf("  Our result    : %0.9f", best.value))
		webPrint(fmt.Sprintf("  Verification  : %0.9f  (math.Cbrt)", trueVal))
		webPrint(fmt.Sprintf("  Difference    : %0.9f", math.Abs(best.value-trueVal)))
	}

	webPrint("")
	webPrint("  ── The perfect power pair behind this result ─────")
	webPrint(fmt.Sprintf("  Anchor (smaller) PP : %s  (root %d)",
		pkg.FormatIntWithCommas(int64(best.smallerPP)), best.rootOfSmaller))
	webPrint(fmt.Sprintf("  Bracket PP          : %s  (root %d)",
		pkg.FormatIntWithCommas(int64(best.largerPP)), best.rootOfLarger))
	webPrint(fmt.Sprintf("  PP ratio            : %s / %s = %.6f  (target: %d)",
		pkg.FormatIntWithCommas(int64(best.largerPP)),
		pkg.FormatIntWithCommas(int64(best.smallerPP)),
		float64(best.largerPP)/float64(best.smallerPP),
		workpiece))
	webPrint("")

	highRatio := float64(best.rootOfLarger) / float64(best.rootOfSmaller)
	lowRatio := float64(best.rootOfLarger-1) / float64(best.rootOfSmaller)
	average := (highRatio + lowRatio) / 2.0

	webPrint(fmt.Sprintf("  Root ratio (high)   : %d / %d = %.9f",
		best.rootOfLarger, best.rootOfSmaller, highRatio))
	webPrint(fmt.Sprintf("  Root ratio (low)    : %d / %d = %.9f",
		best.rootOfLarger-1, best.rootOfSmaller, lowRatio))
	webPrint(fmt.Sprintf("  Average of high+low : %.9f", average))
	if radical == 2 {
		webPrint(fmt.Sprintf("  True value          : %.9f  (math.Sqrt)", trueVal))
	} else {
		webPrint(fmt.Sprintf("  True value          : %.9f  (math.Cbrt)", trueVal))
	}
	webPrint("")
	webPrint("  You can verify this yourself: divide the two perfect")
	if radical == 2 {
		webPrint("  squares above, then take the ratio of their roots.")
	} else {
		webPrint("  cubes above, then take the ratio of their roots.")
	}
	webPrint("  Averaging the high and low ratios gives an even")
	webPrint("  better approximation than either alone.")

	webPrint(fmt.Sprintf("  Candidates found: %d (best pdiff: %0.8f)",
		len(results), best.pdiff))
	webPrint(fmt.Sprintf("  Completed in    : %s", elapsed.Round(time.Millisecond)))
	webPrint("  ─────────────────────────────────────────────────")
	webPrint("")
	webPrint("  (Verification values from math.Sqrt/Cbrt are shown")
	webPrint("   only to confirm accuracy. The calculation above")
	webPrint("   used integer arithmetic exclusively.)")

	if workpiece == 2 && radical == 3 {
		webPrint("")
		webPrint("  * — THE DELIAN PROBLEM — *")
		webPrint("  The cube root of 2 tormented ancient Greek")
		webPrint("  geometers for over two thousand years.")
		webPrint("  ")
		webPrint("  ✦ We just solved it in milliseconds ✦")
		webPrint("  ")
		webPrint("  ✦✦ Archimedes would have wept ✦✦")
		webPrint(" ")
		webPrint("  The cube root of two was finally proved impossible")
		webPrint("  with compass and straightedge in 1837 by")
		webPrint("  Pierre Wantzel — over 2,200 years after the")
		webPrint("  Greeks first posed it.")
		webPrint(" ")
	}
}

// calculatePrecision determines the precision window based on radical and workpiece.
func calculatePrecision(radical, workpiece int) int {
	var precision int
	if radical == 2 {
		precision = int(math.Sqrt(float64(workpiece))) / 3
		precision = max(precision, 4)
		precision = min(precision, 500)
	} else {
		precision = workpiece * 12
		precision = max(precision, 600)
		precision = min(precision, 3000)
	}
	return precision
}