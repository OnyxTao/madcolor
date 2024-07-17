package misc

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Jan 2 15:04:05 2006 MST
// DATE_OCPI time format for DateTime 2015-06-29T20:39:09
const DATE_OCPI = "2006-01-02T15:04:05"

var emptyString = ""

// SafeString returns either the pointer to the string,
// or a pointer to the emply string if the string is
// unset
func SafeString(test *string) (safe *string) {
	if IsStringSet(test) {
		return test
	}
	return &emptyString
}

// DeferError
// accounts for an at-close function that
// returns an error at its close
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
		_, _ = fmt.Fprintf(os.Stderr,
			"[%s] error in DeferError from file: %s line %d\n"+
				" error: %s\n\t(may be harmless!)",
			time.Now().UTC().Format(time.RFC822),
			file, line, err.Error())
	}
}

// IsStringSet -- returns true iff string is neither nil nor empty
func IsStringSet(s *string) (isSet bool) {
	if nil != s && "" != *s {
		return true
	}
	return false
}

// Ternary -- convert a true/false condition into the appropriate value
// or, Go, why did you take my ternary operator?
func Ternary(key bool, trueVal interface{}, falseVal interface{}) (val interface{}) {
	if key {
		return trueVal
	}
	return falseVal
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
