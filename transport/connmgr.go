package transport

import (
	"bytes"
	"fmt"

	"io"

	"net"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/container/concurrent"
)

type Connection struct {
	TransportType Type
	Conn          net.Conn //interface type
}

var (
	//store all connections of any transport type
	allConnections *concurrent.List = concurrent.NewList()

	//to communicate with comsumers of sip message
	sipMsgChan chan coding.SipMessage = make(chan coding.SipMessage)

	ErrTransport error = fmt.Errorf("some error happened on transport")
)

//This function is blocking.
//Serves as interface toward transport users.
func FetchSipMessage() *coding.SipMessage {
	var msg coding.SipMessage
	msg = <-sipMsgChan
	return &msg
}

//should be called in a go routine, since it is blocking
// It seems that tcp conn may also meet EOF even though connection is not closed.(need verify further)
func handleNewData(conn *Connection) {
	trace.Trace("enter handleNewData", conn)
	defer trace.Trace("exit handleNewData", conn)
	var buf []byte = make([]byte, coding.MaxUdpPacketLen)

	for {
		switch conn.TransportType {
		case TCP:
			laddr := conn.Conn.LocalAddr()
			raddr := conn.Conn.RemoteAddr()
			trace.Debug("handleNewData", "laddr", laddr, "raddr", raddr)
			msg, err := coding.FetchSipMessageFromReader(conn.Conn, true)
			if err != nil {
				if err != io.EOF {
					//conn.Close() TODO: consider whether this statement is required
					return
				}
			}
			//this connection is over
			if msg == nil {
				//conn.Close() TODO: consider whether this statement is required
				return
			}
			msg.LocalAddr = laddr
			msg.RemoteAddr = raddr
			sipMsgChan <- *msg
			//this connection is over
			if err == io.EOF {
				//conn.Close() TODO: consider whether this statement is required
				return
			}

		//get a whole message from the connections and output the message
		case UDP:
			laddr := conn.Conn.LocalAddr()
			n, raddr, err := conn.ReadFrom(buf)

			if err != nil {
				trace.Debug("read failed from udp", conn, err) //UDP connection, not return
				continue
			}
			udpReader := bytes.NewReader(buf[:n])
			msg, err := coding.FetchSipMessageFromReader(udpReader, false)
			if err != nil {
				if err != io.EOF {
					trace.Error("UDP server socket encounters unexpected error", err)
					return
				}
			}

			//this connection is over
			if msg == nil {
				return
			}

			msg.LocalAddr = laddr
			msg.RemoteAddr = raddr
			sipMsgChan <- *msg
			//this connection is over
			if err == io.EOF {
				return
			}
		default:
			trace.Error("not implemented")
		}
	}
}

//should be called in a go routine, since it is blocking
func handleNewConn(svr *Server) {
	for {
		conn, err := svr.Accept()
		util.ErrorPanic(err)
		trace.Debug("new connection from",
			conn.Conn.RemoteAddr(), "come to", conn.Conn.LocalAddr())
		allConnections.Add(conn)
		go handleNewData(conn)
	}
}

func Send(msg *coding.SipMessage, transportType Type) error {
	trace.Trace("enter Send")
	defer trace.Trace("exit Send")
	switch transportType {
	case TCP:
		return sendTcp(msg)
	case UDP:
		return sendUdp(msg)
	default:

		return ErrTransport
	}

}

//This function is blocking
func sendTcp(msg *coding.SipMessage) error {
	trace.Trace("enter sendTcp")
	defer trace.Trace("exit sendTcp")
	var conn *Connection
	var err error
	//find connection
	result, ok := allConnections.FindItemBy(sameRemoteAddrFunc(msg.RemoteAddr))
	if !ok {
		conn, err = CreateTcpConnection(msg.RemoteAddr.String())
		if err != nil {
			return err
		}

	} else {
		conn = result.(*Connection)
	}

	//convert SipMessage for transfer on the wire
	//put the load on the wire
	sent, err := conn.Write([]byte(msg.String()))
	if err != nil {
		return err
	}
	// temporary assume the message is sent successfully by one attempt
	// in future may use a loop
	if sent != len(msg.String()) {
		return ErrTransport
	}
	return nil
}

func sendUdp(msg *coding.SipMessage) error {
	trace.Trace("enter sendUdp")
	defer trace.Trace("exit sendUdp")
	var conn *Connection
	var err error
	//find connection
	result, ok := allConnections.FindItemBy(sameLocalAddrFunc(msg.LocalAddr))
	if !ok {
		return ErrTransport

	}
	conn = result.(*Connection)

	//convert SipMessage for transfer on the wire
	//put the load on the wire
	sent, err := conn.WriteTo([]byte(msg.String()), msg.RemoteAddr)

	if err != nil {
		return err
	}
	// temporary assume the message is sent successfully by one attempt
	if sent != len(msg.String()) {
		return ErrTransport
	}
	return nil
}

//for finding a established connection using only remote address
func sameRemoteAddrFunc(remoteAddr net.Addr) func(interface{}) bool {
	return func(conn interface{}) bool {
		c := conn.(*Connection)

		raddr := c.Conn.RemoteAddr()
		if remoteAddr.String() == raddr.String() {
			return true
		}
		return false
	}

}

//for finding a udp server(also represented by connection)  using only local address
func sameLocalAddrFunc(localAddr net.Addr) func(interface{}) bool {
	return func(conn interface{}) bool {
		c := conn.(*Connection)

		laddr := c.Conn.LocalAddr()
		if localAddr.String() == laddr.String() {
			return true
		}
		return false
	}

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
	trace.Trace("enter Close", conn)
	defer trace.Trace("exit Close", conn)
	conn.Conn.Close() //ignore error
	allConnections.RemoveItem(conn)
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

//raddr(remote addr) is like 127.0.0.1:5060
func CreateTcpConnection(raddr string) (*Connection, error) {
	trace.Trace("enter CreateTcpConnection")
	defer trace.Trace("exit CreateTcpConnection")
	conn, err := net.Dial(TCP.String(), raddr)
	if err != nil {
		return nil, err
	}
	result := &Connection{Conn: conn, TransportType: TCP}

	allConnections.Add(result)
	go handleNewData(result)
	return result, nil
}
