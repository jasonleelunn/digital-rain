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

type colour struct {
	bright string
	faint  string
}

type state struct {
	styledChars []rune
	baseChars   []rune
	positions   map[int]int
}

const (
	RESET       string = "\033[0m"
	CLEAR       string = "\033[%dA\033[%dD"
	HIDE_CURSOR string = "\033[?25l"
	SHOW_CURSOR string = "\033[?25h"

	BLACK_BG      string = "\033[48;5;16m"
	BLACK         string = "\033[38;5;16m"
	FAINT_RED     string = "\033[38;5;52m"
	BRIGHT_RED    string = "\033[38;5;196m"
	FAINT_GREEN   string = "\033[38;5;22m"
	BRIGHT_GREEN  string = "\033[38;5;46m"
	FAINT_YELLOW  string = "\033[38;5;58m"
	BRIGHT_YELLOW string = "\033[38;5;226m"
	FAINT_BLUE    string = "\033[38;5;24m"
	BRIGHT_BLUE   string = "\033[38;5;51m"
)

var (
	colours = map[string]colour{
		"red":    makeColour(FAINT_RED, BRIGHT_RED),
		"green":  makeColour(FAINT_GREEN, BRIGHT_GREEN),
		"yellow": makeColour(FAINT_YELLOW, BRIGHT_YELLOW),
		"blue":   makeColour(FAINT_BLUE, BRIGHT_BLUE),
	}

	colourChoice = flag.String("colour", "green", "Set the foreground colour")
)

func clearTerminal(width int, height int) {
	clearString := fmt.Sprintf(CLEAR, height, width)
	os.Stdout.Write([]byte(clearString))
}

func makeColour(faint string, bright string) colour {
	return colour{
		bright,
		faint,
	}
}

func getColour() colour {
	colour, ok := colours[*colourChoice]

	if !ok {
		// NOTE: passed an invalid colour, return the default
		return colours["green"]
	}

	return colour
}

func main() {
	flag.Parse()

	// Set up signal channel to capture interrupt signals
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		// ensure we set the users cursor visible again when quit
		os.Stdout.Write([]byte(SHOW_CURSOR))
		os.Exit(1)
	}()

	colour := getColour()

	// TODO: handle changing terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))

	if err != nil {
		log.Fatal(err)
	}

	state := state{
		positions: map[int]int{},
	}

	for col := range width {
		pos := rand.Intn(height)
		state.positions[col] = pos

		for range height {
			// random char (ASCII decimal 48 to 122)
			char := rune(rand.Intn(123-48) + 48)

			state.baseChars = append(state.baseChars, char)
		}
	}

	os.Stdout.Write([]byte(HIDE_CURSOR))

	for {
		state.styledChars = []rune{}

		for row := range height {
			for col := range width {

				pos := col + (row * width)
				char := state.baseChars[pos]

				var updatedChar []rune

				switch state.positions[col] {
				case row:
					updatedChar = append([]rune(colour.bright), char)
				case row + 1, row + 2, row + 3, row + 4, row + 5:
					updatedChar = append([]rune(colour.faint), char)
				default:
					updatedChar = append([]rune(BLACK), char)
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

		// TODO: timing options
		time.Sleep(time.Millisecond * 100)
	}
}
