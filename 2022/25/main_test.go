package main

import (
	"testing"
)

const (
	expected1 = "2=-1=0"
)

func TestChallenge1(t *testing.T) {
	actual := runChallenge(1)

	if actual != expected1 {
		t.Errorf("Wrong result! Expected: %v, actual: %v", expected1, actual)
	}
}

func TestDecimalToSNAFU(t *testing.T) {
	expected := map[int]string{
		1:         "1",
		2:         "2",
		3:         "1=",
		4:         "1-",
		5:         "10",
		6:         "11",
		7:         "12",
		8:         "2=",
		9:         "2-",
		10:        "20",
		15:        "1=0",
		20:        "1-0",
		2022:      "1=11-2",
		12345:     "1-0---0",
		314159265: "1121-1110-1=0",
	}
	for dec, snafu := range expected {
		calculated := DecimalToSNAFU(dec)
		if calculated != snafu {
			t.Errorf("Wrong result! Expected: %v, calculated: %v", snafu, calculated)
		}
	}
}

func TestSNAFUToDecimal(t *testing.T) {
	expected := map[string]int{
		"1=-0-2": 1747,
		"12111":  906,
		"2=0=":   198,
		"21":     11,
		"2=01":   201,
		"111":    31,
		"20012":  1257,
		"112":    32,
		"1=-1=":  353,
		"1-12":   107,
		"12":     7,
		"1=":     3,
		"122":    37,
	}
	for snafu, dec := range expected {
		calculated := SNAFUToDecimal(snafu)
		if calculated != dec {
			t.Errorf("Wrong result! Expected: %v, calculated: %v", dec, calculated)
		}
	}
}

func Benchmark1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runChallenge(1)
	}
}
