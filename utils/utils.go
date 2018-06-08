package utils

import (
	"os"
	"unicode"
	"fmt"
	"math"
)

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err)  {
		return false
	}
	return true
}

func ParseFloat(text string) float64 {
	if len(text) < 6 {
		return 0
	}
	i1 := float64(text[0] - '0')*10000
	i2 := float64(text[2] - '0')*1000
	i3 := float64(text[3] - '0')*100
	i4 := float64(text[4] - '0')*10
	i5 := float64(text[5] - '0')

	result := (i1 + i2 + i3 + i4 + i5)/10000
	return result
}

func FormatName(name string) string {
	data := []rune(name)
	for _, v := range data {
		if !unicode.Is(unicode.Scripts["Han"], v) {
			name += "."
		}
	}
	for i := len(data); i < 16; i++ {
		name += ".."
	}
	return name
}

func GetAbsText(val2 float64) string {
	valText2 := ""
	if math.Abs(val2) >= 10 {
		valText2 = fmt.Sprintf("%0.1f", val2 + 0.05);
	} else {
		valText2 = fmt.Sprintf("%0.2f", val2 + 0.005);
	}
	if val2 >= 0 {
		valText2 = "\033[1;31m" + "+" + valText2 + "\033[0m"
	} else {
		valText2 = "\033[1;32m" + valText2 + "\033[0m"
	}
	return valText2
}

func GetAbsText2(val2 float64) string {
	valText2 := ""
	if math.Abs(val2) >= 10 {
		valText2 = fmt.Sprintf("%0.1f", val2 + 0.05);
	} else {
		valText2 = fmt.Sprintf("%0.2f", val2 + 0.005);
	}
	if val2 >= 0 {
		valText2 = "\033[1;31m" + "+" + valText2 + "\033[0m"
	} else {
		valText2 = "\033[1;32m" + valText2 + "\033[0m"
	}
	return valText2
}