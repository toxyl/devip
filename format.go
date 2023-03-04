package main

import "fmt"

const (
	RESET     = "\033[0m"
	BOLD      = "\033[1m"
	ITALIC    = "\033[3m"
	UNDERLINE = "\033[4m"
	RED       = 1
	GREEN     = 2
	YELLOW    = 3
	BLUE      = 4
)

// as applies ANSI escape codes to format a string with bold, italic, and/or underline styles.
func as(str string, bold, italic, underline bool) string {
	formatStr := ""
	if bold {
		formatStr += BOLD
	}
	if italic {
		formatStr += ITALIC
	}
	if underline {
		formatStr += UNDERLINE
	}
	return formatStr + str + RESET
}

// asColor applies ANSI escape codes to color a string and optionally apply bold, italic, and/or underline styles.
func asColor(str string, color int, bold, italic, underline bool) string {
	return "\033[" + fmt.Sprint(30+color) + "m" + as(str, bold, italic, underline)
}

// asError applies ANSI escape codes to color a string red and make it bold.
func asError(str string) string {
	return asColor(str, RED, true, false, false)
}

// asOK applies ANSI escape codes to color a string green and make it bold.
func asOK(str string) string {
	return asColor(str, GREEN, true, false, false)
}

// asWarning applies ANSI escape codes to color a string yellow and make it italic.
func asWarning(str string) string {
	return asColor(str, YELLOW, false, true, false)
}

// asIP applies ANSI escape codes to color a string blue.
func asIP(str string) string {
	return asColor(str, BLUE, false, false, false)
}
