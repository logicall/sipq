package config

import (
	"flag"
	"strings"
)

var (
	ConfigFile        string
	ConfigFileExample bool
	TransportType     string
	ScenarioFile      string
)

func init() {
	flag.StringVar(&ConfigFile, "f", "", "the config file of sipq")
	flag.BoolVar(&ConfigFileExample, "cfg-example", false, "show a config file example")
	flag.StringVar(&TransportType, "t", "udp", "transport type to run the scenario")
	flag.StringVar(&ScenarioFile, "s", "", "the text file that contains description of scenarios")

	flag.Parse()
}

func IsStreamTransport() bool {
	s := strings.ToLower(TransportType)
	switch s {
	case "udp":
		return false
	default:
		return true
	}
}
