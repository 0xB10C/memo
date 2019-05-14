package logger

import (
	"log"
	"os"
	"time"
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
	traceHandle := os.Stdout
	infoHandle := os.Stdout
	warningHandle := os.Stdout
	errorHandle := os.Stderr

	logFlags := log.Ldate | log.Ltime | log.Lshortfile

	Trace = log.New(traceHandle, "TRACE: ", logFlags)
	Info = log.New(infoHandle, "INFO:  ", logFlags)
	Warning = log.New(warningHandle, "WARN:  ", logFlags)
	Error = log.New(errorHandle, "ERROR: ", logFlags)

	Trace.Println("Setup Logger")
}

func TrackTime(start time.Time, fname string) {
	elapsed := time.Since(start)
	Trace.Println(elapsed.String() + " \t for " + fname)
}
