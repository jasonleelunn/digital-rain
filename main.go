package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"golang.org/x/term"
)

type ansiEscapeCode string

const (
	RESET  ansiEscapeCode = "\033[0m"
	BOLD   ansiEscapeCode = "\033[1m"
	RED    ansiEscapeCode = "\033[31m"
	GREEN  ansiEscapeCode = "\033[32m"
	YELLOW ansiEscapeCode = "\033[33m"
	BLUE   ansiEscapeCode = "\033[34m"
)

var (
	colours = map[string]ansiEscapeCode{
		"red":    RED,
		"green":  GREEN,
		"yellow": YELLOW,
		"blue":   BLUE,
	}

	colourChoice = flag.String("colour", "green", "Set the output colour")
)

func clearTerminal() {
	// TODO: improve this to work cross platform
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getColour() ansiEscapeCode {
	colour, ok := colours[*colourChoice]

	if !ok {
		// NOTE: passed an invalid colour, return the default
		return GREEN
	}

	return colour
}

func main() {
	flag.Parse()

	colour := getColour()

	// TODO: handle changing terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		log.Fatal(err)
	}

	totalChars := width * height

	var state []rune

	for {
		clearTerminal()

		for range width {
			// empty space char
			char := rune(32)

			if rand.Intn(8) == 1 {
				// random char (ASCII decimal 48 to 122)
				char = rune(rand.Intn(123-48) + 48)
			}

			// prepend new char
			state = append([]rune{char}, state...)
		}

		// TODO: improve this...
		if len(state) >= totalChars {
			state = state[0:totalChars]
		}

		os.Stdout.WriteString(string(colour) + string(state) + string(RESET))

		time.Sleep(time.Millisecond * 100)
	}
}
