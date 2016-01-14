package config

import (
	"flag"
	"strings"
)

const (
	defaultIP   = "127.0.0.1"
	defaultPort = 5060
)

var (
	ConfigFile        string
	ConfigFileExample bool
	TransportType     string
	ScenarioFile      string
	RemoteIP          string
	RemotePort        int
	LocalIP           string
	LocalPort         int
)

func IsStreamTransport() bool {
	s := strings.ToLower(TransportType)
	switch s {
	case "udp":
		return false
	default:
		return true
	}
}

func init() {
	flag.StringVar(&ConfigFile, "config-file", "", "the config file of sipq")
	flag.BoolVar(&ConfigFileExample, "config-file-example", false, "show a config file example")
	flag.StringVar(&TransportType, "transport-type", "udp", "transport type to run the scenario")
	flag.StringVar(&ScenarioFile, "scenario-file", "", "the text file that contains description of scenarios")
	flag.StringVar(&RemoteIP, "remote-ip", defaultIP, "the IP address of the peer side")
	flag.IntVar(&RemotePort, "remote-port", defaultPort, "the IP port of the peer side")
	flag.StringVar(&RemoteIP, "local-ip", defaultIP, "the local IP address ")
	flag.IntVar(&RemotePort, "local-port", defaultPort, "the local IP port ")

	flag.Parse()
}
