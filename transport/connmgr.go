package transport

import (
	"bytes"
	"fmt"

	"io"

	"net"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/concurrent"
)

var (
	//store all connections of any transport type
	allConnections *concurrent.List = concurrent.NewList()

	//to communicate with comsumers of sip message
	sipMsgChan chan coding.SipMessage = make(chan coding.SipMessage)

	ErrTransport error = fmt.Errorf("some error happened on transport")
)

//for finding a established connection using local and remote address
func sameLocalRemoteAddrFunc(localAddr net.Addr, remoteAddr net.Addr) func(interface{}) bool {
	return func(conn interface{}) bool {
		c := conn.(*Connection)
		laddr := c.Conn.LocalAddr()
		raddr := c.Conn.RemoteAddr()
		if localAddr.String() == laddr.String() && remoteAddr.String() == raddr.String() {
			return true
		}
		return false
	}

}

//This function is blocking.
//Serves as interface toward transport users.
func FetchSipMessage() *coding.SipMessage {
	var msg coding.SipMessage
	msg = <-sipMsgChan
	return &msg
}

//This function is blocking.
//Serves as interface toward transport users.
func SendTcp(msg *coding.SipMessage) (int, error) {
	//find connection
	result, ok := allConnections.FindItemBy(sameLocalRemoteAddrFunc(msg.LocalAddr, msg.RemoteAddr))
	if !ok {
		return -1, ErrTransport
	}
	conn := result.(*Connection)

	//convert SipMessage for transfer on the wire
	//put the load on the wire
	return conn.Write([]byte(msg.String()))
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
