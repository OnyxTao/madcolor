package misc

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const CLOSE_BUFFER_SIZE = 4

var atCloseMutex sync.Mutex
var atClose []func() error
var atCloseName []string

// init initializes the atClose and atCloseName
// slices with a predefined capacity of CLOSE_BUFFER_SIZE.
func init() {
	atClose = make([]func() error, 0, CLOSE_BUFFER_SIZE)
	atCloseName = make([]string, 0, CLOSE_BUFFER_SIZE)
}

// AtCloseErr registers a function that returns an error
// to be called when the application is closing.
func AtCloseErr(f func() error) {
	atCloseMutex.Lock()
	defer atCloseMutex.Unlock()
	atClose = append(atClose, f)
	atCloseName = append(atCloseName, GetThingName(f))
}

// AtClose registers a function to be called upon program termination.
// The function is appended to the atClose slice.
func AtClose(f func()) {
	atCloseMutex.Lock()
	defer atCloseMutex.Unlock()
	atClose = append(atClose, func() error { f(); return nil })
	atCloseName = append(atCloseName, GetThingName(f))
}

// FinishClose runs all functions in the atClose slice
// in reverse order. Logs function names and errors if
// flagDebug or flagVerbose is set.
func FinishClose() {
	atCloseMutex.Lock()
	defer atCloseMutex.Unlock()
	var err error
	if flagDebug {
		_, _ = logPrintf("Number of AtClose/AtCloseErr functions is %d (started with capacity %d)\n",
			len(atClose), CLOSE_BUFFER_SIZE)
	}
	for ix := len(atClose) - 1; ix >= 0; ix-- {
		err = (atClose[ix])()
		/* if flagDebug || flagVerbose {
			_, _ = printf("AtClose running function %s\n",
				atCloseName[ix])
		} */
		if nil != err {
			_, _ = logPrintf("AtClose function %s failed because %s\n",
				atCloseName[ix], err.Error())
		}
	}
}

// HandleSignal waits for an OS signal from the signalChan channel and
// acts upon receiving the signal, exiting the program. This allows for
// registered at-close routines to execute even if the program is killed.
func HandleSignal(signalChan <-chan os.Signal) {
	sig := <-signalChan
	_, _ = logPrintf("Got signal %v, exiting immediately\n", sig)
	fatal(-2)
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
		_, _ = logPrintf("[%s] error in DeferError from file: %s line %d\n"+
			" error: %s\n\t(may be harmless!)",
			time.Now().UTC().Format(time.RFC822),
			file, line, err.Error())
	}
}
