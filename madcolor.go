package main

import (
	"bufio"
	"io"
	htmlColor "madcolor/htmlcolor"
	"madcolor/misc"
	"os"
	"path"
	"strings"
)

// main runs the colorize program.
// It initializes the log file, closes it when done,
// initializes the command line flags, initializes and
// opens the input source, initializes and opens the output
// destination, initializes the HTML color names, and
// colorizes the input using the colorize function.
func main() {
	var bw *bufio.Writer
	var br *bufio.Reader

	initLog("madcolor.log")
	defer closeLog()
	initFlags()
	htmlColor.HtmlColorsInitialize()

	br = getInput()

	f := getOutput()
	defer misc.DeferError(f.Close)
	bw = bufio.NewWriter(f)
	defer misc.DeferError(bw.Flush)

	colorize(br, bw)
}

// getOutput returns a *os.File that represents the output destination.
// If the `FlagOutput` variable is set, `getOutput` creates a file with the specified
// name in the directory specified by `FlagOutputDir` and returns it. If opening the
// file encounters an error, it logs the error message using `xLog.Printf` and calls
// `myFatal` to exit the program. If the `FlagOutput` variable is not set, `getOutput`
// returns `os.Stdout`.
func getOutput() (f *os.File) {
	var fn string
	var err error
	if misc.IsStringSet(&FlagOutput) {
		fn = path.Join(FlagOutputDir, FlagOutput)
		f, err = os.Create(fn)
		if err != nil {
			xLog.Printf("Could not open %s because %s", &fn, err.Error())
			myFatal()
		}
	} else {
		f = os.Stdout
	}
	return f
}

// getInput returns a *bufio.Reader that reads from either the
// file specified by the `FlagInput` variable or from a string
// specified by the `FlagText` variable. If the `FlagInput` variable
// is set, `getInput` opens the file and creates a `bufio.Reader` to
// read from it. If opening the file encounters an error, it logs
// the error message using `xLog.Printf` and calls `myFatal` to exit
// the program. If the `FlagInput` variable is not set, `getInput`
// creates a `bufio.Reader` to read from the string specified by the
// `FlagText` variable. The `bufio.Reader` is then returned.
func getInput() (br *bufio.Reader) {
	if misc.IsStringSet(&FlagInput) {
		f, err := os.Open(FlagInput)
		if err != nil {
			xLog.Printf("Could not open %s because %s", &FlagInput, err.Error())
			myFatal()
		}
		br = bufio.NewReader(f)
	} else {
		br = bufio.NewReader(strings.NewReader(FlagText))
	}
	return br
}

// colorize performs colorization on input runes and writes the colorized output to the provided bufio.Writer.
// It starts the colorization with the "<div>" tag and ends it with the "</div>\n" tag.
// For each input rune, it generates a random color or invents a color based on the provided flags.
// The color and anti-color (if FlagAntiColor is true) are applied to the corresponding span tags.
// Each span tag encloses a single input rune.
// If FlagDebug and FlagVerbose are both true, it logs the character and the randomly generated color.
// If there is an error while writing to the output, it logs the error and calls myFatal to exit the program.
// The function assumes that the input rune stream is provided by the bufio.Reader and the colorized output
// should be written to the bufio.Writer.
func colorize(in *bufio.Reader, bw *bufio.Writer) {
	_, _ = bw.WriteString("<div>")
	hex := "#00000"
	colorName := "random"

	var r rune
	var err error
	for r, _, err = in.ReadRune(); err != nil; r, _, err = in.ReadRune() {
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
	if nil != err && err != io.EOF {
		xLog.Printf("Failed to write colorized string to output because %s", err.Error())
		myFatal()
	}
	_, _ = bw.WriteString("</div>\n")
}
