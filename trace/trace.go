package trace

import (
	"log"
	"os"
)

var Trace *log.Logger

func init() {
	Trace = log.New(os.Stdout, "sipq: ", log.LstdFlags|log.Lshortfile)

}
