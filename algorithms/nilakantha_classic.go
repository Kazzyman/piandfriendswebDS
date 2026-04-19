package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// NilakanthaClassic runs a two-phase demonstration:
// Phase 1: float64 with concurrent goroutines, showing the float64 wall.
// Phase 2: big.Float breaking through that wall.
func NilakanthaClassic(done chan bool, webPrint func(string), n1, n2 int) {
	const bw = 50

	webPrint(pkg.BoxLine("  NILAKANTHA TWO-PHASE DEMO  ", bw))
	webPrint(pkg.BoxLine("  π = 3 + 4/(2·3·4) - 4/(4·5·6) + ...  ", bw))
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint(fmt.Sprintf("  Phase 1: %s terms (float64, concurrent)",
		pkg.FormatIntWithCommas(int64(n1))))
	webPrint("  Ceiling: ~15 digits")
	webPrint(fmt.Sprintf("  Phase 2: %s additional terms (big.Float)",
		pkg.FormatIntWithCommas(int64(n2))))
	webPrint("  Ceiling: arbitrary precision")
	webPrint("")

	// ── Phase 1: float64 with goroutines ───────────────────────────────
	webPrint("COLOR:yellow:  ── PHASE 1 BEGIN ──")
	webPrint("")

	phase1Start := time.Now()

	ch := make(chan float64, n1)
	sum := 3.0

	// Launch goroutines for each term
	for k := 1; k <= n1; k++ {
		go func(kk int) {
			j := float64(2 * kk)
			term := 4.0 / (j * (j + 1) * (j + 2))
			if kk%2 == 0 {
				term = -term
			}
			ch <- term
		}(k)
	}

	// Collect results
	bestDigits := 0
	for k := 1; k <= n1; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		case term := <-ch:
			sum += term

			if k%1000 == 0 || k == n1 {
				// Count correct digits
				correct := 0
				ref := "3.14159265358979323846"
				s := fmt.Sprintf("%.20f", sum)
				for i := 0; i < len(ref) && i < len(s); i++ {
					if s[i] != ref[i] {
						break
					}
					correct++
				}
				if correct > 2 {
					correct -= 2
				} else {
					correct = 0
				}

				if correct > bestDigits {
					bestDigits = correct
					webPrint(fmt.Sprintf("  *** %d correct digits!", correct))
				}

				pct := float64(k) / float64(n1) * 100
				elapsed := time.Since(phase1Start)
				webPrint(fmt.Sprintf("UPDATE:  π ≈ %.15f  [%.0f%%]  terms: %s  %s",
					sum, pct,
					pkg.FormatIntWithCommas(int64(k)),
					elapsed.Round(time.Millisecond)))
			}

			// Pace for visibility
			if k < 100 {
				time.Sleep(10 * time.Millisecond)
			}
		}
	}

	phase1Elapsed := time.Since(phase1Start)
	phase1Final := sum

	webPrint("")
	webPrint(pkg.BoxSep(bw))
	webPrint(fmt.Sprintf("  PHASE 1 COMPLETE: %s", phase1Elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Final: %.15f", phase1Final))
	webPrint(fmt.Sprintf("  Correct digits: %d", bestDigits))
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint("  float64 has given everything it has.")
	webPrint("  The wall at ~15 digits is real.")
	webPrint("")
	webPrint("  Phase 2 will now run in big.Float and break through.")
	webPrint("")

	time.Sleep(2 * time.Second)

	// ── Phase 2: big.Float ─────────────────────────────────────────────
	webPrint("COLOR:cyan:  ── PHASE 2 BEGIN ──")
	webPrint("")

	phase2Start := time.Now()
	prec := uint(512)

	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	d1 := new(big.Float).SetPrec(prec).SetFloat64(2.0)
	d2 := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	d3 := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	// Recompute from scratch in big.Float
	firstTerm := new(big.Float).SetPrec(prec).Quo(four,
		new(big.Float).Mul(d1, new(big.Float).Mul(d2, d3)))
	piB := new(big.Float).SetPrec(prec).Add(three, firstTerm)

	wallBroken := false
	totalTerms := n1 + n2

	for k := 2; k <= totalTerms; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		d1.Add(d1, two)
		d2.Add(d2, two)
		d3.Add(d3, two)

		term := new(big.Float).SetPrec(prec).Quo(four,
			new(big.Float).Mul(d1, new(big.Float).Mul(d2, d3)))

		if k%2 == 0 {
			piB.Sub(piB, term)
		} else {
			piB.Add(piB, term)
		}

		if k%10000 == 0 || k == totalTerms {
			// Count correct digits using CrossVerify (indirect count)
			showDigits := 22
			// Simple character-by-character check against known π start
			correct := 0
			ref := "3.14159265358979323846"
			s := piB.Text('f', showDigits)
			for i := 0; i < len(ref) && i < len(s); i++ {
				if s[i] != ref[i] {
					break
				}
				correct++
			}
			if correct > 2 {
				correct -= 2
			} else {
				correct = 0
			}

			if !wallBroken && correct > 15 {
				wallBroken = true
				webPrint("")
				webPrint("COLOR:cyan:  ╔══════════════════════════════════════════╗")
				webPrint("COLOR:cyan:  ║  !! THE FLOAT64 WALL HAS BEEN BROKEN !!  ║")
				webPrint("COLOR:cyan:  ╚══════════════════════════════════════════╝")
				webPrint("")
			}

			pct := float64(k-n1) / float64(n2) * 100
			if pct < 0 {
				pct = 0
			}
			elapsed := time.Since(phase2Start)
			piStr := piB.Text('f', 20)

			webPrint(fmt.Sprintf("UPDATE:  π ≈ %s  [%.0f%%]  term: %s  %s  (%d digits)",
				piStr, pct,
				pkg.FormatIntWithCommas(int64(k)),
				elapsed.Round(time.Millisecond),
				correct))
		}
	}

	phase2Elapsed := time.Since(phase2Start)

	// Final digit count
	correct := 0
	ref := "3.141592653589793238462643383279"
	s := piB.Text('f', 30)
	for i := 0; i < len(ref) && i < len(s); i++ {
		if s[i] != ref[i] {
			break
		}
		correct++
	}
	if correct > 2 {
		correct -= 2
	} else {
		correct = 0
	}

	webPrint("")
	webPrint(pkg.BoxSep(bw))
	webPrint(fmt.Sprintf("  PHASE 2 COMPLETE: %s", phase2Elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Total terms: %s", pkg.FormatIntWithCommas(int64(totalTerms))))
	webPrint(fmt.Sprintf("  Verified digits: %d", correct))
	webPrint(fmt.Sprintf("  π = %s", piB.Text('f', correct+1)))
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint("  Kerala school, c. 1530. Predates Newton by 150 years.")
	webPrint("  Still climbing.")
}