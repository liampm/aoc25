package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type mode string

const (
	FRESH     mode = "fresh"
	AVAILABLE mode = "available"
)

type containsResult int

const (
	UNDER containsResult = -1
	IN    containsResult = 0
	OVER  containsResult = 1
)

type idRange struct {
	start int
	end   int
}

func FromBytes(b []byte) (idRange, error) {

	parts := strings.Split(string(b), "-")
	if len(parts) != 2 {
		return idRange{}, fmt.Errorf("invalid range format: %s", b)
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		return idRange{}, fmt.Errorf("failed to parse start: %v", err)
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return idRange{}, fmt.Errorf("failed to parse end: %v", err)
	}
	return idRange{start, end}, nil
}

func (r idRange) contains(n int) containsResult {
	if n < r.start {
		return UNDER
	}
	if n > r.end {
		return OVER
	}
	return IN
}

func (r idRange) overlaps(other idRange) bool {
	return r.start <= other.end && other.start <= r.end
}

func (r idRange) merge(other idRange) (idRange, bool) {
	if !r.overlaps(other) {
		return r, false
	}
	if r.start > other.start {
		r.start = other.start
	}
	if r.end < other.end {
		r.end = other.end
	}
	return r, true
}

func (r idRange) length() int {
	return r.end - r.start + 1
}

var enableDebug *bool

func main() {
	inputFile := flag.String("f", "input.txt", "input data file")
	enableDebug = flag.Bool("d", false, "debug mode")
	flag.Parse()

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	processingMode := FRESH
	ranges := make([]idRange, 0)
	freshCount := 0
	freshAndAvailableCount := 0

	debug("Processing fresh data...\n")
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			//
			break
		} else if err != nil {
			log.Fatalf("failed to read line: %v", err)
		}

		if len(line) == 0 {

			ranges = sortAndCompactRanges(ranges)
			for _, r := range ranges {
				freshCount += r.length()
			}

			debug("Switching mode\n")
			debug("Processing available data...\n")
			processingMode = AVAILABLE
			continue
		}

		switch processingMode {
		case FRESH:
			r, err := FromBytes(line)
			if err != nil {
				log.Fatal(err)
			}
			ranges = append(ranges, r)

		case AVAILABLE:

			availableId, err := strconv.Atoi(string(line))
			if err != nil {
				log.Fatalf("failed to parse ID: %v", err)
			}
			for _, r := range ranges {
				c := r.contains(availableId)
				if c == UNDER {
					break
				} else if c == OVER {
					continue
				} else if c == IN {
					freshAndAvailableCount++
				}
			}
		}
	}

	fmt.Printf("Fresh count: %d\n", freshCount)
	fmt.Printf("Fresh and available count: %d\n", freshAndAvailableCount)
}

func sortAndCompactRanges(ranges []idRange) []idRange {
	slices.SortFunc(ranges, rangeSort)
	compacted := make([]idRange, 0)

	var lastRange idRange
	for i, r := range ranges {
		if i == 0 {
			lastRange = r
			continue
		}
		mergedRange, changed := r.merge(lastRange)
		if !changed {
			compacted = append(compacted, lastRange)
		}
		lastRange = mergedRange
	}
	compacted = append(compacted, lastRange)
	debug("ranges: %v\n", ranges)
	debug("compacted: %v\n", ranges)
	return compacted
}

func rangeSort(a, b idRange) int {
	if a.start < b.start {
		return -1
	} else if a.start > b.start {
		return 1
	}
	if a.end < b.end {
		return -1
	} else if a.end > b.end {
		return 1
	}
	return 0
}

func debug(format string, a ...any) (n int, err error) {
	if *enableDebug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
