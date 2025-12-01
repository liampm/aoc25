package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type direction int

const (
	L direction = -1
	R direction = 1
)

func main() {
	countAny := flag.Bool("any", false, "count any pass of 0")
	flag.Parse()

	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to open file: %w", err))
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	const pointsOnDial = 100
	runningValue := 50
	timesAtZero := 0
	lastChar := '\n'
	dir := R
	var moveBy strings.Builder

	for {
		c, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
			return
		}

		if lastChar == '\n' {
			switch c {
			case 'L':
				dir = L
			case 'R':
				dir = R
			default:
				log.Fatal(fmt.Sprintf("Incorrectly formatted line: Should be 'L' or 'R', got %q", c))
			}
			lastChar = c
			continue
		}

		lastChar = c

		if c == '\n' {
			moveByInt, moveByErr := strconv.Atoi(moveBy.String())
			if moveByErr != nil {
				log.Fatal(fmt.Errorf("Unable to parse moveBy (%s) to integer: %w", moveBy.String(), moveByErr))
			}
			moveBy.Reset()

			if *countAny {
				timesAtZero += moveByInt / pointsOnDial // 1 pass of zero for each full rotation
			}
			// apply the remaining moves after removing full rotations
			moveByInt %= pointsOnDial
			if moveByInt == 0 {
				// no more moves
				continue
			}
			oldValue := runningValue
			moveByInt *= int(dir)
			runningValue += moveByInt

			if runningValue > (pointsOnDial - 1) { // overflowing
				runningValue -= pointsOnDial
				if *countAny {
					timesAtZero++
					continue
				}
			} else if runningValue < 0 { // underflowing
				runningValue = pointsOnDial + runningValue
				if *countAny && oldValue != 0 { // hasn't passed zero if it started there
					timesAtZero++
					continue
				}
			}

			if runningValue == 0 {
				timesAtZero++
			}

		} else {
			moveBy.WriteRune(c)
		}
	}
	fmt.Printf("Times at zero: %d\n", timesAtZero)
}
