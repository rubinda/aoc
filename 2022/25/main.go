package main

import (
	_ "embed"
	"fmt"
	"math"
	"strings"
)

var (
	//go:embed example.in
	input string
)

var (
	// bitConversionDecimal converts a single SNAFU bit to a decimal digit.
	bitConversionDecimal = map[rune]int{
		'2': 2,
		'1': 1,
		'0': 0,
		'-': -1,
		'=': -2,
	}
	// bitConversionSNAFU converts a single decimal digit it to a SNAFU bit.
	bitConversionSNAFU = map[int]rune{
		0: '0',
		1: '1',
		2: '2',
		3: '=',
		4: '-',
	}
)

// reverse returns the reverse of a string.
func reverse(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// SNAFUToDecimal converts a SNAFU number to the decimal value.
func SNAFUToDecimal(snafuValue string) int {
	decimalValue := 0
	for b, bitVal := range reverse(snafuValue) {
		decimalValue += bitConversionDecimal[bitVal] * int(math.Pow(5, float64(b)))
	}
	return decimalValue
}

// DecimalToSNAFU converts a decimal number to the SNAFU value.
func DecimalToSNAFU(value int) string {
	base5 := make([]int, 0)
	remainder := 0
	for value > 0 {
		remainder = value % 5
		base5 = append(base5, remainder)
		value /= 5
	}
	carry := 0
	snafu := ""
	for _, b := range base5 {
		b += carry
		carry = 0
		if b > 4 {
			b = 0
			carry = 1
		} else if b > 2 {
			carry = 1
		}
		snafu = string(bitConversionSNAFU[b]) + snafu
	}
	if carry == 1 {
		snafu = "1" + snafu
	}
	return snafu
}

// runChallenge returns the desired output for the day's challenge.
func runChallenge(challengePart int) string {
	if challengePart == 1 {
		snafus := strings.Split(input, "\n")
		fuelRequirement := 0
		for _, snafu := range snafus {
			fuelRequirement += SNAFUToDecimal(snafu)
		}
		return DecimalToSNAFU(fuelRequirement)
	}
	return ""
}

func main() {
	fmt.Println(runChallenge(1))
}
