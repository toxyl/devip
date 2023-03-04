package main

import "fmt"

// as applies ANSI escape codes to format a string with bold, italic, and/or underline styles.
func as(str string, bold, italic, underline bool) string {
	formatStr := ""
	if bold {
		formatStr += "\033[1m"
	}
	if italic {
		formatStr += "\033[3m"
	}
	if underline {
		formatStr += "\033[4m"
	}
	return formatStr + str + "\033[0m"
}

// asColor applies ANSI escape codes to color a string and optionally apply bold, italic, and/or underline styles.
func asColor(str string, color int, bold, italic, underline bool) string {
	formatStr := "\033[" + fmt.Sprint(30+color) + "m"
	if bold {
		formatStr += "\033[1m"
	}
	if italic {
		formatStr += "\033[3m"
	}
	if underline {
		formatStr += "\033[4m"
	}
	return formatStr + str + "\033[0m"
}

// asBold applies ANSI escape codes to make a string bold and optionally apply italic and/or underline styles.
func asBold(str string, italic bool, underline bool) string {
	return as(str, true, italic, underline)
}

// asItalic applies ANSI escape codes to make a string italic and optionally apply bold and/or underline styles.
func asItalic(str string, bold bool, underline bool) string {
	return as(str, bold, true, underline)
}

// asUnderline applies ANSI escape codes to make a string underlined and optionally apply bold and/or italic styles.
func asUnderline(str string, bold bool, italic bool) string {
	return as(str, bold, italic, true)
}

// asError applies ANSI escape codes to color a string red and make it bold.
func asError(str string) string {
	return asColor(str, 1, true, false, false)
}

// asWarning applies ANSI escape codes to color a string yellow and make it italic.
func asWarning(str string) string {
	return asColor(str, 3, false, true, false)
}

// asOK applies ANSI escape codes to color a string green and make it bold.
func asOK(str string) string {
	return asColor(str, 2, true, false, false)
}

// asNeutral applies ANSI escape codes to color a string blue.
func asNeutral(str string) string {
	return asColor(str, 4, false, false, false)
}
