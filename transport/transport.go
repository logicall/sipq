// transport.go
package transport

type Type string

const (
	TCP  Type = "tcp"
	UDP  Type = "udp"
	SCTP Type = "sctp"
	TLS  Type = "tls"
)

func (tt Type) String() string {
	return string(tt)
}
