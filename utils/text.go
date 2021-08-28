package utils

import "strings"

func WrapLine(text string, lineWidth int) (wrapped string) {
	if text == "" {
		return ""
	}

	words := strings.Split(text, " ")
	if len(words) == 0 {
		return
	}
	wrapped = words[0]
	spaceLeft := lineWidth - len(wrapped)
	for _, word := range words[1:] {
		if LenString(word)+1 > spaceLeft {
			wrapped += "\n" + word
			spaceLeft = lineWidth - LenString(word)
		} else {
			wrapped += " " + word
			spaceLeft -= 1 + LenString(word)
		}
	}

	return
}

func WrapText(text string, lineWidth int) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return ""
	}
	for i := range lines {
		lines[i] = WrapLine(lines[i], lineWidth)
	}

	return strings.Join(lines, "\n")

}
