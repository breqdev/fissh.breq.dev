package fishes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/exp/rand"
)

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
		lines := strings.Split(string(fish), "\n")
		numLines := len(lines)
		maxLength := 0
		for _, line := range lines {
			length := len(line)
			if length > maxLength {
				maxLength = length
			}
		}

		if numLines < maxHeight && maxLength < maxWidth {
			new_fish := ""
			for _, line := range lines {
				if len(line) < maxWidth {
					line += strings.Repeat(" ", maxLength-len(line))
				}
				new_fish += line + "\n"
			}
			return new_fish
		}
	}

	// TODO this is bad
	return ""
}
