// sipq.go
package main

import (
	"flag"
	"fmt"

	"sipq/config"
	"sipq/trace"
	"sipq/transport"
	"sipq/util"
)

func main() {
	trace.Trace.Println("sipq started")
	defer trace.Trace.Println("sipq exited")

	var configFile *string = flag.String("f", "", "the config file of sipq")
	var configFileExample *bool = flag.Bool("cfg-example", false, "show a config file example")
	flag.Parse()

	if *configFileExample {
		fmt.Println(config.DefaultConfig)
		return
	}

	var err error
	if *configFile != "" {
		config.TheExeConfig, err = config.ReadExeConfigFile(*configFile)
		util.ErrorPanic(err)
	}
	transport.AllServers = transport.StartServers(config.TheExeConfig)

}
