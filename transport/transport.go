// transport.go
package transport

import (
	"bufio"
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

type Server struct {
	TransportType TransportType
	Listener      net.Listener
}

type Connection struct {
	Conn net.Conn //interface type
}

func (conn *Connection) Reader() *bufio.Reader {
	return bufio.NewReader(conn.Conn)
}

func (conn *Connection) Writer() *bufio.Writer {
	return bufio.NewWriter(conn.Conn)
}

func (conn *Connection) Close() {
	conn.Conn.Close() //ignore error
}

func (svr *Server) String() string {
	return svr.TransportType.String()
}

func (svr *Server) Close() {
	svr.Listener.Close() //ignore error
}

func (svr *Server) Accept() (*Connection, error) {
	conn, err := svr.Listener.Accept()
	if err != nil {
		return nil, err
	}
	result := &Connection{Conn: conn}
	return result, nil
}

//laddr(local addr) is like 127.0.0.1:5060
func CreateTcpServer(laddr string) (*Server, error) {
	server := &Server{TransportType: TransportTypeTCP}

	listener, err := net.Listen(TransportTypeTCP.String(), laddr)
	if err != nil {
		return nil, err
	}
	server.Listener = listener

	return server, nil

}

//raddr(remote addr) is like 127.0.0.1:5060
func CreateTcpConnection(raddr string) (*Connection, error) {
	conn, err := net.Dial(TransportTypeTCP.String(), raddr)
	if err != nil {
		return nil, err
	}
	result := &Connection{conn}
	return result, nil
}
