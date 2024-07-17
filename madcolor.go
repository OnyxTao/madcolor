package main

import (
	"bufio"
	htmlColor "madcolor/htmlcolor"
	"madcolor/misc"
	"os"
)

func main() {
	initLog("madcolor.log")
	defer closeLog()
	initFlags()
	htmlColor.HtmlColorsInitialize()

	bw := bufio.NewWriter(os.Stdout)
	defer misc.DeferError(bw.Flush)

	_, _ = bw.WriteString("<div>")
	for _, r := range FlagText {
		_, _ = bw.WriteString("<span style=\"color:")
		colorName, hex := htmlColor.RandomColor(3 * 160)
		_, _ = bw.WriteString(hex)
		_, _ = bw.WriteString(";\">")
		_, _ = bw.WriteRune(r)
		_, _ = bw.WriteString("</span>")
		if FlagDebug && FlagVerbose {
			xLog.Printf("char %c random color %s", r, colorName)
		}
	}
	_, _ = bw.WriteString("</div>\n")
}
