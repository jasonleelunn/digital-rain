package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

type ansiEscapeCode string

type state struct {
	styledChars []rune
	baseChars   []rune
	positions   map[int]int
}

const (
	RESET  ansiEscapeCode = "\033[0m"
	BOLD   ansiEscapeCode = "\033[1m"
	RED    ansiEscapeCode = "\033[31m"
	GREEN  ansiEscapeCode = "\033[32m"
	YELLOW ansiEscapeCode = "\033[33m"
	BLUE   ansiEscapeCode = "\033[34m"
	WHITE  ansiEscapeCode = "\033[37m"

	DIM      ansiEscapeCode = "\033[38;2;11;32;17m"
	BLACK_BG ansiEscapeCode = "\033[48;2;0;0;0m"

	CLEAR       ansiEscapeCode = "\033[%dA\033[%dD"
	HIDE_CURSOR ansiEscapeCode = "\033[?25l"
	SHOW_CURSOR ansiEscapeCode = "\033[?25h"
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

func clearTerminal(width int, height int) {
	clearString := fmt.Sprintf(string(CLEAR), height, width)
	os.Stdout.Write([]byte(clearString))
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
	os.Stdout.Write([]byte(HIDE_CURSOR))

	// Set up signal channel to capture interrupt signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		// ensure we set the users cursor visible again when quit
		os.Stdout.Write([]byte(SHOW_CURSOR))
		os.Exit(1)
	}()

	flag.Parse()

	// colour := getColour()

	// TODO: handle changing terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		log.Fatal(err)
	}

	// totalChars := width * height

	state := state{
		positions: map[int]int{},
	}

	for col := range width {
		pos := rand.Intn(height)
		state.positions[col] = pos

		for range height {
			// random char (ASCII decimal 48 to 122)
			char := rune(rand.Intn(123-48) + 48)

			charWithColour := append([]rune(DIM), char)

			state.baseChars = append(state.baseChars, char)
			state.styledChars = append(state.styledChars, charWithColour...)
		}
	}

	for {
		state.styledChars = []rune{}

		for row := range height {
			for col := range width {

				pos := col + (row * width)
				char := state.baseChars[pos]

				var updatedChar []rune

				if row == state.positions[col] {
					updatedChar = append([]rune(WHITE), char)
					// updatedChar = append([]rune(BOLD), updatedChar...)
				} else {
					updatedChar = append([]rune(DIM), char)
				}

				state.styledChars = append(state.styledChars, updatedChar...)

			}

		}

		for col := range width {
			if state.positions[col] >= height {
				state.positions[col] = 0
			} else {
				state.positions[col]++
			}
		}

		// var sb strings.Builder
		//
		// sb.WriteString()
		//
		output := string(BLACK_BG) + string(state.styledChars) + string(RESET)

		// NOTE: clear the previous output just before we paint the new output
		// to try and prevent flickering
		clearTerminal(width, height)
		os.Stdout.WriteString(output)

		time.Sleep(time.Millisecond * 200)
	}
}
