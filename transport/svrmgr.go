package transport

import (
	"net"

	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/container/concurrent"
)

type Server struct {
	TransportType Type
	Listener      net.Listener //used by TCP/TLS/SCTP instead of UDP
	UdpConn       *Connection  //used by UDP
}

var allServers *concurrent.List = concurrent.NewList()

//start the server and keep the server in all server list
func StartServer(ip string, port int, transportType Type) error {
	addr := util.AddrStr(ip, port)
	trace.Debug("starting server", transportType, addr)
	var server *Server
	var err error
	switch transportType {
	case TCP:
		server, err = createTcpServer(addr)
		go handleNewConn(server)

	case UDP:
		server, err = CreateUdpServer(addr)

	}

	if err != nil {
		return err
	}
	allServers.Add(server)
	return nil
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
func createTcpServer(laddr string) (*Server, error) {
	server := &Server{TransportType: TCP}

	listener, err := net.Listen(TCP.String(), laddr)
	if err != nil {
		return nil, err
	}
	server.Listener = listener

	return server, nil

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
