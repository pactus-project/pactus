package util

import (
	"strconv"
	"strings"
)

func FormatIntWithDelimiters(num int64) string {
	numStr := strconv.FormatInt(num, 10)

	return formatNumberString(numStr)
}

func FormatFloatWithDelimiters(num float64, prec int) string {
	numStr := strconv.FormatFloat(num, 'f', prec, 64)

	parts := strings.Split(numStr, ".")
	numStr = parts[0]
	numStr = formatNumberString(numStr)

	if len(parts) > 1 {
		numStr += "." + parts[1]
	}

	return numStr
}

func formatNumberString(numStr string) string {
	var formattedNum string
	if strings.HasPrefix(numStr, "-") {
		formattedNum = "-"
		numStr = numStr[1:]
	}
	for i, c := range numStr {
		if (i > 0) && (len(numStr)-i)%3 == 0 {
			formattedNum += ","
		}
		formattedNum += string(c)
	}

	return formattedNum
}
