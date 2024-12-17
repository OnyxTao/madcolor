package misc

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"reflect"
	"runtime"
	"strings"
)

// DATE_OCPI time format for DateTime 2015-06-29T20:39:09
// Jan 2 15:04:05 2006 MST
const DATE_OCPI = "2006-01-02T15:04:05"

var emptyString = ""

// SafeString returns either the pointer to the string,
// or a pointer to the empty string if the string is
// unset
func SafeString(test *string) (safe *string) {
	if IsStringSet(test) {
		return test
	}
	return &emptyString
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

// GetThingName returns the name of the function passed as an interface.
func GetThingName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func RecordString(outFileName string, inTx <-chan string, wgDone func()) {
	defer wgDone()

	bout, err := os.OpenFile(outFileName,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if nil != err {
		logPrintf("Failed to open %s because %s\n",
			outFileName, err.Error())
		fatal()
	}
	defer DeferError(bout.Close)

	bw := bufio.NewWriterSize(bout, 1024*8)
	defer DeferError(bw.Flush)

	for val := range inTx {
		_, err = bw.WriteString(val)
		if nil != err {
			logPrintf("failed to write string %s to file %s because %s\n",
				val, outFileName, err.Error())
			fatal()
		}
		_, err = bw.WriteRune('\n')
		if nil != err {
			logPrintf("failed to write newline following string %s to file %s because %s\n",
				val, outFileName, err.Error())
			fatal()
		}
	}

	if flagDebug {
		logPrintf("Finished output to file %s\n", outFileName)
	}

	/***** deferred actions
	 * flush buffered writer
	 * close output file
	 * signal waitgroup.Done
	 ******/
}

// LogPrintf logs a formatted message using the provided format
// string and optional arguments. The
// expectation is that a program with a custom logging method may
// register that function with the misc package, and then it may
// It internally delegates to the logPrintf function for output handling.
func LogPrintf(format string, a ...interface{}) {
	_, _ = logPrintf(format, a...)
}

// Fatal terminates the application with an optional return code,
// defaulting to a standard fatal behavior if unspecified. The
// expectation is that a program with a custom fatal() method may
// register that function with the misc package, and then it may
// be called from other packages.
func Fatal(retcode ...int) {
	fatal(retcode...)
}
