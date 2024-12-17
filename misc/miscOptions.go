package misc

import (
	"fmt"
	"os"
)

/**************
 * internal flags to mirror the caller's settings
 * to provide better debug / verbose info
 **************/

var flagDebug = false
var flagVerbose = false
var logPrintf = defaultPrintf
var fatal = defaultFatal

func IsDebug() bool   { return flagDebug }
func IsVerbose() bool { return flagVerbose }

// OptionDebug sets the debug flag to the specified value and returns the old value of the debug flag.
func OptionDebug(debug bool) (old bool) {
	old = flagDebug
	flagDebug = debug
	return old
}

// OptionVerbose sets the verbose flag to the specified value, and returns the old value of the verbose flag.
func OptionVerbose(verbose bool) (old bool) {
	old = flagVerbose
	flagVerbose = verbose
	return old
}

// OptionPrintf allows setting a custom printf function for logging.
// It takes a function `f` with the same signature as the default `printf`
// function and returns the old `printf` function. Useful when there's a
// custom logger function compatible with printf to integrate external
// package messages with the program's logfile.
func OptionPrintf(f func(format string, a ...interface{}) (n int, err error)) (old func(format string, a ...interface{}) (n int, err error)) {
	old = logPrintf
	logPrintf = f
	return old
}

// defaultPrintf writes a formatted string to stderr. It returns the number of bytes
// written and any write error encountered. It is a default print function, almost
// always overwritten by either safeLogPrint or xLog.Printf ... but just in case ...
func defaultPrintf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// OptionFatal sets a custom exit routine to ensure
// cleanup and log closure happens as cleanly as possible
func OptionFatal(f func(retcode ...int)) (old func(retcode ...int)) {
	old = fatal
	fatal = f
	return old
}

func defaultFatal(retcode ...int) {
	rc := 0
	if len(retcode) > 0 {
		rc = retcode[0]
	}
	FinishClose()
	os.Exit(rc)
}
