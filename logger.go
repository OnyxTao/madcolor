package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var xLogFile *os.File
var xLogBuffer *bufio.Writer
var xLog log.Logger

// flushLog flushes the log buffer if it is not nil.
// If flushing the buffer results in an error, it logs the error message to standard output.
// This function is typically used to ensure that all log messages are written before shutting down the logging service.
func flushCloseLog() {
	if nil != xLogBuffer {
		err := xLogBuffer.Flush()
		if nil != err {
			_, _ = fmt.Fprintf(os.Stdout, "huh? could not flush xLogBuffer because %s", err.Error())
		}
	}
}

var closeLogMutex sync.Mutex

// closeLog shuts the logging service down
// cleanly, flushing buffers (and thus
// preserving the most likely error of
// interest)
func closeLog() {

	var err error = nil

	closeLogMutex.Lock()
	{
		if nil != xLogBuffer {
			flushCloseLog()
			xLogBuffer = nil
		}

		if nil != xLogFile {
			err = xLogFile.Close()
			xLogFile = nil
		}
	}
	closeLogMutex.Unlock()

	if nil != err {
		safeLogPrintf(err.Error())
	}
}

// initLog initializes the log file and log buffer.
// It opens the log file with the specified name, creating it if it does not exist, and truncates it if it does exist.
// If opening the log file encounters an error, it logs the error message to standard output using safeLogPrintf.
// It creates a new bufio.Writer to be used as the log buffer and sets the log writers to the standard output and the log buffer.
// It sets the log flags to include the date, time, UTC, and short file.
// It resolves the absolute path of the log file and logs it using safeLogPrintf.
// This function is typically called at the initialization of the logging service.
// The log file name should be passed as the lfName argument.
func initLog(lfName string) {
	var err error
	var logWriters = make([]io.Writer, 0, 2)
	xLogFile, err = os.OpenFile(lfName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		safeLogPrintf("error opening log file %s because %s",
			lfName, err.Error())
	}

	xLogBuffer = bufio.NewWriter(xLogFile)
	logWriters = append(logWriters, os.Stdout)
	logWriters = append(logWriters, xLogBuffer)
	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(io.MultiWriter(logWriters...))

	// logPath, err := filepath.Abs(xLogFile.Name())
	if nil != err {
		safeLogPrintf("huh? could not resolve logfilename %s because %s",
			xLogFile.Name(), err.Error())
		myFatal()
	}
	// safeLogPrintf("Logfile set to %s", logPath)

}

var myFatalMutex sync.Mutex

// myFatal is meant to close the program, and close the
// log files properly. Go doesn't support optional arguments,
// but variadic arguments allow finessing this. myFatal() gets
// a default RC of -1, and that's overridden by the first int
// in the slice of integers argument (which is present
// even if the length is 0).
//
// At some point, might create a more
// thorough at-close routine and register closing the file
// and log as part of the things to do 'at close'.
func myFatal(rcList ...int) {
	var rc int = -1
	myFatalMutex.Lock()
	// the app never releases this lock, because the
	// app is exiting for an fatal error. Only one
	// abnormal program termination at a time.
	// any threads waiting on this lock are similarly
	// in a fatal condition.

	// default rc is -1, but that *might* be
	// overridden by the caller
	if len(rcList) > 0 {
		rc = rcList[0]
	}

	// if this is an expected exit
	// this doesn't need to be logged
	if rc != 0 || FlagDebug {
		_, srcFile, srcLine, ok := runtime.Caller(1)
		if ok {
			srcFile = filepath.Base(srcFile)
			safeLogPrintf("\n\t\t/*** myFatal called ***/\n"+
				"\tfrom file:line %12s:%04d\n"+
				"\t\t/*** myFatal ended ***/", srcFile, srcLine)
		} else {
			safeLogPrintf("\n\t\t/*** myFatal called ***/\n" +
				"\tbut could not get stack information for caller\n" +
				"\t\t/*** myFatal ended ***/")
		}
	}
	closeLog()
	os.Exit(rc)
}

// safeLogPrintf may be called in lieu of xLog.Printf() if there
// is a possibility the log may not be open. If the log is
// available, well and good. Otherwise, print the message to
// STDERR.
var safeLogPrintfMutex sync.Mutex

func safeLogPrintf(format string, a ...any) {
	safeLogPrintfMutex.Lock()
	defer safeLogPrintfMutex.Unlock()
	if nil != xLogBuffer && nil != xLogFile {
		xLog.Printf(format, a...)
	} else {
		_, _ = fmt.Fprintf(os.Stderr,
			"\n\tSAFELOG\n"+format+"\n",
			a...)
	}
}

// debugMapStringString is a function that takes a map of string keys and string values as input.
// It generates a formatted string representation of the map for debugging purposes.
// The output string includes the size of the map and all key-value pairs in a tabular format.
// Each key-value pair is displayed in a separate line, with the key and value left-aligned in columns of width 20.
// This function does not return any value, it simply writes the debug string to standard output.
// Example usage:
//
//	params := map[string]string{
//	    "key1": "value1",
//	    "key2": "value2",
//	    "key3": "value3",
//	}
//	debugMapStringString(params)
//
// Output:
//
//	got map[string]string size 3
//	[ key1               ][ value1             ]
//	[ key2               ][ value2             ]
//	[ key3               ][ value3             ]
func debugMapStringString(params map[string]string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n\tgot map[string]string size %d\n", len(params)))
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("\t[ %-20s ][ %-20s ]\n", k, v))
	}
}

type NLVWriter struct {
	writer *bufio.Writer
}

func NewNLVWriter(b *bufio.Writer) (q *NLVWriter) {
	a := NLVWriter{b}
	return &a
}
func (w *NLVWriter) WriteString(s ...string) {
	for _, str := range s {
		_, err := w.writer.WriteString(str)
		if nil != err {
			xLog.Printf("failed to write %s to bufio.writer because %s",
				str, err.Error())
		}
	}
}
func (w *NLVWriter) Write(s ...[]byte) {
	for _, str := range s {
		_, err := w.writer.Write(str)
		if nil != err {
			xLog.Printf("failed to write %s to bufio.writer because %s",
				string(str), err.Error())
		}
	}
}
func (w *NLVWriter) WriteRune(s ...rune) {
	for _, str := range s {
		_, err := w.writer.WriteRune(str)
		if nil != err {
			xLog.Printf("failed to write %s to bufio.writer because %s",
				string(str), err.Error())
		}
	}
}
