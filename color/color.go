package color

import "fmt"

const consoleEscape = "\x1b"

type ConsoleColor int

const (
	None ConsoleColor = iota
	Red
	Green
	Yellow
	Blue
	Purple
)

func Greenf(format string, a ...any) string {
	return color(Green) + fmt.Sprintf(format, a...) + color(None)
}

func Yellowf(format string, a ...any) string {
	return color(Yellow) + fmt.Sprintf(format, a...) + color(None)
}

func Redf(format string, a ...any) string {
	return color(Red) + fmt.Sprintf(format, a...) + color(None)
}

func color(color ConsoleColor) string {
	if color == None {
		return fmt.Sprintf("%s[%dm", consoleEscape, color)
	}

	return fmt.Sprintf("%s[3%dm", consoleEscape, color)
}
