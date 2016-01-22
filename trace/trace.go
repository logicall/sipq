package trace

import (
	"io"
	"log"
	"os"
)

//Specify the present trace level
//By default only print error level.
var TraceLevel int = 0
var LogWriter io.Writer = os.Stdout

//Note, log is different from normal program output.
//log may be turned off, but program output may not.
var TraceLevels []string = []string{
	"error",   //log error message and exit the program
	"warning", //log warning message and user shall be cautious
	"info",    //log overall figures of the program(e.g. number of concurrent connections). It should not per transaction basis.
	"trace",   // only used to log the enter and exit of a function
	"debug",   // usd to log any kinf of debug print
}

var loggers []*log.Logger = []*log.Logger{
	log.New(LogWriter, "sipq: error", log.LstdFlags|log.Lshortfile),
	log.New(LogWriter, "sipq: warning", log.LstdFlags|log.Lshortfile),
	log.New(LogWriter, "sipq: info", log.LstdFlags|log.Lshortfile),
	log.New(LogWriter, "sipq: trace", log.LstdFlags|log.Lshortfile),
	log.New(LogWriter, "sipq: debug", log.LstdFlags|log.Lshortfile),
}

func Error(args ...interface{}) {
	printLog(0, TraceLevel, args...)
}
func Warning(args ...interface{}) {
	printLog(1, TraceLevel, args...)
}
func Info(args ...interface{}) {
	printLog(2, TraceLevel, args...)
}
func Trace(args ...interface{}) {
	printLog(3, TraceLevel, args...)
}
func Debug(args ...interface{}) {
	printLog(4, TraceLevel, args...)
}

func printLog(selfLogLevel, globalLogLevel int, args ...interface{}) {
	if selfLogLevel <= globalLogLevel {
		loggers[selfLogLevel].Println(args...)
		if selfLogLevel <= 0 {
			panic("critical error")
		}
	}
}

func init() {
	// In future, may init LogWriter to a file

}
