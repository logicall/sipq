// sipq.go
package main

import (
	"fmt"

	"github.com/henryscala/sipq/config"
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
	transport.AllServers = transport.StartServers(config.TheExeConfig)

	//wait forever here
	var exit chan bool = make(chan bool)
	<-exit
}
