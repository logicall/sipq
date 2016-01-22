// transport.go
package transport

type TransportType string

const (
	TCP  TransportType = "tcp"
	UDP  TransportType = "udp"
	SCTP TransportType = "sctp"
	TLS  TransportType = "tls"
)

func (tt TransportType) String() string {
	return string(tt)
}
