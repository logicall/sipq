// sipq.go
package main

import (
	"os"
	"time"

	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/scenario"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/transport"
	"github.com/henryscala/sipq/util"
)

func main() {
	trace.Trace.Println("sipq started")
	defer trace.Trace.Println("sipq exited")

	var err error

	trace.Trace.Println("local ip", config.LocalIP, "local port", config.LocalPort, "transport type", config.TransportType)
	err = transport.StartServer(config.LocalIP, config.LocalPort, transport.TransportType(config.TransportType))
	if err != nil {
		util.ErrorPanic(err)
	}
	err = scenario.LoadFile(config.ScenarioFile)
	if err != nil {
		util.ErrorPanic(err)
	}

	var scenarioSuccess chan bool = make(chan bool)

	ascenario := scenario.New()
	go ascenario.Run(scenarioSuccess)

	scenarioTimeout := time.After(config.TimeLimit)
	select {
	case sucess := <-scenarioSuccess:
		if sucess {
			trace.Trace.Println("scenario succeed")
			os.Exit(0)
		} else {
			trace.Trace.Println("scenario failed")
			os.Exit(1)
		}
	case <-scenarioTimeout:
		trace.Trace.Println("scenario timeout")
		os.Exit(2)
	}

}
