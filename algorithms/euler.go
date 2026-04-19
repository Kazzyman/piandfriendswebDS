package algorithms

import (
	"fmt"
	"math"
	"time"

	"piandfriends/pkg"
)

// EulersNumber demonstrates the convergence of (1 + 1/n)^n to e.
// Includes the compound interest narrative from the original.
func EulersNumber(done chan bool, webPrint func(string)) {
	webPrint(pkg.BoxLine("  EULER'S NUMBER  ", 50))
	webPrint(pkg.BoxSep(50))
	webPrint("")
	webPrint("LARGE:Euler's Number ℯ")
	webPrint("LARGE:ℯ = (1 + 1/n)^n")
	webPrint("LARGE:... as n → ∞")
	webPrint("")

	ns := []float64{9, 99, 999, 9999, 99999999999}

	for _, n := range ns {
		select {
		case <-done:
			return
		default:
		}

		sum := 1.0 + 1.0/n
		e := math.Pow(sum, n)

		webPrint(fmt.Sprintf("LARGE:%0.45f", e))
		webPrint(fmt.Sprintf("LARGE:  calculated with n = %0.f", n))
		webPrint("LARGE: ")
		time.Sleep(500 * time.Millisecond)
	}

	webPrint("LARGE:")
	webPrint("LARGE:2.718281828459045... is Euler's Number")
	webPrint("LARGE:")
	webPrint("LARGE:An account starts with $1.00 at 100% annual interest.")
	webPrint("LARGE:Compounded once: $2.00")
	webPrint("LARGE:Compounded daily: $2.714567...")
	webPrint("LARGE:Compounded continuously: ℯ dollars")
	webPrint("LARGE:")
	webPrint("LARGE:Bernoulli noticed this limit in the 17th century.")
	webPrint("LARGE:It became known as e, Euler's Number.")
	webPrint("")
	webPrint("LARGE: ")
	webPrint("LARGE: ")
	webPrint("LARGE: ")



	webPrint("LARGE:   Again:")
		webPrint("LARGE: ")

	webPrint("LARGE:2.71828182845904523536028747135266249775724")
	webPrint("LARGE:That, is Euler's Number from the web")
		webPrint("LARGE: ")

	webPrint("LARGE:2.718281828 is the dollar value of $1 compounded")
	webPrint("LARGE:continuously for one year.")
		webPrint("LARGE: ")

	webPrint("LARGE:2.714567 is from daily compound interest which is")
	webPrint("LARGE:near-enough to continuous interest.")
			webPrint("LARGE: ")

	webPrint("LARGE:An account starts with $1.00 and pays 100 percent")
	webPrint("LARGE: interest per year. If the interest is credited once,")
			webPrint("LARGE: ")

	webPrint("LARGE:at the end of the year, the value of the account at year-")
	webPrint("LARGE:end will be $2.00. What happens if the interest is")
	webPrint("LARGE:computed and credited more frequently during the year?")
	webPrint("LARGE: ")

	webPrint("LARGE:If the interest is credited twice in the year, the interest")
	webPrint("LARGE:rate for each 6 months will be 50%, so the initial $1 is")
	webPrint("LARGE:multiplied by 1.5 twice, yielding $2.25 at the end of the")
	webPrint("LARGE:year. Compounding quarterly yields $2.44140625, and")
	webPrint("LARGE:compounding monthly yields $2.613035 = ")
	webPrint("LARGE:  $1.00 × (1 + 1/12)^12  Generally, if there are n")
	webPrint("LARGE:compounding intervals, the interest for each interval will")
	webPrint("LARGE:be 100%/n and the value at the end of the year will be")
	webPrint("LARGE:  $1.00 × (1 + 1/n)^n.")
	webPrint("LARGE: ")

	webPrint("LARGE:Bernoulli noticed that this sequence approaches a limit")
	webPrint("LARGE:(the force of interest) with larger n and, thus, smaller")
	webPrint("LARGE:compounding intervals. Compounding weekly (n = 52) yields")
	webPrint("LARGE:  $2.692596..., while compounding daily (n = 365) yields")
	webPrint("LARGE:  $2.714567... (approximately two cents more). The limit as")
	webPrint("LARGE: n grows large is the number that came to be known as e. ")
	webPrint("LARGE: ")

	webPrint("LARGE:That is, with continuous compounding, the account value will")
	webPrint("LARGE:reach $2.718281828")

	webPrint(pkg.BoxSep(50))
}