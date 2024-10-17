package color

import "fmt"

const consoleEscape = "\x1b"

type ConsoleColor int

const (
	None ConsoleColor = iota
	ColorRed
	colorGreen
	colorYellow
	colorBlue
	colorPurple
)

func Green(format string, a ...any) string {
	return color(colorGreen) + fmt.Sprintf(format, a...) + color(None)
}

func Yellow(format string, a ...any) string {
	return color(colorYellow) + fmt.Sprintf(format, a...) + color(None)
}

func Red(format string, a ...any) string {
	return color(ColorRed) + fmt.Sprintf(format, a...) + color(None)
}

func color(color ConsoleColor) string {
	if color == None {
		return fmt.Sprintf("%s[%dm", consoleEscape, color)
	}

	return fmt.Sprintf("%s[3%dm", consoleEscape, color)
}
