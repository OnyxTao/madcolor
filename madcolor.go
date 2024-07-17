package main

import (
	htmlcolors "madcolor/htmlcolor"
	"os"
	"strings"
)

func main() {
	initLog("madcolor.log")
	defer closeLog()
	initFlags()
	htmlcolors.HtmlColorsInitialize()

	var sb strings.Builder

	sb.WriteString("<div>")
	for _, r := range FlagText {
		sb.WriteString("<span style=\"color:")
		colorName, hex := htmlcolors.RandomColor(3 * 160)
		sb.WriteString(hex)
		sb.WriteString(";\">")
		sb.WriteRune(r)
		sb.WriteString("</span>")
		if FlagDebug && FlagVerbose {
			xLog.Printf("char %c random color %s", r, colorName)
		}
	}
	sb.WriteString("</div>\n")
	_, err := os.Stdout.WriteString(sb.String())
	if err != nil {
		xLog.Printf("could not colorize [%s]\nbecause %s", FlagText, err.Error())
	}
}
