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
func handleNewData(conn *Connection) {
	var buf []byte = make([]byte, coding.MaxUdpPacketLen)

	for {
		switch conn.TransportType {
		case TCP:
			laddr := conn.Conn.LocalAddr()
			raddr := conn.Conn.RemoteAddr()
			msg, err := coding.FetchSipMessageFromReader(conn.Conn, true)
			if err != nil && err != io.EOF {
				conn.Close()
				return
			}
			msg.LocalAddr = laddr
			msg.RemoteAddr = raddr
			if err == io.EOF {
				conn.Close()
				return
			}
			sipMsgChan <- *msg

		//get a whole message from the connections and output the message
		case UDP:
			laddr := conn.Conn.LocalAddr()
			n, raddr, err := conn.ReadFrom(buf)

			if err != nil {
				trace.Trace.Println("read failed from udp", conn, err) //UDP connection, not return
				continue
			}
			udpReader := bytes.NewReader(buf[:n])
			msg, err := coding.FetchSipMessageFromReader(udpReader, false)
			if err != nil && err != io.EOF {
				trace.Trace.Fatalln("UDP server socket encounters unexpected error", err)
				return
			}

			msg.LocalAddr = laddr
			msg.RemoteAddr = raddr
			sipMsgChan <- *msg
		default:
			trace.Trace.Fatalln("not implemented")
		}
	}
}

//should be called in a go routine, since it is blocking
func handleNewConn(svr *Server) {
	for {
		conn, err := svr.Accept()
		util.ErrorPanic(err)
		trace.Trace.Println("new connection from",
			conn.Conn.RemoteAddr(), "come to", conn.Conn.LocalAddr())
		allConnections.Add(conn)
		go handleNewData(conn)
	}
}

func Send(msg *coding.SipMessage, transportType TransportType) error {
	switch transportType {
	case TCP:
		return sendTcp(msg)
	case UDP:
		return sendUdp(msg)
	}
	util.ErrorPanic(ErrTransport) //NOT IMPLEMENTED
	return ErrTransport

}

//This function is blocking
func sendTcp(msg *coding.SipMessage) error {

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
