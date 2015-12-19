// transport.go
package transport

import (
	_ "errors"
	"net"
)

type TransportType string

const (
	TransportTypeTCP  TransportType = "tcp"
	TransportTypeUDP                = "udp"
	TransportTypeSCTP               = "sctp"
	TransportTypeTLS                = "tls"
)

func (tt TransportType) String() string {
	return string(tt)
}

type ListenServer struct {
	TransportType TransportType
	Listener      net.Listener
}

func (ls *ListenServer) String() string {
	return ls.TransportType.String()
}

func createTcpServer(laddr string) (*ListenServer, error) {
	server := &ListenServer{TransportType: TransportTypeTCP}

	listener, err := net.Listen(string(TransportTypeTCP), laddr)
	if err != nil {
		return nil, err
	}
	server.Listener = listener

	return server, nil

}
