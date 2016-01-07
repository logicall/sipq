package transport

import (
	"bytes"

	"io"

	"net"
	"sync"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
)

type Connections struct {
	connsLock sync.Mutex
	conns     []*Connection
}

func (conns *Connections) Remove(conn *Connection) {
	conns.connsLock.Lock()
	defer conns.connsLock.Unlock()

	var i int
	for i = 0; i < len(conns.conns); i++ {
		if conns.conns[i] == conn {
			break
		}
	}
	if i >= len(conns.conns) {
		trace.Trace.Fatalln(conn, "does not exist")
	}
	conns.conns = append(conns.conns[0:i], conns.conns[i+1:]...)
}

func (conns *Connections) Add(conn *Connection) {
	conns.connsLock.Lock()
	defer conns.connsLock.Unlock()
	conns.conns = append(conns.conns, conn)
}

//finding a established connection using local and remote address
func (conns *Connections) Find(localAddr net.Addr, remoteAddr net.Addr) *Connection {
	conns.connsLock.Lock()
	defer conns.connsLock.Unlock()
	for _, conn := range conns.conns {
		laddr := conn.Conn.LocalAddr()
		raddr := conn.Conn.RemoteAddr()
		if localAddr.String() == laddr.String() && remoteAddr.String() == raddr.String() {
			return conn
		}
	}
	return nil
}

//store all connections of any transport type
var allConnections *Connections = &Connections{}

//to communicate with comsumers of sip message
var sipMsgChan chan coding.SipMessage = make(chan coding.SipMessage)

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
	conn := allConnections.Find(msg.LocalAddr, msg.RemoteAddr)
	//convert SipMessage for transfer on the wire
	//put the load on the wire
	return conn.Write([]byte(msg.Stringify()))
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
