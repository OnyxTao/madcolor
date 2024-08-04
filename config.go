package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/spf13/pflag"
	"madcolor/misc"
)

const DEFAULTCOLORTEXT = "We promptly judged antique ivory buckles for the next prize !@#$%^*(){}[];:.,?"

// wordSepNormalizeFunc all options are lowercase, so
// ... lowercase they shall be
func wordSepNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(strings.ToLower(name))
}

var nFlags *pflag.FlagSet

/* secret flags */
var FlagSlow bool

/* standard flags */

var FlagHelp bool
var FlagQuiet bool
var FlagVerbose bool
var FlagDebug bool

/* program flags */
var FlagBackgroundColor string
var FlagContrast int8
var FlagText string
var FlagInventColor bool
var FlagAntiColor bool
var FlagOutput string
var FlagOutputDir string
var FlagInput string
var FlagClip bool
var FlagStdout bool
var FlagPipe bool
var FlagDistance int8 = 20
var FlagClipboardBuffer bool

// initFlags initializes the command line flags for the program.
// It sets up the flag set, defines the flags, and parses the command line arguments.
func initFlags() {
	var err error

	// [output-name][opt-name]
	hideFlags := make(map[string]string, 8)

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.SetNormalizeFunc(wordSepNormalizeFunc)

	// secret flags
	nFlags.StringVarP(&FlagOutputDir, "output-dir", "", "",
		"output directory for output file")
	// hideFlags["FlagOutputDir"] = "output-dir"

	// standard flags

	nFlags.BoolVarP(&FlagSlow, "slow", "", false,
		"Add some time between http calls (do not hammer server)")

	nFlags.BoolVarP(&FlagDebug, "debug", "d",
		false, "Enable additional informational and operational logging output for debug purposes")

	nFlags.BoolVarP(&FlagVerbose, "verbose", "v",
		false, "Supply additional run messages; use --debug for more information")

	nFlags.BoolVarP(&FlagHelp, "help", "h",
		false, "Display help message and usage information")

	nFlags.BoolVarP(&FlagQuiet, "quiet", "q",
		true, "Suppress log output to stdout and stderr (output still goes to logfile)")

	// program flags

	nFlags.Int8VarP(&FlagContrast, "contrast", "c", int8(minContrast),
		"minimum relative contrast between foreground and background")

	nFlags.Int8VarP(&FlagDistance, "distance", "D", int8(minColorDistance),
		"minimum relative contrast between foreground and background")

	nFlags.StringVarP(&FlagBackgroundColor, "background-color", "b", "white",
		"Background color. Ignored for --anti.")

	nFlags.BoolVarP(&FlagClipboardBuffer, "buff", "", false,
		"buffer mode -- convert text in the clipboard buffer")

	nFlags.BoolVarP(&FlagPipe, "pipe", "p", false,
		"Pipe mode; read from STDIN, write to STDOUT, all other io disabled.")

	// this USUALLY defaults to TRUE, but is set to FALSE if the user does not
	// explicitly specify it when using --output.
	nFlags.BoolVarP(&FlagStdout, "stdout", "", true,
		"Write to STDOUT as well as the output file")

	nFlags.BoolVarP(&FlagClip, "nopaste", "", true,
		"Suppress paste of buffer to clipboard (if clipboard is available)")

	nFlags.BoolVarP(&FlagAntiColor, "anti", "a", false,
		"Set the colorspace background to the foreground complement "+
			"or something random with minimum contrast (see --contrast)")

	nFlags.BoolVarP(&FlagInventColor, "invent", "I", false,
		"randomly generate colors (rather than randomly select websafe colors)")

	nFlags.StringVarP(&FlagText, "text", "t",
		DEFAULTCOLORTEXT, "Text to colorize")

	nFlags.StringVarP(&FlagInput, "input", "i",
		"", "Input file to colorize, defaults to stdin")

	nFlags.StringVarP(&FlagOutput, "output", "o",
		"", "Write colorized text to file instead of STDOUT. Use --stdout if output should go both to file and STDOUT.")

	for flagName, optName := range hideFlags {
		err = nFlags.MarkHidden(optName)
		if nil != err {
			xLog.Printf("could not mark option %s as %s hidden because %s\n",
				optName, flagName, err.Error())
			myFatal()
		}
	}

	// Fetch and load the program flags
	err = nFlags.Parse(os.Args[1:])
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", nFlags.FlagUsagesWrapped(75))
		xLog.Fatalf("\nerror parsing flags because: %s\n%s %s\n%s\n\t%v\n",
			err.Error(),
			"  common issue: 2 hyphens for long-form arguments,",
			"  1 hyphen for short-form argument",
			"  Program arguments are: ",
			os.Args)
	}

	if FlagClipboardBuffer {
		flagSet("nopaste", "false")
		flagSet("pipe", "false")
		flagSet("input", "")
	}

	if FlagPipe {
		flagSet("stdout", "true")
		flagSet("input", "")
		flagSet("output", "")
		flagSet("quiet", "true")
	}

	// do quietness
	// only write to logfile not stderr
	// for debug and verbose messages
	if FlagQuiet {
		xLog.SetOutput(xLogBuffer)
		// messages only to logfile, not stderr
	}

	if FlagClip && !modeClipboardAvailable {
		if FlagVerbose {
			xLog.Printf("This system does not offer a clipboard to paste to! --nopaste enabled!")
		}
		flagSet("nopaste", "false")
	}

	if FlagDebug && FlagVerbose {
		xLog.Println("\t\t/*** start program flags ***/\n")
		nFlags.VisitAll(logFlag)
		xLog.Println("\t\t/***   end program flags ***/")
	}

	if FlagHelp {
		var err1, err2 error
		_, thisCmd := filepath.Split(os.Args[0])
		_, err1 = fmt.Fprint(os.Stdout, "\n", "usage for ", thisCmd, ":\n")
		_, err2 = fmt.Fprintf(os.Stdout, "%s\n", nFlags.FlagUsagesWrapped(75))
		if nil != err1 || nil != err2 {
			xLog.Printf("huh? can't write to os.stdout because\n%s",
				misc.ConcatenateErrors(err1, err2).Error())
		}
		UsageMessage()
		_, _ = fmt.Fprintf(os.Stdout, "\t please see USAGE.MD for details")
		myFatal(0)
	}

	if FlagVerbose {
		errMsg := ""
		user, host, err := misc.UserHostInfo()
		if nil != err {
			errMsg = " (encountered error " + err.Error() + ")"
		}
		xLog.Printf("Verbose mode active (all debug and informative messages) for %s@%s%s",
			user, host, errMsg)
	}

	if FlagDebug && FlagVerbose {
		_, exeName := filepath.Split(os.Args[0])
		exeName = strings.TrimSuffix(exeName, filepath.Ext(exeName))
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			xLog.Printf("huh? Could not read build information for %s "+
				"-- perhaps compiled without module support?", exeName)
		} else {
			xLog.Printf("\n***** %s BuildInfo: *****\n%s\n%s\n",
				exeName, bi.String(), strings.Repeat("*", 22+len(exeName)))
		}
	}

	// Override the default TRUE setting for FlagStdout iff FlagStdout was not set by user
	if misc.IsStringSet(&FlagOutput) && !nFlags.Changed("stdout") {
		flagSet("stdout", "false")
	}

	if !misc.IsStringSet(&FlagOutput) && !FlagClip {
		if !FlagStdout {
			flagSet("stdout", "true")
		}
	}

}

// logFlag -- This writes out to the logger the value of a
// particular flag. Called indirectly. `Write()` is used
// directly to prevent wierd interactions with backslash
// in filenames
func logFlag(flag *pflag.Flag) {
	var sb strings.Builder
	sb.WriteString(" flag ")
	sb.WriteString(flag.Name)
	sb.WriteString(" has value [")
	sb.Write([]byte(flag.Value.String()))
	sb.WriteString("] with default [")
	sb.Write([]byte(flag.DefValue))
	sb.WriteString("]\n")
	_, _ = xLog.Writer().Write([]byte(sb.String()))
}

// UsageMessage prints useful information to the log
// Example usage:
//
//	UsageMessage()
func UsageMessage() {
	xLog.Printf("Useful Information Here")
}

// flagSet sets the value of a flag in the given flag set.
// If an error occurs while setting the flag value, it logs an error message
// using the xLog package and calls the myFatal function.
func flagSet(flag string, val string) {
	err := nFlags.Set(flag, val)
	if nil != err {
		xLog.Printf("huh? Could not set Flag [%s] to [%s] because %s",
			flag, val, err.Error())
		myFatal()
	}
}
