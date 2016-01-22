package trace

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	//Specify the present trace level
	//By default only print error level.
	TraceLevel int       = 0
	LogWriter  io.Writer = os.Stdout
	pid        string    = fmt.Sprint(os.Getpid())
	prefix     string    = "sipq"
)

//Note, log is different from normal program output.
//log may be turned off, but program output may not.
var traceLevels []string = []string{
	"error",   //log error message and exit the program
	"warning", //log warning message and user shall be cautious
	"info",    //log overall figures of the program(e.g. number of concurrent connections). It should not per transaction basis.
	"trace",   // only used to log the enter and exit of a function
	"debug",   // usd to log any kinf of debug print
}

func getFileNameLineNum() (fileName string, lineNo string) {
	_, fileName, line, _ := runtime.Caller(3)
	return fileName, fmt.Sprint(line)
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
		fileName, lineNo := getFileNameLineNum()
		fileName = filepath.Base(fileName)
		fileNameLineNo := fileName + ":" + lineNo

		var argList []interface{} = []interface{}{
			prefix, pid, traceLevels[selfLogLevel], time.Now().Format("15:04:05.000"), fileNameLineNo, "|",
		}
		argList = append(argList, args...)

		fmt.Fprintln(LogWriter, argList...)

		if selfLogLevel <= 0 {
			panic("critical error")
		}
	}
}

func init() {
	// In future, may init LogWriter to a file

}
