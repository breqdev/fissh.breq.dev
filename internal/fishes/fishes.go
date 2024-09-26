package fishes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

const Fish = `
                ,"(
               ////\                           _
              (//////--,,,,,_____            ,"
            _;"""----/////_______;,,        //
__________;"o,-------------......"""""` + "`" + `'-._/(
      ""'==._.__,;;;;"""           ____,.-.==
             "-.:______,...;---""/"   "    \(
                 '-._      ` + "`" + `-._("   ctr     \\
                     '-._                    '._
`

func GetFish(maxWidth int, maxHeight int) string {
	// fish are stored in the fishes/*.txt files

	// get the list of fish files
	files, err := os.ReadDir("fishes")
	if err != nil {
		fmt.Println(err)
	}

	// shuffle the list of fish files
	rand.Seed(uint64(time.Now().UnixNano()))
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })

	// iterate through them
	for _, file := range files {
		// read in the fish
		fish, err := os.ReadFile("fishes/" + file.Name())
		if err != nil {
			fmt.Println(err)
		}

		// check if the fish fits in the terminal
		numLines := 0
		maxLength := 0
		leadingSpaces := -1

		for _, line := range strings.Split(string(fish), "\n") {
			numLines += 1

			if len(line) > maxLength {
				maxLength = len(line)
			}

			lineLeadingSpaces := 0
			for _, char := range line {
				if char == ' ' {
					lineLeadingSpaces += 1
				} else {
					break
				}
			}

			if len(line) > lineLeadingSpaces {
				if leadingSpaces == -1 {
					leadingSpaces = lineLeadingSpaces
				} else {
					if lineLeadingSpaces < leadingSpaces {
						leadingSpaces = lineLeadingSpaces
					}
				}
			}
		}

		if numLines < maxHeight && maxLength < maxWidth {
			new_fish := ""
			for _, line := range strings.Split(string(fish), "\n") {
				if len(line) > leadingSpaces {
					new_fish += line[leadingSpaces:] + "\n"
				} else {
					new_fish += "\n"
				}
			}
			return new_fish
		}
	}

	// TODO this is bad
	return ""
}
