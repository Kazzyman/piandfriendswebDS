package algorithms

import (
	"fmt"
	"math/big"
	"time"
	"piandfriends/pkg"
)

// NilakanthaClassic runs a two-phase dramatic demonstration:
// Phase 1: float64 with concurrent goroutines, slamming into the 15-digit wall.
// Phase 2: big.Float breaking through to 20+ digits.
//
// This is a fixed demonstration. No user inputs—just watch the tragedy and triumph.
func NilakanthaClassic(done chan bool, webPrint func(string)) {
	const bw = 50

	// Fixed parameters for maximum drama
	const phase1Terms = 500000   // Enough to slam the wall repeatedly
	const phase2Terms = 10000000 // 10 million: breakthrough to ~21 digits

	webPrint(pkg.BoxLine("  NILAKANTHA: THE WALL  ", bw))
	webPrint(pkg.BoxLine("  A Tragedy in Two Acts  ", bw))
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint("  Act I: float64")
	webPrint("  Act II: big.Float")
	webPrint("")
	webPrint("  Nilakantha Somayaji's series (c. 1530):")
	webPrint("  π = 3 + 4/(2·3·4) - 4/(4·5·6) + 4/(6·7·8) - ...")
	webPrint("")
	webPrint(pkg.BoxSep(bw))
	webPrint("")

	// ──────────────────────────────────────────────────────────────────
	// ACT I: THE FALL OF FLOAT64
	// ──────────────────────────────────────────────────────────────────

	webPrint("COLOR:cyan:" + pkg.BoxLine("  ACT I: float64  ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  The March Toward the Wall  ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxSep(bw))
	webPrint("")
	webPrint(fmt.Sprintf("  Computing %s terms with float64...",
		pkg.FormatIntWithCommas(int64(phase1Terms))))
	webPrint("")
	webPrint("  Each milestone locks in another decimal digit.")
	webPrint("  But there is a ceiling. Watch closely.")
	webPrint("")

	phase1Start := time.Now()

	ch := make(chan float64, phase1Terms)

	// k=1 term: 4/(2*3*4) = 1/6
	sum := 3.0 + 1.0/6.0

	// Launch goroutines for k=2 to phase1Terms
	for k := 2; k <= phase1Terms; k++ {
		go func(kk int) {
			d1 := float64(2 * kk)
			d2 := float64(2*kk + 1)
			d3 := float64(2*kk + 2)
			term := 4.0 / (d1 * d2 * d3)
			if kk%2 == 0 {
				term = -term
			}
			ch <- term
		}(k)
	}

	bestDigits := 0
	wallAnnounced := false
	lastMilestone := 0

	milestoneMessages := map[int]string{
		1:  "3.1 — the journey begins",
		2:  "3.14 — the famous digits appear",
		3:  "3.141 — three digits secured",
		4:  "3.1415 — four digits",
		5:  "3.14159 — five correct decimal digits",
		6:  "3.141592 — six digits",
		7:  "3.1415926 — seven digits. Confidence grows.",
		8:  "3.14159265 — eight digits",
		9:  "3.141592653 — nine digits",
		10: "3.1415926535 — ten digits. float64 is strong.",
		11: "3.14159265358 — eleven digits",
		12: "3.141592653589 — twelve digits",
		13: "3.1415926535897 — thirteen digits",
		14: "3.14159265358979 — fourteen. The wall nears.",
		15: "3.141592653589793 — FIFTEEN DIGITS.",
	}

	for k := 2; k <= phase1Terms; k++ {
		select {
		case <-done:
			webPrint("  The curtain falls early.")
			return
		case term := <-ch:
			sum += term

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

			// Milestone fanfare
			if correct > bestDigits {
				bestDigits = correct
				if msg, ok := milestoneMessages[correct]; ok {
					webPrint(fmt.Sprintf("  ✦ %s", msg))
					lastMilestone = k
				}
			}

			// Wall detection
			if !wallAnnounced && bestDigits >= 15 && k > lastMilestone+10000 {
				wallAnnounced = true
				webPrint("")
				webPrint("COLOR:red:" + pkg.BoxLine("  !! THE WALL !!  ", bw))
				webPrint("COLOR:red:  float64 can go no further.")
				webPrint("COLOR:red:  The remaining digits are noise.")
				webPrint("COLOR:red:  More iterations change nothing.")
				webPrint("")
			}

			// Progress updates
			if k%50000 == 0 || k == phase1Terms {
			pct := float64(k) / float64(phase1Terms) * 100
				piDisplay := fmt.Sprintf("%.15f", sum)
				if wallAnnounced {
					// Show the stagnation dramatically
					piDisplay = piDisplay + " (unchanging)"
				}
				webPrint(fmt.Sprintf("  [%3.0f%%] %s terms: π ≈ %s",
					pct,
					pkg.FormatIntWithCommas(int64(k)),
					piDisplay))
			}
		}
	}

	phase1Elapsed := time.Since(phase1Start)

	webPrint("")
	webPrint(pkg.BoxSep(bw))
	webPrint("  ACT I COMPLETE")
	webPrint(fmt.Sprintf("  Time: %s", phase1Elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Final: %.15f", sum))
	webPrint("  Correct digits: 15 — AND NO MORE.")
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint("  float64 has given everything it has.")
	webPrint("  The 53-bit mantissa cannot represent")
	webPrint("  more than ~15-16 decimal digits precisely.")
	webPrint("  This is a hardware limit, not a mathematical one.")
	webPrint("")
	webPrint("  But wait. The story is not over.")
	webPrint("")

	time.Sleep(3 * time.Second)

	// ──────────────────────────────────────────────────────────────────
	// ACT II: THE REDEMPTION OF BIG.FLOAT
	// ──────────────────────────────────────────────────────────────────

	webPrint("COLOR:green:" + pkg.BoxLine("  ACT II: big.Float  ", bw))
	webPrint("COLOR:green:" + pkg.BoxLine("  Breaking the Wall  ", bw))
	webPrint("COLOR:green:" + pkg.BoxSep(bw))
	webPrint("")
	webPrint(fmt.Sprintf("  Computing %s terms with big.Float...",
		pkg.FormatIntWithCommas(int64(phase2Terms))))
	webPrint("")
	webPrint("  big.Float uses arbitrary precision.")
	webPrint("  The wall does not exist here.")
	webPrint("")

	phase2Start := time.Now()
	prec := uint(512)

	three := new(big.Float).SetPrec(prec).SetFloat64(3.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)

	firstTerm := new(big.Float).SetPrec(prec).Quo(four,
		new(big.Float).SetPrec(prec).SetFloat64(2*3*4))
	piB := new(big.Float).SetPrec(prec).Add(three, firstTerm)

	wallBroken := false
	bestBigDigits := 0

	for k := 2; k <= phase2Terms; k++ {
		select {
		case <-done:
			webPrint("  The curtain falls early.")
			return
		default:
		}

		d1 := new(big.Float).SetPrec(prec).SetFloat64(float64(2 * k))
		d2 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 1))
		d3 := new(big.Float).SetPrec(prec).SetFloat64(float64(2*k + 2))

		denom := new(big.Float).SetPrec(prec).Mul(d1,
			new(big.Float).SetPrec(prec).Mul(d2, d3))
		term := new(big.Float).SetPrec(prec).Quo(four, denom)

		if k%2 == 0 {
			piB.Sub(piB, term)
		} else {
			piB.Add(piB, term)
		}

		if k%500000 == 0 || k == phase2Terms {
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

			if !wallBroken && correct > 15 {
				wallBroken = true
				webPrint("")
				webPrint("COLOR:yellow:  ╔══════════════════════════════════════════╗")
				webPrint("COLOR:yellow:  ║  !! THE WALL IS BROKEN !!              ║")
				webPrint("COLOR:yellow:  ║  big.Float sees what float64 could not ║")
				webPrint("COLOR:yellow:  ╚══════════════════════════════════════╝")
				webPrint("")
			}

			if correct > bestBigDigits {
				bestBigDigits = correct
				if correct > 15 {
					webPrint(fmt.Sprintf("  ✦ %d digits — BEYOND THE WALL", correct))
				}
			}

			pct := float64(k) / float64(phase2Terms) * 100
			piStr := piB.Text('f', correct+2)
			if len(piStr) > 25 {
				piStr = piStr[:25] + "..."
			}

			webPrint(fmt.Sprintf("  [%3.0f%%] %s terms: π ≈ %s (%d digits)",
				pct,
				pkg.FormatIntWithCommas(int64(k)),
				piStr,
				correct))
		}
	}

	phase2Elapsed := time.Since(phase2Start)

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
	webPrint("  ACT II COMPLETE")
	webPrint(fmt.Sprintf("  Time: %s", phase2Elapsed.Round(time.Millisecond)))
	webPrint(fmt.Sprintf("  Total terms: %s",
		pkg.FormatIntWithCommas(int64(phase2Terms))))
	webPrint(fmt.Sprintf("  Verified digits: %d", correct))
	webPrint(fmt.Sprintf("  π = %s", piB.Text('f', correct+1)))
	webPrint(pkg.BoxSep(bw))
	webPrint("")
	webPrint("  The wall was never a mathematical limit.")
	webPrint("  It was a hardware limit. A representation limit.")
	webPrint("")
	webPrint("  Nilakantha's series, born in Kerala c. 1530,")
	webPrint("  continues to give. It only needed the right")
	webPrint("  tools to be heard.")
	webPrint("")
	webPrint("  Curtain call.")
	webPrint(pkg.BoxSep(bw))
}