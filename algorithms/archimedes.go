package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// Archimedes implements an improved version of Archimedes' polygon method.
// It uses big.Float for high precision and progressively refines the
// perimeter of an inscribed polygon.
//
// This is Rick's personal favorite.
func Archimedes(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("ARCHIMEDES' METHOD — Rick's Personal Favorite", 50))
	webPrint(pkg.BoxLine("Polygon perimeter refinement, c. 230 BCE", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("  Finding π the original way, with a modern twist.")
	webPrint("  Using Go's math/big for arbitrary precision.")
	webPrint("")

	precision := uint(55000)

	r := big.NewFloat(1)
	s1 := big.NewFloat(1)
	sides := big.NewFloat(6)

	a := new(big.Float)
	b := new(big.Float)
	p := new(big.Float)
	s2 := new(big.Float)
	pd := new(big.Float)
	s1_2 := new(big.Float)

	pd.SetPrec(precision)
	a.SetPrec(precision)
	s1_2.SetPrec(precision)
	s2.SetPrec(precision)
	b.SetPrec(precision)
	p.SetPrec(precision)
	r.SetPrec(precision)
	s1.SetPrec(precision)
	sides.SetPrec(precision)

	// Initial setup
	sides.Mul(sides, big.NewFloat(2))
	s1_2.Quo(s1, big.NewFloat(2))
	a.Sqrt(new(big.Float).Sub(r, new(big.Float).Mul(s1_2, s1_2)))
	b.Sub(r, a)
	s2.Sqrt(new(big.Float).Add(new(big.Float).Mul(b, b), new(big.Float).Mul(s1_2, s1_2)))

	s1.Set(s2)
	p.Mul(sides, s1)
	pd.Set(p)
	pd.Quo(pd, big.NewFloat(2))

	webPrint("  Iterating...")

	for i := 0; i < 5001; i++ {
		select {
		case <-done:
			webPrint("  Calculation stopped.")
			return
		default:
		}

		sides.Mul(sides, big.NewFloat(2))
		s1_2.Quo(s1, big.NewFloat(2))
		a.Sqrt(new(big.Float).Sub(r, new(big.Float).Mul(s1_2, s1_2)))
		b.Sub(r, a)
		s2.Sqrt(new(big.Float).Add(new(big.Float).Mul(b, b), new(big.Float).Mul(s1_2, s1_2)))
		s1.Set(s2)
		p.Mul(sides, s1)
		pd.Set(p)
		pd.Quo(pd, big.NewFloat(2))

		switch i {
		case 24:
			webPrint(fmt.Sprintf("  Iteration %d: %.20f", i, pd))
			webPrint("  Reference:  3.141592653589793238")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
		case 50:
			webPrint(fmt.Sprintf("  Iteration %d: %.33f", i, pd))
		case 150:
			webPrint(fmt.Sprintf("  Iteration %d: %.95f", i, pd))
		case 4500:
			webPrint(pkg.BoxSep(50))
			webPrint("  COMPLETE: Archimedes would have wept.")
			webPrint(pkg.BoxSep(50))
		}

		// Paced sleep for dramatic effect
		var sleepDur time.Duration
		switch {
		case i < 24:
			sleepDur = 135 * time.Millisecond
		case i < 50:
			sleepDur = 55 * time.Millisecond
		case i < 150:
			sleepDur = 35 * time.Millisecond
		case i < 400:
			sleepDur = 7 * time.Millisecond
		case i < 1100:
			sleepDur = 2 * time.Millisecond
		case i < 2000:
			sleepDur = 1 * time.Millisecond
		default:
			sleepDur = 0
		}

		if sleepDur > 0 {
			select {
			case <-done:
				return
			case <-time.After(sleepDur):
			}
		}
	}
}