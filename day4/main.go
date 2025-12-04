package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type symbol string

const (
	EMPTY                   symbol = "."
	PAPER                   symbol = "@"
	ACCESSIBLE              symbol = "x"
	ACCESSIBILITY_THRESHOLD        = 4
)

var enableDebug *bool

func main() {
	inputFile := flag.String("f", "input.txt", "input file path")
	repeatParses := flag.Bool("r", false, "repeated mode")
	enableDebug = flag.Bool("d", false, "debug mode")
	flag.Parse()

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	line, _, err := reader.ReadLine()
	if err == io.EOF {
		log.Fatal("File is empty")
	} else if err != nil {
		log.Fatalf("failed to read line: %v", err)
	}

	var rows [][]bool

	rowNum := 0
	rows = append(rows, lineToCells(line))

	accessibleCount := 0
	lastParseCount := 0
	finished := false
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			finished = true
		} else if err != nil {
			log.Fatalf("failed to read line: %v", err)
		}
		if !finished {
			rows = append(rows, lineToCells(line))
			if len(rows[rowNum]) != len(rows[rowNum+1]) {
				log.Fatalf("row %d has different length than row %d", rowNum, rowNum+1)
			}
		}

		lastParseCount += processRow(&rows, rowNum, *repeatParses)
		rowNum++
		if finished {
			break
		}
	}

	accessibleCount += lastParseCount

	if *repeatParses {
		for {
			debug("\nParse again\n")
			lastParseCount = 0
			for i, _ := range rows {
				lastParseCount += processRow(&rows, i, *repeatParses)
			}
			accessibleCount += lastParseCount
			if lastParseCount == 0 {
				break
			}
		}
	}
	fmt.Printf("Accessible count: %d\n", accessibleCount)
}

func processRow(rows *[][]bool, rowNum int, live bool) int {
	accessibleCount := 0
	for i, c := range (*rows)[rowNum] {

		adjacentCount := 0
		if !c {
			debug(string(EMPTY))
			continue
		}

		columnCount := len((*rows)[rowNum])
		hasRowBefore := rowNum > 0
		hasRowAfter := rowNum < len(*rows)-1
		hasColumnBefore := i > 0
		hasColumnAfter := i < columnCount-1

		if hasColumnBefore {
			if hasRowBefore && (*rows)[rowNum-1][i-1] {
				adjacentCount++
			}
			if (*rows)[rowNum][i-1] {
				adjacentCount++
			}
			if hasRowAfter && (*rows)[rowNum+1][i-1] {
				adjacentCount++
			}
		}
		if hasRowBefore && (*rows)[rowNum-1][i] {
			adjacentCount++
		}
		if hasRowAfter && (*rows)[rowNum+1][i] {
			adjacentCount++
		}
		if hasColumnAfter {
			if hasRowBefore && (*rows)[rowNum-1][i+1] {
				adjacentCount++
			}
			if (*rows)[rowNum][i+1] {
				adjacentCount++
			}
			if hasRowAfter && (*rows)[rowNum+1][i+1] {
				adjacentCount++
			}
		}

		if adjacentCount < ACCESSIBILITY_THRESHOLD {
			debug(string(ACCESSIBLE))
			if live {
				(*rows)[rowNum][i] = false // The paper can be removed so update for next evaluation
			}
			accessibleCount++
		} else {
			debug(string(PAPER))
		}
	}
	debug("\n")
	return accessibleCount
}

func lineToCells(line []byte) []bool {
	columnCount := len(line)
	cells := make([]bool, columnCount)

	for i, c := range strings.Split(string(line), "") {
		cells[i] = c == string(PAPER)
	}

	return cells
}

func debug(format string, a ...any) (n int, err error) {
	if *enableDebug {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}
