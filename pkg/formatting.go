package pkg

import (
	"math/big"
	"strconv"
	"strings"
)

func FormatBigFloatWithCommas(num *big.Float) string {
	numInt, _ := num.Int(nil)
	numStr := numInt.String()

	prefix := ""
	if strings.HasPrefix(numStr, "-") {
		prefix = "-"
		numStr = numStr[1:]
	}

	var result strings.Builder
	for i, ch := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}
	return prefix + result.String()
}

func FormatIntWithCommas(num int64) string {
	numStr := strconv.FormatInt(num, 10)

	prefix := ""
	if strings.HasPrefix(numStr, "-") {
		prefix = "-"
		numStr = numStr[1:]
	}

	var result strings.Builder
	for i, ch := range numStr {
		if i > 0 && (len(numStr)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}
	return prefix + result.String()
}

func FormatFloat64WithCommas(num float64) string {
	numStr := strconv.FormatFloat(num, 'f', -1, 64)

	prefix := ""
	if strings.HasPrefix(numStr, "-") {
		prefix = "-"
		numStr = numStr[1:]
	}

	parts := strings.Split(numStr, ".")
	intPart := parts[0]

	var result strings.Builder
	for i, ch := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(ch)
	}

	if len(parts) > 1 {
		result.WriteRune('.')
		result.WriteString(parts[1])
	}

	return prefix + result.String()
}

func BoxSep(width int) string {
	return "+" + strings.Repeat("-", width) + "+"
}

func BoxLine(content string, width int) string {
	runeLen := len([]rune(content))
	if runeLen > width {
		content = string([]rune(content)[:width])
		runeLen = width
	}
	return "|" + content + strings.Repeat(" ", width-runeLen) + "|"
}