package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
)

//go:embed challenge.in
var input string

// Represents a CPU instruction
const (
	addx = "addx"
	noop = "noop"
)

// Represents display constant
const (
	darkPixel = "░"
	litPixel  = "█"
)

// CRT represent a display for visual data output.
type CRT struct {
	width  int
	height int
	screen [][]string
}

// NewCRT returns a new state of the art CRT display.
func NewCRT(width, height int) *CRT {
	crt := &CRT{
		width:  width,
		height: height,
		screen: make([][]string, height),
	}
	for y := range crt.screen {
		crt.screen[y] = make([]string, width)
	}
	return crt
}

// DrawPixel updates a single pixel based on current CPU values.
func (crt *CRT) DrawPixel(cpuCycle, cpuRegister int) {
	row := cpuCycle / crt.width
	col := cpuCycle % crt.width
	// cpuRegister is our middle sprite position (+- 1)
	// R=1 -> 0 1 2
	// R=15 -> 14 15 16
	// etc.
	// Basically if column matches register (+-1) => pixel is lit
	crt.screen[row][col] = darkPixel
	if col+1 == cpuRegister || col-1 == cpuRegister || col == cpuRegister {
		crt.screen[row][col] = litPixel
	}
}

// Output returns a formated output of the display contents.
func (crt *CRT) Output() string {
	out := ""
	for y := 0; y < crt.height; y++ {
		for x := 0; x < crt.width; x++ {
			out += crt.screen[y][x]
		}
		out += "\n"
	}
	return out
}

// executionTimes contains number of cycles it takes to complete a command
var executionTimes = map[string]int{
	addx: 2,
	noop: 1,
}

// CPU is a virtual processor unit.
type CPU struct {
	cycle         int
	register      int
	signalStrengh int
	display       *CRT
}

// MakeCPU initializes and returns a new virtual CPU.
func MakeCPU() *CPU {
	return &CPU{
		cycle:    0,
		register: 1,
		display:  NewCRT(40, 6),
	}
}

// Execute runs a single instruction on the virtual CPU.
func (cpu *CPU) Execute(op instruction) {
	for op.cycles > 0 {
		op.cycles--
		cpu.display.DrawPixel(cpu.cycle, cpu.register)
		cpu.cycle++
		if cpu.cycle%40 == 20 {
			// fmt.Printf("R=%d, C=%d \n", cpu.register, cpu.cycle)
			cpu.signalStrengh += cpu.register * cpu.cycle
		}

		switch op.command {
		case noop:
			// Do nothing
		case addx:
			if op.cycles == 0 {
				cpu.register += op.argument
			}
		}
	}

}

// RunCode executes operations in order as they are passed.
func (cpu *CPU) RunCode(code []instruction) {
	for _, op := range code {
		cpu.Execute(op)
	}
}

// instruction represents our virtual CPU command description.
type instruction struct {
	command  string
	argument int
	cycles   int
}

// parseAssemblyLike creates instructions from the input files.
func parseAssemblyLike(code string) []instruction {
	lines := strings.Split(code, "\n")
	insts := make([]instruction, len(lines))

	for i, line := range lines {
		f := strings.Fields(line)
		cmd := f[0]
		var arg int
		if len(f) > 1 {
			arg, _ = strconv.Atoi(f[1])
		}
		insts[i] = instruction{
			command:  cmd,
			argument: arg,
			cycles:   executionTimes[cmd],
		}
	}
	return insts
}

// runChallenge returns the desired output for the days challenge.
// May print additional information to stdout.
func runChallenge() *CPU {
	instructions := parseAssemblyLike(input)
	cpu := MakeCPU()
	cpu.RunCode(instructions)
	return cpu
}

func main() {
	cpu := runChallenge()
	fmt.Println("====Part 1=====")
	fmt.Println(cpu.signalStrengh)
	fmt.Println("====Part 2=====")
	fmt.Println(cpu.display.Output())
}
