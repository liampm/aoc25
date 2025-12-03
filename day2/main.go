package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {

	blockNum := flag.Int("n", 0, "the number of repeating blocks - 0 means any")
	inputFile := flag.String("input", "input.txt", "input file path")
	flag.Parse()

	input, err := os.ReadFile(*inputFile)
	if err != nil {
		panic(err)
	}
	ranges := strings.Split(strings.TrimSpace(string(input)), ",")

	var blockSizer func(int) []int
	if *blockNum > 0 {
		blockSizer = blockSizerForNBlocks(*blockNum)
	} else {
		blockSizer = factorise // A block of each possible size
	}

	var invalidIds []int
	total := 0
	for _, rangeStr := range ranges {
		totalInRange, invalidIdsInRange, err := invalidInRange(rangeStr, blockSizer)
		if err != nil {
			log.Printf("Error processing range %s: %v", rangeStr, err)
			continue
		}
		total += totalInRange
		invalidIds = append(invalidIds, invalidIdsInRange...)
	}

	fmt.Printf("Invalid IDs: %v\n", invalidIds)
	fmt.Printf("Total: %d\n", total)
}

func invalidInRange(rangeStr string, blockSizer func(int) []int) (total int, ids []int, err error) {
	rangeParts := strings.Split(rangeStr, "-")
	if len(rangeParts) != 2 {
		err = fmt.Errorf("invalid range format: %s", rangeStr)
		return
	}
	start, err := strconv.Atoi(rangeParts[0])
	if err != nil {
		return
	}
	end, err := strconv.Atoi(rangeParts[1])
	if err != nil {
		return
	}
	if start >= end {
		err = fmt.Errorf("invalid range: %s", rangeStr)
		return
	}

	startDigits := len(rangeParts[0])

	if start < 10 {
		start = 11 // first repeated sequence
		startDigits = 2
	}

	if start > end {
		return // Silently skip this invalid range when the start has been coerced
	}

	digits := startDigits
	blockSizes := blockSizer(digits)
	for i := start; i <= end; i++ {
		// Revalidate the possible block sizes each time we increment a unit
		if i != start && i == int(math.Pow(10, float64(digits))) {
			digits++
			blockSizes = blockSizer(digits)
		}

		if blockSizes == nil {
			continue
		}

		id := []byte(strconv.Itoa(i))

		if !validateID(id, digits, blockSizes) {
			ids = append(ids, i)
			total += i
		}
	}
	return total, ids, nil
}

func validateID(id []byte, digits int, blockSizes []int) bool {
	for _, blockSize := range blockSizes {
		numberOfParts := digits / blockSize // e.g 6 digits has 3 parts when split into 2 digit blocks
		isInvalid := true
		for j := 0; j < (numberOfParts - 1); j++ {
			if !bytes.Equal(id[j*blockSize:(j+1)*blockSize], id[(j+1)*blockSize:(j+2)*blockSize]) {
				isInvalid = false
				break
			}
		}
		if isInvalid {
			return false // An ID can only be invalid once .e.g. 2222 shouldn't count twice
		}
	}
	return true
}

func factorise(n int) (factors []int) {
	if n == 0 || n == 1 {
		return nil
	}
	for i := 1; i < n; i++ {
		if n%i == 0 {
			factors = append(factors, i)
		}
	}
	return factors
}

func blockSizerForNBlocks(nBlocks int) func(int) []int {
	return func(inputLength int) []int {
		if inputLength%nBlocks == 0 {
			return []int{inputLength / nBlocks}
		} else {
			return nil
		}
	}
}
