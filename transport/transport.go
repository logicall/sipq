// transport.go
package transport

import (
	"net"
)

type TransportType string

const (
	ServerAddress string = "127.0.0.1:50600"
	ClientAddress string = "127.0.0.1:50700"
)

const (
	TCP  TransportType = "tcp"
	UDP  TransportType = "udp"
	SCTP TransportType = "sctp"
	TLS  TransportType = "tls"
)

func (tt TransportType) String() string {
	return string(tt)
}

type Server struct {
	TransportType TransportType
	Listener      net.Listener //used by TCP/TLS/SCTP instead of UDP
	UdpConn       *Connection  //used by UDP
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
func (conn *Connection) Write(buf []byte) (int, error) {
	var tcpConn *net.TCPConn = conn.Conn.(*net.TCPConn)
	n, err := tcpConn.Write(buf)
	return n, err
}

func (conn *Connection) Close() {
	conn.Conn.Close() //ignore error
	allConnections.RemoveItem(conn)
}

func (svr *Server) String() string {
	return svr.TransportType.String()
}

func (svr *Server) Close() {
	//TODO remove all connections belong to this server
	if svr.Listener != nil {
		svr.Listener.Close() //ignore error
		svr.Listener = nil
	}

	if svr.UdpConn != nil {
		svr.UdpConn.Close()
		svr.UdpConn = nil
	}

}

func (svr *Server) Accept() (*Connection, error) {
	conn, err := svr.Listener.Accept()
	if err != nil {
		return nil, err
	}
	result := &Connection{Conn: conn, TransportType: TCP}
	return result, nil
}

func createUdpConnection(laddr string) (*Connection, error) {
	svrConn := &Connection{TransportType: UDP}
	var err error
	udpAddr, err := net.ResolveUDPAddr(UDP.String(), laddr)
	if err != nil {
		return nil, err
	}

	svrConn.Conn, err = net.ListenUDP(UDP.String(), udpAddr)

	if err != nil {
		return nil, err
	}

	allConnections.Add(svrConn)
	go handleNewData(svrConn)

	return svrConn, nil
}

//used for both client and server
func CreateUdpServer(laddr string) (*Server, error) {
	conn, err := createUdpConnection(laddr)
	if err != nil {
		return nil, err
	}
	server := &Server{TransportType: UDP, UdpConn: conn}
	return server, nil
}

//laddr(local addr) is like 127.0.0.1:5060
func CreateTcpServer(laddr string) (*Server, error) {
	server := &Server{TransportType: TCP}

	listener, err := net.Listen(TCP.String(), laddr)
	if err != nil {
		return nil, err
	}
	server.Listener = listener

	return server, nil

}

//raddr(remote addr) is like 127.0.0.1:5060
func CreateTcpConnection(raddr string) (*Connection, error) {
	conn, err := net.Dial(TCP.String(), raddr)
	if err != nil {
		return nil, err
	}
	result := &Connection{Conn: conn, TransportType: TCP}

	allConnections.Add(result)
	go handleNewData(result)
	return result, nil
}
