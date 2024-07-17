package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"madcolor/misc"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

// wordSepNormalizeFunc all options are lowercase, so
// ... lowercase they shall be
func wordSepNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(strings.ToLower(name))
}

var nFlags *pflag.FlagSet

/* secret flags */
var FlagSlow bool
var FlagMaxBrightness int
var FlagMinBrightness int

/* standard flags */

var FlagHelp bool
var FlagQuiet bool
var FlagVerbose bool
var FlagDebug bool

/* program flags */
var FlagText string
var FlagInventColor bool
var FlagAntiColor bool

// initFlags initializes the command line flags for the program.
// It sets up the flag set, defines the flags, and parses the command line arguments.
func initFlags() {
	var err error

	hideFlags := make(map[string]string, 8)

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.SetNormalizeFunc(wordSepNormalizeFunc)

	// secret flags
	nFlags.IntVarP(&FlagMaxBrightness, "max", "", 160,
		"maximum total brightness of any foreground color")

	nFlags.IntVarP(&FlagMaxBrightness, "min", "", 0,
		"maximum total brightness of any foreground color")

	nFlags.BoolVarP(&FlagSlow, "slow", "", false,
		"Add some time between http calls (do not hammer server)")

	nFlags.BoolVarP(&FlagDebug, "debug", "d",
		false, "Enable additional informational and operational logging output for debug purposes")

	nFlags.BoolVarP(&FlagVerbose, "verbose", "v",
		false, "Supply additional run messages; use --debug for more information")

	nFlags.BoolVarP(&FlagHelp, "help", "h",
		false, "Display help message and usage information")

	nFlags.BoolVarP(&FlagQuiet, "quiet", "q",
		false, "Suppress log output to stdout and stderr (output still goes to logfile)")

	// program flags
	nFlags.BoolVarP(&FlagAntiColor, "anti", "a", false,
		"Set the colorspace background to the foreground complement")

	nFlags.BoolVarP(&FlagInventColor, "invent", "i", false,
		"randomly generate colors (rather than randomly select websafe colors)")

	nFlags.StringVarP(&FlagText, "text", "t",
		"", "Text to colorize")

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

	// do quietness setup first
	// only write to logfile not stderr
	// for debug and verbose messages
	if FlagQuiet {
		xLog.SetOutput(xLogBuffer)
		// messages only to logfile, not stderr
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
		_, _ = fmt.Fprintf(os.Stdout, "\t please see USAGE.MD for ")
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
