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
	hex := "#00000"
	colorName := "random"
	for _, r := range FlagText {
		_, _ = bw.WriteString("<span style=\"color:")
		if FlagInventColor {
			hex = htmlColor.InventColor(3*FlagMinBrightness, 3*FlagMaxBrightness)
		} else {
			colorName, hex = htmlColor.RandomColor(3 * FlagMaxBrightness)
		}
		_, _ = bw.WriteString(hex)
		if FlagAntiColor {
			_, _ = bw.WriteString("; background-color: ")
			_, _ = bw.WriteString(htmlColor.AntiColor(hex))
		}
		_, _ = bw.WriteString(";\">")
		_, _ = bw.WriteRune(r)
		_, _ = bw.WriteString("</span>")
		if FlagDebug && FlagVerbose {
			xLog.Printf("char %c random color %s", r, colorName)
		}
	}
	_, _ = bw.WriteString("</div>\n")
}
