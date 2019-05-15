package logger

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/0xb10c/memo/memod/config"
)

var (
	// Trace logs trace messages
	Trace *log.Logger
	// Info logs info messages
	Info *log.Logger
	// Warning logs warning messages
	Warning *log.Logger
	// Error logs error messages
	Error *log.Logger
)

func init() {

	traceHandle := ioutil.Discard
	if config.GetBool("log.enableTrace") {
		traceHandle = os.Stdout
	}

	infoHandle := os.Stdout
	warningHandle := os.Stdout
	errorHandle := os.Stderr

	logFlags := log.Ldate | log.Ltime //log.Lshortfile

	Trace = log.New(traceHandle, Dim("TRACE: "), logFlags)
	Info = log.New(infoHandle, Blue("INFO: "), logFlags)
	Warning = log.New(warningHandle, Yellow("WARN: "), logFlags)
	Error = log.New(errorHandle, Red("ERROR: "), logFlags)

	Info.Println("Setup logger. Logging Trace:", config.GetBool("log.enableTrace"))
}

// TrackTime tracks the time a function takes till return and logs it to Trace
func TrackTime(start time.Time, funcname string) {
	elapsed := time.Since(start)
	Trace.Println(funcname + " took " + elapsed.String())
}
