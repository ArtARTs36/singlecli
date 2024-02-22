package color

import "fmt"

const consoleEscape = "\x1b"

type ConsoleColor int

const (
	colorNone ConsoleColor = iota
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorPurple
)

func Green(msg string) string {
	return color(colorGreen) + msg + color(colorNone)
}

func Yellow(msg string) string {
	return color(colorYellow) + msg + color(colorNone)
}

func color(color ConsoleColor) string {
	if color == colorNone {
		return fmt.Sprintf("%s[%dm", consoleEscape, color)
	}

	return fmt.Sprintf("%s[3%dm", consoleEscape, color)
}
