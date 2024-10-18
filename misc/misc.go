package misc

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Jan 2 15:04:05 2006 MST
// DATE_OCPI time format for DateTime 2015-06-29T20:39:09
const DATE_OCPI = "2006-01-02T15:04:05"

var emptyString = ""

var mDebug = false
var mVerbose = false
var xLog *log.Logger = nil
var mFatal func(...int) = miscExit

func reportError(msg ...string) {
	for _, m := range msg {
		n := SafeString(&m)
		if nil != xLog {
			xLog.Println(n)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, n)
		}
	}
}

// miscExit terminates the current program with the
// given exit code, or 0 if no code is provided. It's
// a default, in case it's not overridden by MyFatal()
// or something similar from the main package. Must
// call indirectly because the function takes an
// array of exit values to permit a naked call from
// MyFatal().
func miscExit(code ...int) {
	var rc int = -4
	if len(code) <= 0 {
		rc = code[0]
	}
	os.Exit(rc)
}

// SetOptions configures the debug mode, verbose mode, and logger to
// primary logging tool, as well as the exit function. These options
// are NOT required to be set; there are reasonable default:
// debug == false
// verbose == false
// logging goes to os.Stderr
// fatal function goes to os.Exit(-4).
func SetOptions(debug bool, verbose bool, mainLogger *log.Logger, fatal func(...int)) {
	mDebug = debug
	mVerbose = verbose
	// do not point to nil functions :-)
	if nil != mainLogger {
		xLog = mainLogger
	}
	if nil != fatal {
		mFatal = fatal
	}
}

// SafeString returns either the pointer to the string,
// or a pointer to the empty string if the string is
// unset
func SafeString(test *string) (safe *string) {
	if IsStringSet(test) {
		return test
	}
	return &emptyString
}

// DeferError accounts for an at-close function that
// returns an error function, so it cannot be simply
// deferred
func DeferError(f func() error) {
	err := f()
	if nil != err {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		} else {
			file = filepath.Base(file)
		}
		msg := fmt.Sprintf("[%s] error in DeferError from file: %s line %d\n"+
			" error: %s\n\t(may be harmless!)",
			time.Now().UTC().Format(time.RFC822),
			file, line, err.Error())
		reportError(msg)
	}
}

// IsStringSet -- returns true iff string is neither nil nor empty
func IsStringSet(s *string) (isSet bool) {
	if nil != s && "" != *s {
		return true
	}
	return false
}

// UserHostInfo returns the current username, current hostname and an error, as appropriate
func UserHostInfo() (userName string, hostName string, err error) {
	var ui *user.User
	ui, err = user.Current()
	if nil != err {
		return "",
			"",
			errors.New(fmt.Sprintf("UserHostInfo failed to get user.Current() because %s",
				err.Error()))
	}
	hostName, err = os.Hostname()
	if nil != err {
		return ui.Name, "",
			errors.New(fmt.Sprintf("UserHostInfo failed to get os.Hostname() because %s",
				err.Error()))
	}
	return ui.Name, hostName, nil
}

// ConcatenateErrors concatenates a variable number of error
// instances into a single error message, enumerating each error.
func ConcatenateErrors(errList ...error) error {
	if nil == errList {
		return nil
	}
	var sb strings.Builder
	ix := 1
	for _, err := range errList {
		if err == nil {
			continue
		}
		sb.WriteString(fmt.Sprintf("\n%02d.\t%s", ix, err.Error()))
		ix++
	}
	if sb.Len() > 0 {
		return errors.New(sb.String())
	}
	return nil
}

// RecordString writes strings received from the inTx channel to a specified file in outDir with outFileName.
// When the function completes, it calls the provided wgDone function. Each string is written followed by a newline.
// If an error occurs during file operations, the function panics. The function is intended to run in a
// goroutine, and is closed by closing the `inTx` channel.
func RecordString(outDir string, outFileName string, inTx <-chan string, wgDone func()) {
	defer wgDone()

	bout, err := os.OpenFile(path.Join(outDir, outFileName),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		msg := fmt.Sprintf("Failed to open %s because %s\n",
			path.Join(outDir, outFileName), err.Error())
		reportError(msg)
		mFatal(-3)
	}
	defer DeferError(bout.Close)

	bw := bufio.NewWriterSize(bout, 1024*4)
	defer DeferError(bw.Flush)

	for val := range inTx {
		_, err = bw.WriteString(val)
		if nil != err {
			msg := fmt.Sprintf("failed to write string %s to file %s because %s\n",
				val, outFileName, err.Error())
			reportError(msg)
			mFatal(-3)
		}
		err = bw.WriteByte('\n')
		if nil != err {
			msg := fmt.Sprintf("failed to write newline following string %s to file %s because %s\n",
				val, outFileName, err.Error())
			reportError(msg)
			mFatal(-3)
		}
	}

	/***********
	 * DEFERRED ACTIONS
	 * flush buffer_out
	 * close file
	 * call wgDone (waitGroup done)
	 **********/

}
