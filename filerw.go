package main

import (
	"bufio"
	"madcolor/misc"
	"os"
	"path"
)

// recordString writes a string to a file specified by outFileName
// inTx is a channel that receives strings to be written to the file
// wg is a pointer to a sync.WaitGroup that is used to signal when the function is done
func recordString(outFileName string, inTx <-chan string, wgDone func()) {
	defer wgDone()

	bout, err := os.OpenFile(path.Join(FlagOutputDir, outFileName),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		xLog.Printf("Failed to open %s because %s\n",
			path.Join(FlagOutputDir, outFileName), err.Error())
		myFatal()
	}
	defer misc.DeferError(bout.Close)

	bw := bufio.NewWriterSize(bout, 1024*8)
	defer misc.DeferError(bw.Flush)

	for val := range inTx {
		_, err = bw.WriteString(val)
		if nil != err {
			xLog.Printf("failed to write string %s to file %s because %s\n",
				val, outFileName, err.Error())
			myFatal()
		}
		_, err = bw.WriteRune('\n')
		if nil != err {
			xLog.Printf("failed to write newline following string %s to file %s because %s\n",
				val, outFileName, err.Error())
			myFatal()
		}
	}

	if FlagDebug {
		xLog.Printf("Finished output to file %s\n", outFileName)
	}

	/***** deferred actions
	 * flush buffered writer
	 * close output file
	 * signal waitgroup.Done
	 ******/
}
