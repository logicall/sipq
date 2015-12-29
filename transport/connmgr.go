package transport

import (
	"bufio"
	"sipq/coding"
	"sipq/trace"
	"sipq/util"
	"sync"
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

//store all connections of any transport type
var allConnections *Connections = &Connections{}

//to communicate with comsumers of sip message
var sipMsgChan chan coding.SipMessage = make(chan coding.SipMessage)

func fetchSipMessageFromReader(reader *bufio.Reader) (*coding.SipMessage, error) {
	reader.ReadString(coding.LF[0])
	//TODO
	return nil, nil
}

//should be called in a go routine, since it is blocking
func handleNewData(conn *Connection) {
	var buf []byte
	var reader *bufio.Reader //a struct instead of an interface

	for {
		switch conn.TransportType {
		case TCP:
			if reader == nil {
				reader = conn.Reader()
			}
			msg, err := fetchSipMessageFromReader(reader)
			if err != nil {
				conn.Close()
				return
			}
			sipMsgChan <- *msg

		//TODO
		//get a whole message from the connections and output the message
		case UDP:
			if buf == nil {
				buf = make([]byte, coding.MaxUdpPacketLen)
			}
			n, addr, err := conn.ReadFrom(buf)
			if err != nil {
				trace.Trace.Fatalln("read failed from udp", conn) //UDP connection, not return
				continue
			}
			var msg coding.SipMessage
			msg.Raw = string(buf[0:n])
			msg.LocalAddr = conn.Conn.LocalAddr()
			msg.RemoteAddr = addr
			sipMsgChan <- msg
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
