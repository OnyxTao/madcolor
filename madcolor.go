package main

import (
	"bufio"
	"bytes"
	"golang.design/x/clipboard"
	"io"
	htmlColor "madcolor/htmlcolor"
	"madcolor/misc"
	"os"
	"path"
	"strings"
)

var modeClipboardAvailable = false

const MAXATTEMPTS = 100

// very tweaky values, colordistance 85 / contrast 0.36
const MINCOLORDISTANCE float64 = 85.0
const MINCONTRAST float64 = 0.36

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
	err := clipboard.Init()
	if nil != err {
		modeClipboardAvailable = true
	} else {
		modeClipboardAvailable = false
		if FlagDebug || FlagVerbose {
			xLog.Printf("Clipboard not available because %s", err.Error())
		}
	}
	initFlags()
	htmlColor.Initialize()

	br = getInput()

	f := getOutput()
	defer misc.DeferError(f.Close)
	bw = bufio.NewWriter(f)
	defer misc.DeferError(bw.Flush)

	colorize(br, bw)
}

// getOutput returns a *os.File that represents the output destination.
// If the `FlagOutput` variable is set, `getOutput` creates a file with
// the specified // name in the directory specified by `FlagOutputDir`
// and returns it. If opening the file encounters an error, it logs the
// error message using `xLog.Printf` and calls `myFatal` to exit the
// program. If the `FlagOutput` variable is not set, `getOutput` returns
// `os.Stdout`. It does not return a buffered writer, because there would be
// no way to close the underlying file.
func getOutput() (f *os.File) {
	var fn string
	var err error
	if misc.IsStringSet(&FlagOutput) {
		fn = path.Join(FlagOutputDir, FlagOutput)
		f, err = os.Create(fn)
		if err != nil {
			xLog.Printf("Could not open %s because %s", fn, err.Error())
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
		return bufio.NewReader(f)
	}
	return bufio.NewReader(strings.NewReader(FlagText))
}

// colorize applies colors to characters read from the input reader
// and writes the colorized output to the output writer. It generates
// a random color and its hexadecimal representation using
// htmlColor.RandomColor. If FlagInventColor is set, it
// generates an invented color in the specified brightness range using
// htmlColor.InventColor. If the FlagDrift is set, it uses the
// antiColor generated in the previous iteration. Otherwise, it
// generates a new random color and its hexadecimal representation.
// If the FlagAntiColor is set, it generates an antiColor using
// htmlColor.AntiColor and checks the contrast and color differentiation.
// If the generated antiColor does not have enough contrast or color
// differentiation, it generates a new one. The function writes the color
// span tag and the colorized character to the output writer. If FlagDebug
// and FlagVerbose are set, it logs the random color for each character.
// The function stops reading if it encounters an error other than io.EOF
// and logs the error.
//
// Parameters:
// - in: the input reader from which characters are read
// - out: the output writer to which colorized output is written
func colorize(in *bufio.Reader, out *bufio.Writer) {
	var r rune
	var err error = nil
	var antiColor string
	var b bytes.Buffer

	if FlagClip {
		var writerList []io.Writer
		writerList = append(writerList, bufio.NewWriter(&b))
		multiWriter := io.MultiWriter(writerList...)
		out = bufio.NewWriter(multiWriter)
	}

	_, _ = out.WriteString("<div>")

	colorName, hex := htmlColor.RandomColor(3 * FlagMaxBrightness)
	_, antiColor = htmlColor.RandomColor(3 * FlagMaxBrightness)

	for r, _, err = in.ReadRune(); err == nil; r, _, err = in.ReadRune() {

		_, _ = out.WriteString("<span style=\"color: ")

		if FlagInventColor {
			hex = htmlColor.InventColor(3*FlagMinBrightness, 3*FlagMaxBrightness)
		} else if FlagDrift {
			hex = antiColor
		} else {
			colorName, hex = htmlColor.RandomColor(3 * FlagMaxBrightness)
		}
		_, _ = out.WriteString(hex)

		if FlagAntiColor || FlagDrift {
			var cnt = 0
			var cd float64
			var ccr float64
			if FlagDrift && !FlagInventColor {
				_, antiColor = htmlColor.RandomColor(3 * FlagMaxBrightness)
			} else {
				antiColor = htmlColor.AntiColor(hex)
			}

			cd, ccr = htmlColor.ColorDistance(hex, antiColor)
			// check contrast & color differentiation

			for (cd < MINCOLORDISTANCE || ccr < MINCONTRAST) && cnt < MAXATTEMPTS {
				xLog.Printf("%s vs %s: cd: %12f  ccr: %13f",
					hex, antiColor, cd, ccr)
				// if not enough contrast, try again ...
				if FlagInventColor {
					antiColor = htmlColor.InventColor(FlagMinBrightness, FlagMaxBrightness)
				} else {
					_, antiColor = htmlColor.RandomColor()
				}
				cd, ccr = htmlColor.ColorDistance(hex, antiColor)
				// but not forever!
				cnt++
			}
			if cnt >= MAXATTEMPTS {
				xLog.Printf(
					"huh? Could not get a good contrasting color for %s %s (tried %d times!)",
					colorName, hex, cnt)
			}
			_, _ = out.WriteString(";padding: 1px 0px 1px 0px; background-color: ")
			_, _ = out.WriteString(antiColor)
		}
		_, _ = out.WriteString(";\">")
		_, _ = out.WriteRune(r)
		_, _ = out.WriteString("</span>")
		if FlagDebug && FlagVerbose {
			xLog.Printf("char %c random color %s", r, colorName)
		}
	}
	if err != io.EOF {
		xLog.Printf("Failed to write colorized string to output because %s", err.Error())
		myFatal()
	}
	_, _ = out.WriteString("</div>\n")

	if FlagClip {
		_ = clipboard.Write(clipboard.FmtText, b.Bytes())
	}

}
