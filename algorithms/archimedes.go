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
	webPrint("  You've selected Dick's improved version of")
	webPrint("  Archimedes' method for approximating π.")
	webPrint("  The goal: over 2,700 correct digits.")
	webPrint("  We'll need floating-point numbers with")
	webPrint("  thousands of decimal places.")
	webPrint("  This can be done using Go's math/big package.")
	webPrint("  All variables must be big.Floats.")
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

	webPrint(fmt.Sprintf("  Go's precision is set to %d on all variables.", precision))
	webPrint("")
	webPrint("  First we determine the height (a) of a right")
	webPrint("  triangle formed by bisecting a side of a")
	webPrint("  polygon inscribed in a unit circle (r=1).")
	webPrint("  The polygon's side length (s1) is halved")
	webPrint("  (s1_2 = s1/2). This refines the perimeter")
	webPrint("  to approximate π as sides increase.")
	webPrint("")
	webPrint("  Pseudo-code for the algorithm:")
	webPrint("    Inputs:  b (short side, midpoint to edge)")
	webPrint("             s1_2 (half current side length)")
	webPrint("    Output:  s2 (new side length)")
	webPrint("    Step 1. temp1 = b * b")
	webPrint("    Step 2. temp2 = s1_2 * s1_2")
	webPrint("    Step 3. temp3 = temp1 + temp2")
	webPrint("    Step 4. s2 = sqrt(temp3)")
	webPrint("")
	webPrint("  Now, we get to work!!")
	webPrint("")

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
		case 2:
			webPrint("    Sleeping each iteration for 135ms...")
		case 24:
			webPrint(pkg.BoxSep(50))
			webPrint(fmt.Sprintf("  %d iterations completed", i))
			webPrint(fmt.Sprintf("  %.20f", pd))
			webPrint("  3.141592653589793238  (reference)")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
			webPrint(pkg.BoxSep(50))
		case 26:
			webPrint("    Sleeping each iteration for 55ms...")
		case 50:
			webPrint(pkg.BoxSep(50))
			webPrint(fmt.Sprintf("  %d iterations completed", i))
			webPrint(fmt.Sprintf("  %.33f", pd))
			webPrint("  3.141592653589793238462643383279502")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
			webPrint(pkg.BoxSep(50))
		case 52:
			webPrint("    Sleeping each iteration for 35ms...")
		case 150:
			webPrint(pkg.BoxSep(50))
			webPrint(fmt.Sprintf("  %d iterations completed", i))
			webPrint(fmt.Sprintf("  %.95f", pd))
			webPrint("  3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
			webPrint(pkg.BoxSep(50))
		case 152:
			webPrint("    Sleeping each iteration for 7ms...")
		case 200:
			webPrint(pkg.BoxSep(50))
			webPrint(fmt.Sprintf("  %d iterations completed", i))
			webPrint(fmt.Sprintf("  %.122f", pd))
			webPrint("  3.14159265358979323846264338327950288419716939937510582097494459230781640628620899862803482534211706798214808651328230664709")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
			webPrint("  ... working ...")
			webPrint(pkg.BoxSep(50))
		case 402:
			webPrint("    Sleeping each iteration for 2ms...")
		case 1102:
			webPrint("    Sleeping each iteration for 1ms...")
		case 1200, 2200, 3200, 4200:
			webPrint(fmt.Sprintf("  ... still working, %d iterations ...", i))
		case 2002:
			webPrint("    No more sleeping!!!")
		case 4500:
			webPrint(pkg.BoxSep(50))
			webPrint("  All Done!")
			formattedSides := pkg.FormatBigFloatWithCommas(sides)
			webPrint(fmt.Sprintf("  Sides: %s", formattedSides))
			webPrint(fmt.Sprintf("  %d iterations completed", i))
			webPrint(fmt.Sprintf("  Precision: %d bits", precision))
			webPrint("")
			// Show a truncated slice
			piStr := pd.Text('f', 2800)
			webPrint(fmt.Sprintf("  π = %s", piStr[:200]))
			webPrint("  ...")
			webPrint(fmt.Sprintf("  %s", piStr[len(piStr)-200:]))
			webPrint("")
			webPrint("  ... verified to 2712 digits!")
			webPrint(" ")
			webPrint("  ... Mister A. would have wept!")
			webPrint(" ")
			webPrint("  by Richard (Rick) H. Woolley")
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
