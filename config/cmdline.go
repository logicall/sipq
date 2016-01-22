package config

import (
	"flag"
	"strings"
	"time"

	"github.com/henryscala/sipq/trace"
)

const (
	defaultIP   = "127.0.0.1"
	defaultPort = 5060
)

var (
	TransportType string
	ScenarioFile  string
	RemoteIP      string
	RemotePort    int
	LocalIP       string
	LocalPort     int
	TimeLimit     time.Duration
	TraceLevel    int
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
	flag.StringVar(&TransportType, "transport-type", "udp", "transport type to run the scenario")
	flag.StringVar(&ScenarioFile, "scenario-file", "", "the text file that contains description of scenarios")
	flag.StringVar(&RemoteIP, "remote-ip", defaultIP, "the IP address of the peer side")
	flag.IntVar(&RemotePort, "remote-port", defaultPort, "the IP port of the peer side")
	flag.StringVar(&LocalIP, "local-ip", defaultIP, "the local IP address ")
	flag.IntVar(&LocalPort, "local-port", defaultPort, "the local IP port ")
	flag.IntVar(&TraceLevel, "trace-level", 0, "0 error 1 warning 2 info 3 trace 4 debug 5 all -1 OFF")
	flag.DurationVar(&TimeLimit, "time-limit", 25*time.Second, "the maximum time the scenario can last. exceeding this limit the case fails")

	flag.Parse()

	trace.TraceLevel = TraceLevel
	trace.Debug("flag parse")
}
