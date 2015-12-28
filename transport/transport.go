// transport.go
package transport

import (
	"bufio"

	"net"
)

type TransportType string

const (
	ServerAddress string = "127.0.0.1:50600"
	ClientAddress string = "127.0.0.1:50700"
)

const (
	TransportTypeTCP  TransportType = "tcp"
	TransportTypeUDP  TransportType = "udp"
	TransportTypeSCTP TransportType = "sctp"
	TransportTypeTLS  TransportType = "tls"
)

func (tt TransportType) String() string {
	return string(tt)
}

type Server struct {
	TransportType TransportType
	Listener      net.Listener
}

type Connection struct {
	TransportType TransportType
	Conn          net.Conn //interface type
}

//used by UDP
func (conn *Connection) WriteTo(buf []byte, addr net.Addr) (int, error) {
	var udpConn *net.UDPConn = conn.Conn.(*net.UDPConn)
	n, err := udpConn.WriteTo(buf, addr)
	return n, err
}

//used by UDP
func (conn *Connection) ReadFrom(buf []byte) (int, net.Addr, error) {
	var udpConn *net.UDPConn = conn.Conn.(*net.UDPConn)
	n, addr, err := udpConn.ReadFrom(buf)
	return n, addr, err
}

//used by TCP
func (conn *Connection) Reader() *bufio.Reader {
	return bufio.NewReader(conn.Conn)
}

//used by TCP
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
	result := &Connection{Conn: conn, TransportType: TransportTypeTCP}
	return result, nil
}

//used for both client and server
func CreateUdpConnection(laddr string) (*Connection, error) {
	svrConn := &Connection{TransportType: TransportTypeUDP}
	var err error
	udpAddr, err := net.ResolveUDPAddr(TransportTypeUDP.String(), laddr)
	if err != nil {
		return nil, err
	}

	svrConn.Conn, err = net.ListenUDP(TransportTypeUDP.String(), udpAddr)

	if err != nil {
		return nil, err
	}

	return svrConn, nil
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
	result := &Connection{Conn: conn, TransportType: TransportTypeTCP}
	return result, nil
}
