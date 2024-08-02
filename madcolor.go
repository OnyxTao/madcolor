package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"strings"

	"golang.design/x/clipboard"
	htmlColor "madcolor/htmlcolor"
	"madcolor/misc"
)

var modeClipboardAvailable = false

// const MAXATTEMPTS = 100

const DEBUGTEXTLOG = "madcolorDebugText.log"

var minContrast = 60
var minColorDistance = 33

func initializeClipboard() {
	err := clipboard.Init()
	if nil == err {
		modeClipboardAvailable = true
	} else {
		modeClipboardAvailable = false
		if FlagDebug || FlagVerbose {
			xLog.Printf("Clipboard not available because %s", err.Error())
		}
	}
}

// main runs the colorize program.
// It initializes the log file, closes it when done,
// initializes the command line flags, initializes and
// opens the input source, initializes and opens the output
// destination, initializes the HTML color names, and
// colorizes the input using the colorize function.
func main() {
	var err error
	var writerList []io.Writer
	var clipboardBuffer bytes.Buffer

	// SETUP *****************************

	initLog("madcolor.log")
	defer closeLog()

	initializeClipboard()

	initFlags()

	htmlColor.Initialize(FlagDebug)

	br := getInput()

	f := getOutput()
	defer misc.DeferError(f.Close)
	writerList = append(writerList, f)

	if FlagClip {
		clipboardBuffer.Grow(16 * 1024)
		writerList = append(writerList, &clipboardBuffer)
	}

	if FlagDebug {
		f, err := os.Open(DEBUGTEXTLOG)
		if err != nil {
			xLog.Printf("could not open %s because: %s", DEBUGTEXTLOG, err)
			myFatal()
		}
		defer misc.DeferError(f.Close)
		writerList = append(writerList, f)
	}

	mw := bufio.NewWriter(io.MultiWriter(writerList...))
	colorize(br, mw)
	err = mw.Flush()
	if nil != err {
		xLog.Printf("huh? Could not flush bytes from buffered multiwriter because %s",
			err.Error())
		myFatal()
	}

	if FlagClip {
		_ = clipboard.Write(clipboard.FmtText, clipboardBuffer.Bytes())
	}

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

	if FlagPipe {
		return os.Stdout
	}

	if misc.IsStringSet(&FlagOutput) {
		fn = path.Join(FlagOutputDir, FlagOutput)
		f, err = os.Create(fn)
		if err != nil {
			xLog.Printf("Could not open file %s because %s", fn, err.Error())
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
	if FlagPipe {
		return bufio.NewReader(os.Stdin)
	}

	if misc.IsStringSet(&FlagInput) {
		f, err := os.Open(FlagInput)
		if err != nil {
			xLog.Printf("Could not open %s because %s", FlagInput, err.Error())
			myFatal()
		}
		return bufio.NewReader(f)
	}

	return bufio.NewReader(strings.NewReader(FlagText))
}

// colorize applies color to each character read from the input.
// It calculates the background color if FlagBackgroundColor is not set,
// using a random color or a random named color from the htmlColor package.
// It then generates a foreground color that has enough contrast and distance
// from the background color.
// If FlagAntiColor is set, it generates a random background color for each character.
// It writes the colorized text to the output using the NLVWriter type.
// It wraps the colorized text in `<div>` tags.
//
// Parameters:
// - in: Input reader for reading characters
// - out: Output writer for writing colorized text
//
// Variables:
// - r: Represents each character read from the input (rune type)
// - err: Error returned from reading characters from the input
// - fg: Foreground color string
// - bg: Background color string
// - w: Writer that writes the colorized text to the output
// - colorName: Name of the generated foreground color
//
// Returns: None
func colorize(in *bufio.Reader, out *bufio.Writer) {
	var r rune
	var err error = nil
	var fg, bg string
	var w = NewNLVWriter(out)
	var colorName = ""

	// figure out the background color if not FlagAnticolor
	if !FlagAntiColor && !misc.IsStringSet(&FlagBackgroundColor) {
		err = nil
		if FlagInventColor {
			bg = htmlColor.RandColor()
			err = nFlags.Set("background-color", bg)
		} else {
			_, colorName, bg = htmlColor.RandNamedColor()
			err = nFlags.Set("background-color", colorName)
		}
		xLog.Printf("warning: FlagBackgroundColor was unset. "+
			"Should not happen! Creating a default background color ... %s %s",
			bg, colorName)
		if nil != err {
			xLog.Printf("Failed to set background-color color flag because %s", err.Error())
			myFatal()
		}
	}

	w.WriteString("<div>")

	for r, _, err = in.ReadRune(); err == nil; r, _, err = in.ReadRune() {

		w.WriteString("<span style=\"color: ")

		if FlagAntiColor {
			if FlagInventColor {
				bg = htmlColor.RandColor()
			} else {
				_, _, bg = htmlColor.RandNamedColor()
			}
		}

		if FlagInventColor {
			fg, _ = htmlColor.InventColor(bg, minContrast, minColorDistance)
		} else {
			colorName, fg = htmlColor.RandomColor(bg, minContrast, minColorDistance)
		}

		w.WriteString(fg)

		if FlagAntiColor {
			w.WriteString("; padding: 1px 0px; background-color: ")
			w.WriteString(bg)
		}

		w.WriteString(";\">")
		w.WriteRune(r)
		w.WriteString("</span>")
	}
	w.WriteString("</div>\n")
}
