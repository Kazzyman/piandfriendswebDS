package algorithms

import (
	"fmt"
	"math"
	"math/big"
	"time"

	"piandfriends/pkg"
)

// Chudnovsky computes π using the Chudnovsky algorithm.
// This is the algorithm used for world-record π calculations.
// Convergence is extremely rapid: ~14 digits per term.
func Chudnovsky(done chan bool, webPrint func(string), digits int) {
	webPrint(pkg.BoxLine("  CHUDNOVSKY ALGORITHM  ", 50))
	webPrint(pkg.BoxLine("  World-record π calculations  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint(fmt.Sprintf("  Computing %d digits...", digits))
	webPrint("")

	if digits > 10000 {
		webPrint("  That's a lot of pie. This may take a while...")
	}

	start := time.Now()

	// Terms needed: each term gives ~14.18 decimal digits
	terms := int(float64(digits)/14.181647462) + 2

	// Precision in bits
	prec := uint(float64(digits)*math.Log2(10) + float64(digits)*0.1 + 64)

	// Constants
	C := new(big.Float).SetPrec(prec)
	C.Mul(big.NewFloat(426880), new(big.Float).SetPrec(prec).Sqrt(big.NewFloat(10005)))

	K := big.NewInt(6)
	K12 := big.NewInt(12)
	L := new(big.Float).SetPrec(prec).SetFloat64(13591409)
	LC := new(big.Float).SetPrec(prec).SetFloat64(545140134)
	X := new(big.Float).SetPrec(prec).SetFloat64(1)
	XC := new(big.Float).SetPrec(prec).SetFloat64(-262537412640768000)
	M := new(big.Float).SetPrec(prec).SetFloat64(1)
	Sum := new(big.Float).SetPrec(prec).SetFloat64(13591409)

	bigOne := big.NewInt(1)
	bigI := big.NewInt(0)

	for i := 0; i < terms; i++ {
		select {
		case <-done:
			webPrint("  Stopped.")
			return
		default:
		}

		// L calculation
		L.Add(L, LC)

		// X calculation
		X.Mul(X, XC)

		// M calculation
		kpower3 := new(big.Int).Exp(K, big.NewInt(3), nil)
		ktimes16 := new(big.Int).Mul(K, big.NewInt(16))
		mtop := new(big.Int).Sub(kpower3, ktimes16)

		iPlusOne := new(big.Int).Add(bigI, bigOne)
		mbot := new(big.Int).Exp(iPlusOne, big.NewInt(3), nil)

		mtopF := new(big.Float).SetPrec(prec).SetInt(mtop)
		mbotF := new(big.Float).SetPrec(prec).SetInt(mbot)
		mtmp := new(big.Float).SetPrec(prec).Quo(mtopF, mbotF)
		M.Mul(M, mtmp)

		// Sum calculation
		t := new(big.Float).SetPrec(prec)
		t.Mul(M, L)
		t.Quo(t, X)
		Sum.Add(Sum, t)

		K.Add(K, K12)
		bigI.Add(bigI, bigOne)

		// Progress indicator
		if i > 0 && i%1000 == 0 {
			webPrint(fmt.Sprintf("  ... term %d", i))
		}
	}

	// π = C / Sum
	pi := new(big.Float).SetPrec(prec).Quo(C, Sum)

	elapsed := time.Since(start)

	webPrint("")
	webPrint(pkg.BoxSep(50))

	showDigits := digits
	if showDigits > 100 {
		showDigits = 100
	}
	piStr := pi.Text('f', showDigits+2)

	webPrint(fmt.Sprintf("  π = %s", piStr))
	if digits > 100 {
		webPrint(fmt.Sprintf("  (... and %d more digits)", digits-100))
	}

	// Cross-verify
	verifyMsg := pkg.VerifyAndReport(pi, digits, "Chudnovsky")
	webPrint(verifyMsg)

	webPrint(fmt.Sprintf("  Terms computed: %d", terms))
	webPrint(fmt.Sprintf("  Time: %s", elapsed.Round(time.Millisecond)))
	webPrint("")
	webPrint("  The Chudnovsky brothers developed this algorithm")
	webPrint("  in the late 1980s. It has been used to compute")
	webPrint("  over 300 trillion digits of π.")
	webPrint(pkg.BoxSep(50))
}