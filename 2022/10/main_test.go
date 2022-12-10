package main

import "testing"

const (
	expected1 = 15220
)

var expected2 = [][]string{{litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, darkPixel, darkPixel}, {litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel}, {litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, darkPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel}, {litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel}, {litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel}, {litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, litPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, litPixel, litPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, darkPixel, darkPixel, litPixel, darkPixel, darkPixel, litPixel, darkPixel}}

func TestChallenge1(t *testing.T) {
	cpu := runChallenge()

	if cpu.signalStrengh != expected1 {
		t.Errorf("Wrong result! Expected: %v, actual: %v", expected1, cpu.signalStrengh)
	}
}

func TestChallenge2(t *testing.T) {
	cpu := runChallenge()
	cpuExpected := MakeCPU()
	cpuExpected.display.screen = expected2
	for y := 0; y < cpu.display.height; y++ {
		for x := 0; x < cpu.display.width; x++ {
			if cpu.display.screen[y][x] != expected2[y][x] {
				t.Errorf("Fail on challenge 2:\n====Expected====\n%s==== Actual ====\n%s", cpuExpected.display.Output(), cpu.display.Output())
			}
		}
	}
}
