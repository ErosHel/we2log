package log

import "fmt"

const (
	black = iota + 30
	red
	green
	yellow
	blue
	purple
	cyan
	white
)

func Black(str string) string {
	return textColor(black, str)
}

func Red(str string) string {
	return textColor(red, str)
}
func Yellow(str string) string {
	return textColor(yellow, str)
}
func Green(str string) string {
	return textColor(green, str)
}
func Cyan(str string) string {
	return textColor(cyan, str)
}
func Blue(str string) string {
	return textColor(blue, str)
}
func Purple(str string) string {
	return textColor(purple, str)
}
func White(str string) string {
	return textColor(white, str)
}

func textColor(color int, str string) string {
	return fmt.Sprintf("\x1b[0;%dm%s\x1b[0m", color, str)
}
