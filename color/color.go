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

func Green(format string, a ...any) string {
	return color(colorGreen) + fmt.Sprintf(format, a...) + color(colorNone)
}

func Yellow(format string, a ...any) string {
	return color(colorYellow) + fmt.Sprintf(format, a...) + color(colorNone)
}

func Red(format string, a ...any) string {
	return color(colorRed) + fmt.Sprintf(format, a...) + color(colorNone)
}

func color(color ConsoleColor) string {
	if color == colorNone {
		return fmt.Sprintf("%s[%dm", consoleEscape, color)
	}

	return fmt.Sprintf("%s[3%dm", consoleEscape, color)
}
