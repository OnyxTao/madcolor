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

var minContrast int = 40
var minColorDistance int = 30

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
	var br *bufio.Reader
	var err error
	var writerList []io.Writer
	var buffw bytes.Buffer

	// SETUP *****************************

	initLog("madcolor.log")
	defer closeLog()

	initializeClipboard()

	initFlags()

	htmlColor.Initialize()

	// LOGIC *************************

	br = getInput()

	f := getOutput()
	defer misc.DeferError(f.Close)
	writerList = append(writerList, f)

	if FlagClip {
		buffw.Grow(16 * 1024)
		writerList = append(writerList, &buffw)
	}

	mw := bufio.NewWriter(io.MultiWriter(writerList...))
	colorize(br, mw)

	err = mw.Flush()
	if nil != err {
		xLog.Printf("Could not flush bytes from buffered multiwriter because %s", err.Error())
	}
	if FlagClip {
		_ = clipboard.Write(clipboard.FmtText, buffw.Bytes())
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
	var fg, bg, colorName string
	var w = NewNLVWriter(out)

	w.WriteString("<div>")

	for r, _, err = in.ReadRune(); err == nil; r, _, err = in.ReadRune() {

		w.WriteString("<span style=\"color: ")

		if FlagInventColor {
			fg, bg = htmlColor.InventColor(FlagBackgroundColor, minContrast, minColorDistance)
			colorName = ""
		} else {
			colorName, fg = htmlColor.RandomColor(FlagBackgroundColor, minContrast, minColorDistance)
			bg = FlagBackgroundColor
		}

		_, _ = out.WriteString(fg)

		if FlagAntiColor {
			bg, fg = htmlColor.InventColor(fg, minContrast, minColorDistance)
			w.WriteString(";padding: 1px 0px; background-color: ")
			w.WriteString(bg)
		}
		w.WriteString(";\">")
		w.WriteRune(r)
		w.WriteString("</span>")
		if FlagDebug && FlagVerbose {
			if !misc.IsStringSet(&colorName) {
				colorName = fg
			}
			xLog.Printf("char %c background %s foreground %s", r, bg, colorName)
		}
	}
	w.WriteString("</div>\n")
}
