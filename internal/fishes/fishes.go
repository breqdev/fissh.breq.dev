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

func GetFish(max_width int, max_height int) string {
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
		num_lines := 0
		max_length := 0
		leading_spaces := -1

		for _, line := range strings.Split(string(fish), "\n") {
			num_lines += 1

			if len(line) > max_length {
				max_length = len(line)
			}

			line_leading_spaces := 0
			for _, char := range line {
				if char == ' ' {
					line_leading_spaces += 1
				} else {
					break
				}
			}

			if len(line) > line_leading_spaces {
				if leading_spaces == -1 {
					leading_spaces = line_leading_spaces
				} else {
					if line_leading_spaces < leading_spaces {
						leading_spaces = line_leading_spaces
					}
				}
			}
		}

		if num_lines < max_height && max_length < max_width {
			new_fish := ""
			for _, line := range strings.Split(string(fish), "\n") {
				if len(line) > leading_spaces {
					new_fish += line[leading_spaces:] + "\n"
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
