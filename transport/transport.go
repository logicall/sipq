// transport.go
package transport

import (
	"strings"

	"github.com/henryscala/sipq/trace"
)

type Type int

const (
	TCP     Type = iota
	UDP     Type = iota
	TLS     Type = iota
	SCTP    Type = iota
	INVALID Type = iota
)

func (t Type) String() string {
	switch t {
	case TCP:
		return "tcp"
	case UDP:
		return "udp"
	case SCTP:
		return "sctp"
	case TLS:
		return "tls"
	}
	trace.Error("not implemented")
	return "not implemented"
}

func TypeFromString(transportType string) Type {
	transportType = strings.ToLower(transportType)
	switch {
	case strings.Contains(transportType, "tcp"):
		return TCP
	case strings.Contains(transportType, "udp"):
		return UDP
	case strings.Contains(transportType, "sctp"):
		return SCTP
	case strings.Contains(transportType, "tls"):
		return TLS
	}
	trace.Error("not implemented")
	return INVALID
}
