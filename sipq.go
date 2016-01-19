// sipq.go
package main

import (
	"fmt"
	"os"

	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/scenario"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/transport"
	"github.com/henryscala/sipq/util"
)

func main() {
	trace.Trace.Println("sipq started")
	defer trace.Trace.Println("sipq exited")

	if config.ConfigFileExample {
		fmt.Println(config.DefaultConfig)
		return
	}

	var err error
	if config.ConfigFile != "" {
		config.TheExeConfig, err = config.ReadExeConfigFile(config.ConfigFile)
		util.ErrorPanic(err)
	}

	trace.Trace.Println("local ip", config.LocalIP, "local port", config.LocalPort)
	transport.StartServer(config.LocalIP, config.LocalPort, transport.TransportType(config.TransportType))

	err = scenario.LoadFile(config.ScenarioFile)
	if err != nil {
		util.ErrorPanic(err)
	}

	var scenarioSuccess chan bool = make(chan bool)

	ascenario := scenario.New()
	go ascenario.Run(scenarioSuccess)

	sucess := <-scenarioSuccess
	if sucess {
		trace.Trace.Println("succeed")
		os.Exit(0)
	} else {
		trace.Trace.Println("failed")
		os.Exit(-1)
	}

}
