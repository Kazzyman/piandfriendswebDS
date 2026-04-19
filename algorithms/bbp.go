package algorithms

import (
	"fmt"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// BBP computes π using the Bailey–Borwein–Plouffe formula.
// This implementation evaluates the series sequentially in base-10.
// The formula itself can extract hexadecimal digits individually,
// but this demo computes sequentially for simplicity.
func BBP(done chan bool, webPrint func(string), digits int) {
	webPrint(pkg.BoxLine("  BAILEY–BORWEIN–PLOUFFE (BBP)  ", 50))
	webPrint(pkg.BoxLine("  Discovered 1995 via PSLQ algorithm  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint(fmt.Sprintf("  Computing %d digits...", digits))
	webPrint("  (Sequential base-10 evaluation)")
	webPrint("")

	start := time.Now()

	// Precision: ~3.32 bits per decimal digit, plus margin
	prec := uint(float64(digits)*3.33 + 64)

	pi := new(big.Float).SetPrec(prec).SetFloat64(0.0)

	one := new(big.Float).SetPrec(prec).SetFloat64(1.0)
	four := new(big.Float).SetPrec(prec).SetFloat64(4.0)
	two := new(big.Float).SetPrec(prec).SetFloat64(2.0)

	// Terms needed: roughly digits * 0.83
	terms := int(float64(digits) * 0.83)
	if terms < 10 {
		terms = 10
	}

	for k := 0; k < terms; k++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		// 16^k
		pow16 := new(big.Int).Exp(big.NewInt(16), big.NewInt(int64(k)), nil)
		pow16f := new(big.Float).SetPrec(prec).SetInt(pow16)

		k8 := float64(k * 8)

		t1 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 1)
		t2 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 4)
		t3 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 5)
		t4 := new(big.Float).SetPrec(prec).SetFloat64(k8 + 6)

		term1 := new(big.Float).SetPrec(prec).Quo(four, t1)
		term2 := new(big.Float).SetPrec(prec).Quo(two, t2)
		term3 := new(big.Float).SetPrec(prec).Quo(one, t3)
		term4 := new(big.Float).SetPrec(prec).Quo(one, t4)

		inner := new(big.Float).SetPrec(prec)
		inner.Sub(term1, term2)
		inner.Sub(inner, term3)
		inner.Sub(inner, term4)

		inner.Quo(inner, pow16f)

		pi.Add(pi, inner)

		// Progress indicator
		if k > 0 && k%1000 == 0 {
			webPrint(fmt.Sprintf("  ... term %d", k))
		}
	}

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	// Show result
	showDigits := digits
	if showDigits > 100 {
		showDigits = 100
	}
	piStr := pi.Text('f', showDigits+2)

	webPrint(fmt.Sprintf("  π = %s", piStr))
	if digits > 100 {
		webPrint(fmt.Sprintf("  (... and %d more correct digits)", digits-100))
	}
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  Note: BBP can extract individual hex digits of π")
	webPrint("  without computing preceding digits—a unique property.")
	webPrint(pkg.BoxSep(50))
}