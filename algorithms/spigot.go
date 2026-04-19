package algorithms

import (
	"fmt"
	"strconv"
	"time"

	"piandfriends/pkg"
)

const spigotDigits = 850
const spigotLineWidth = 50
const feynmanPosition = 762

// Spigot runs the Rabinowitz–Wagon spigot algorithm twice.
// Run 1: full speed, no uncertainty display.
// Run 2: honest edition with red '?' for uncertain digits,
//        human-paced for visibility, and Feynman Point detection.
func Spigot(done chan bool, webPrint func(string)) {
	bw := 50

	webPrint("COLOR:cyan:" + pkg.BoxSep(bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  THE RABINOWITZ–WAGON SPIGOT ALGORITHM ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  Pi from integer arithmetic alone ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  1995 — produces digits sequentially ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxSep(bw))
	webPrint("")

	// Run 1
	webPrint("COLOR:yellow:" + pkg.BoxSep(bw))
	webPrint("COLOR:yellow:" + pkg.BoxLine("  RUN 1 — Full speed ", bw))
	webPrint("COLOR:yellow:" + pkg.BoxSep(bw))
	webPrint("")

	run1Start := time.Now()
	if !spigotRun1(done, webPrint) {
		return
	}
	run1Time := time.Since(run1Start)

	webPrint("")
	webPrint("COLOR:yellow:" + pkg.BoxSep(bw))
	webPrint("COLOR:yellow:" + pkg.BoxLine(fmt.Sprintf("  RUN 1 COMPLETE in %s", run1Time.Round(time.Millisecond)), bw))
	webPrint("COLOR:yellow:" + pkg.BoxSep(bw))
	webPrint("")
	webPrint("  Or was that too fast to follow?")
	webPrint("")
	webPrint("  We will now run the algorithm again — a fresh computation —")
	webPrint("  with deliberate pauses so you can see the uncertainty.")
	webPrint("")
	webPrint("COLOR:red:  Uncertain digits appear in red as '?'.")
	webPrint("  Each '?' is overwritten the moment the algorithm resolves it.")
	webPrint("")
	webPrint("COLOR:yellow:  Watch around digit 762. Something famous lives there.")
	webPrint("")

	time.Sleep(3 * time.Second)

	// Run 2
	webPrint("COLOR:green:" + pkg.BoxSep(bw))
	webPrint("COLOR:green:" + pkg.BoxLine("  RUN 2 — Honest Edition ", bw))
	webPrint("COLOR:green:" + pkg.BoxLine("  Uncertainty shown in real time ", bw))
	webPrint("COLOR:green:" + pkg.BoxSep(bw))
	webPrint("")

	if !spigotRun2(done, webPrint) {
		return
	}

	webPrint("")
	webPrint("COLOR:cyan:" + pkg.BoxSep(bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  SPIGOT COMPLETE: 850 digits of π ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxLine("  Integer arithmetic only. ", bw))
	webPrint("COLOR:cyan:" + pkg.BoxSep(bw))
}

func spigotRun1(done chan bool, webPrint func(string)) bool {
	size := spigotDigits*10/3 + 50
	a := make([]int, size)
	for i := range a {
		a[i] = 2
	}

	line := ""
	pre := -1
	nines := 0
	count := 0
	decInserted := false

	addChar := func(ch string) {
		line += ch
		if len([]rune(line)) >= spigotLineWidth {
			webPrint(line)
			line = ""
		}
	}

	emit := func(d int) {
		if count >= spigotDigits {
			return
		}
		if !decInserted && count == 1 {
			addChar(".")
			decInserted = true
		}
		addChar(strconv.Itoa(d))
		count++
	}

	for count < spigotDigits {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum := 0
		for j := size - 1; j >= 0; j-- {
			a[j] *= 10
			sum = a[j] + carriedOver
			quotient := sum / (j*2 + 1)
			a[j] = sum % (j*2 + 1)
			carriedOver = quotient * j
		}
		a[0] = sum % 10
		q := sum / 10

		switch {
		case q == 9:
			nines++
		case q == 10:
			if pre >= 0 {
				emit(pre + 1)
			}
			for i := 0; i < nines && count < spigotDigits; i++ {
				emit(0)
			}
			pre = 0
			nines = 0
		default:
			if pre >= 0 {
				emit(pre)
			}
			for i := 0; i < nines && count < spigotDigits; i++ {
				emit(9)
			}
			pre = q
			nines = 0
		}
	}

	if line != "" {
		webPrint(line)
	}
	return true
}

func spigotRun2(done chan bool, webPrint func(string)) bool {
	size := spigotDigits*10/3 + 50
	a := make([]int, size)
	for i := range a {
		a[i] = 2
	}

	line := ""
	pendingQs := 0
	pre := -1
	nines := 0
	count := 0
	decInserted := false
	feynmanFired := false

	baseDelay := 8 * time.Millisecond
	humanPause := 750 * time.Millisecond

	show := func() {
		if pendingQs > 0 {
			webPrint("UPDATE:HASRED:" + line)
		} else {
			webPrint("UPDATE:" + line)
		}
	}

	newRow := func() {
		webPrint("")
		line = ""
		pendingQs = 0
	}

	showQ := func() {
		if !decInserted && count == 1 {
			line += "."
			decInserted = true
		}
		line += "?"
		pendingQs++
		show()
		time.Sleep(humanPause)
	}

	confirmAndEmit := func(preDigit int, resolvedChar rune) {
		runes := []rune(line)
		qStart := len(runes) - pendingQs
		prefix := make([]rune, qStart)
		copy(prefix, runes[:qStart])

		if preDigit >= 0 && count < spigotDigits {
			if !decInserted && count == 1 {
				prefix = append(prefix, '.')
				decInserted = true
			}
			prefix = append(prefix, rune('0'+preDigit))
			count++
		}
		prefixLen := len(prefix)

		pending := make([]rune, pendingQs)
		copy(pending, runes[qStart:])
		line = string(prefix) + string(pending)
		show()

		for i := 0; pendingQs > 0 && count < spigotDigits; i++ {
			allRunes := []rune(line)
			allRunes[prefixLen+i] = resolvedChar
			line = string(allRunes)
			pendingQs--
			count++
			show()
			time.Sleep(baseDelay / 2)
		}
		pendingQs = 0
		if len([]rune(line)) >= spigotLineWidth {
			newRow()
		}
	}

	webPrint("")

	for count < spigotDigits {
		select {
		case <-done:
			return false
		default:
		}

		carriedOver := 0
		sum := 0
		for j := size - 1; j >= 0; j-- {
			a[j] *= 10
			sum = a[j] + carriedOver
			quotient := sum / (j*2 + 1)
			a[j] = sum % (j*2 + 1)
			carriedOver = quotient * j
		}
		a[0] = sum % 10
		q := sum / 10

		switch {
		case q == 9:
			nines++
			showQ()

		case q == 10:
			savedNines := nines
			carryPre := -1
			if pre >= 0 {
				carryPre = pre + 1
			}
			confirmAndEmit(carryPre, '0')
			nines = 0

			if savedNines > 0 {
				if line != "" {
					newRow()
				}
				webPrint("COLOR:red:  [CARRY] " + strconv.Itoa(savedNines) + " uncertain digit(s) resolved to 0")
				time.Sleep(baseDelay * 3)
				webPrint("")
				webPrint("")
			}
			pre = 0

		default:
			savedNines := nines
			confirmAndEmit(pre, '9')
			nines = 0

			if !feynmanFired && savedNines >= 6 && count >= feynmanPosition-6 && count <= feynmanPosition+6 {
				feynmanFired = true
				if line != "" {
					newRow()
				}
				webPrint("")
				webPrint("COLOR:yellow:  +--------------------------------------------------+")
				webPrint("COLOR:yellow:  | !! THE FEYNMAN POINT !!                          |")
				webPrint("COLOR:yellow:  | Six consecutive 9s at decimal position 762       |")
				webPrint("COLOR:yellow:  | You just watched them hold as '??????'           |")
				webPrint("COLOR:yellow:  | before resolving to 999999.                      |")
				webPrint("COLOR:yellow:  |                                                  |")
				webPrint("COLOR:yellow:  | Feynman: \"nine nine nine nine nine nine...       |")
				webPrint("COLOR:yellow:  |          ...and so on!\"                          |")
				webPrint("COLOR:yellow:  +--------------------------------------------------+")
				webPrint("")
				time.Sleep(4 * time.Second)
				webPrint("")
			}

			pre = q
			time.Sleep(baseDelay)
		}
	}

	if pre >= 0 && count < spigotDigits {
		confirmAndEmit(pre, '9')
	}
	if line != "" {
		newRow()
	}
	return true
}