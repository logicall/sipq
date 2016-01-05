package config

import (
	"flag"
	"strings"
)

var (
	ConfigFile        *string = flag.String("f", "", "the config file of sipq")
	ConfigFileExample *bool   = flag.Bool("cfg-example", false, "show a config file example")
	TransportType     *string = flag.String("t", "udp", "transport type to run the scenario")
	ScenarioFile      *string = flag.String("s", "", "the text file that contains description of scenarios")
)

func init() {
	flag.Parse()
}

func IsStreamTransport() bool {
	if TransportType == nil {
		return false
	}
	s := strings.ToLower(*TransportType)
	switch s {
	case "udp":
		return false
	default:
		return true
	}
}
