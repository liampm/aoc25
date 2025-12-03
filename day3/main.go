package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

var enableDebug *bool

func main() {

	inputFile := flag.String("f", "input.txt", "input file path")
	joltageLength := flag.Int("j", 2, "length of joltage")
	enableDebug = flag.Bool("d", false, "debug mode")
	flag.Parse()

	if *joltageLength < 1 {
		log.Fatalf("joltageLength (%d) must be greater than 0", *joltageLength)
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	var joltageByBank []int
	totalJoltage := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("failed to read byte: %v", err)
		}

		bankJoltage := processBank(line, *joltageLength)

		joltageByBank = append(joltageByBank, bankJoltage)
		totalJoltage += bankJoltage
	}
	fmt.Printf("Joltage by Bank (%d): %v\n", len(joltageByBank), joltageByBank)
	fmt.Printf("Total Joltage: %d\n", totalJoltage)
}

func processBank(input []byte, joltageLength int) int {
	biggestNumbers := make([]int, joltageLength)
	bankLength := len(input)
	if joltageLength > bankLength {
		log.Fatalf("joltageLength (%d) is greater than bankLength (%d)", joltageLength, bankLength)
	}

	for i, r := range input {
		if r == '\n' {
			continue
		}

		digit, err := strconv.Atoi(string(r))
		if err != nil {
			log.Fatalf("failed to convert rune to int: %v", err)
		}

		remainingBankLength := bankLength - i
		remainingJoltageLength := joltageLength
		if remainingBankLength < remainingJoltageLength {
			remainingJoltageLength = remainingBankLength
		}
		forceZeroes := false

		// Go through the remaining joltage digits that could still change
		for j := joltageLength - remainingJoltageLength; j < joltageLength; j++ {
			if forceZeroes {
				biggestNumbers[j] = 0
			} else if digit > biggestNumbers[j] {
				biggestNumbers[j] = digit
				debug("digit: %d, %d: %v\n", digit, j, biggestNumbers)
				forceZeroes = true // Bigger than we had, any remaining smaller units will do
			}
		}
	}

	total := 0
	for i := 0; i < joltageLength; i++ {
		total += int(float64(biggestNumbers[i]) * math.Pow10(joltageLength-i-1))
	}

	return total
}

func debug(format string, a ...any) (n int, err error) {
	if *enableDebug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
