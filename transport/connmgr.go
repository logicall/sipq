package transport

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
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

//if error is EOF, need to handle specially
func fetchSipMessageFromReader_bak(reader io.Reader, transportType TransportType) (*coding.SipMessage, error) {

	var bufReader *bufio.Reader

	bufReader = bufio.NewReader(reader)

	const (
		expectingStartLine int = iota
		expectingHeader
		expectingBody
	)

	var sipMessage *coding.SipMessage = &coding.SipMessage{}
	var line string
	var lineLen int
	var err error
	var state int
	var contentLengthHdr *coding.SipHeaderContentLength
	var hdr coding.SipHeader

	state = expectingStartLine
	for {
		switch state {
		case expectingStartLine:
			line, err = bufReader.ReadString(coding.LF[0])

			if err != nil {
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				//tolerate empty line
				if util.StrTrim(line) == "" {
					continue
				} else {
					return nil, coding.ErrInvalidLine
				}
			}
			if !strings.HasSuffix(line, coding.CRLF) {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			line = util.StrTrim(line)
			if line == "" {
				continue
			}

			startLine, msgType, err := coding.ParseStartLine(line)
			if err != nil {
				trace.Trace.Println(err)
				return nil, err
			}
			sipMessage.MsgType = msgType
			sipMessage.StartLine = startLine
			state = expectingHeader
		case expectingHeader:
			line, err = bufReader.ReadString(coding.LF[0])

			if err != nil {
				if err == io.EOF {
					return sipMessage, err
				}
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			if !strings.HasSuffix(line, coding.CRLF) {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			if lineLen == 2 {
				state = expectingBody
				hdr, err = sipMessage.GetHeader(coding.HdrContentLength)
				switch transportType {
				case UDP, SCTP:
					if err != nil {
						continue
					}
				//content length header is mandatory for stream based transport
				default:

					if err != nil {
						return nil, coding.ErrInvalidMsg
					}

				}
				contentLengthHdr = hdr.(*coding.SipHeaderContentLength)
				continue
			}
			hdr, err = coding.ParseHeader(line)
			if err != nil {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			sipMessage.AddHeader(hdr)

		case expectingBody:
			for {

				b, err := bufReader.ReadByte()
				if err != nil {
					if err == io.EOF {
						switch transportType {
						case UDP, SCTP:
							return sipMessage, err
						default:
							if len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
								return sipMessage, err
							} else {
								trace.Trace.Println(coding.ErrInvalidMsg)
								return nil, coding.ErrInvalidMsg
							}
						}
					}
					trace.Trace.Println(coding.ErrInvalidMsg)
					return nil, coding.ErrInvalidMsg
				}
				if contentLengthHdr != nil && len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
					return sipMessage, nil
				}
				sipMessage.BodyContent = append(sipMessage.BodyContent, b)
			} //inner for
		} //switch
	} //for

	err = fmt.Errorf("unexpected")
	panic(err)
	return nil, err
}

//if error is EOF, need to handle specially
func fetchSipMessageFromReader(reader io.Reader, transportType TransportType) (*coding.SipMessage, error) {

	var bufReader *bufio.Reader

	bufReader = bufio.NewReader(reader)

	const (
		expectingStartLine int = iota
		expectNewLine
		expectingHeader
		expectingFolding
		expectingBody
	)

	var sipMessage *coding.SipMessage = &coding.SipMessage{}
	var line string
	var headerLogicLine string
	var lineLen int
	var err error
	var state int
	var contentLengthHdr *coding.SipHeaderContentLength
	var hdr coding.SipHeader

	state = expectingStartLine
	for {
		switch state {
		case expectingStartLine:
			line, err = bufReader.ReadString(coding.LF[0])

			if err != nil {
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				//tolerate empty line
				if util.StrTrim(line) == "" {
					continue
				} else {
					return nil, coding.ErrInvalidLine
				}
			}
			if !strings.HasSuffix(line, coding.CRLF) {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			line = util.StrTrim(line)
			if line == "" {
				continue
			}

			startLine, msgType, err := coding.ParseStartLine(line)
			if err != nil {
				trace.Trace.Println(err)
				return nil, err
			}
			sipMessage.MsgType = msgType
			sipMessage.StartLine = startLine
			state = expectNewLine
			headerLogicLine = ""
		case expectNewLine:
			line, err = bufReader.ReadString(coding.LF[0])
			if err != nil {
				if err == io.EOF {
					if len(headerLogicLine) > 0 {
						hdr, err = coding.ParseHeader(headerLogicLine)
						if err != nil {
							return nil, coding.ErrInvalidLine
						}
						sipMessage.AddHeader(hdr)
					}
					headerLogicLine = ""
					return sipMessage, err
				}
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			if !strings.HasSuffix(line, coding.CRLF) {
				trace.Trace.Println(coding.ErrInvalidLine)
				return nil, coding.ErrInvalidLine
			}
			if lineLen == 2 {
				state = expectingBody
				continue
			}
			hdr, err = coding.ParseHeader(line)
			// folding start
			if err != nil {
				state = expectingFolding
				continue
			}
			// header start
			state = expectingHeader
		case expectingHeader:
			// handle previous header
			state = expectNewLine
			if len(headerLogicLine) > 0 {
				hdr, err = coding.ParseHeader(headerLogicLine)
				if err != nil {
					return nil, coding.ErrInvalidLine
				}
				sipMessage.AddHeader(hdr)
			}
			headerLogicLine = line

		case expectingFolding:
			state = expectNewLine
			// handle folding line
			headerLogicLine += strings.TrimRight(headerLogicLine, "\r\n")
			headerLogicLine += string(coding.SP)
			headerLogicLine += strings.TrimLeft(line, " \t")

		case expectingBody:
			if len(headerLogicLine) > 0 {
				hdr, err = coding.ParseHeader(headerLogicLine)
				if err != nil {
					return nil, coding.ErrInvalidLine
				}
				sipMessage.AddHeader(hdr)
			}
			headerLogicLine = ""

			hdr, err = sipMessage.GetHeader(coding.HdrContentLength)
			switch transportType {
			case UDP, SCTP:
				if err != nil {
					//it's ok
				}
			//content length header is mandatory for stream based transport
			default:
				if err != nil {
					return nil, coding.ErrInvalidMsg
				}
			}
			contentLengthHdr = hdr.(*coding.SipHeaderContentLength)
			for {

				b, err := bufReader.ReadByte()
				if err != nil {
					if err == io.EOF {
						switch transportType {
						case UDP, SCTP:
							return sipMessage, err
						default:
							if len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
								return sipMessage, err
							} else {
								trace.Trace.Println(coding.ErrInvalidMsg)
								return nil, coding.ErrInvalidMsg
							}
						}
					}
					trace.Trace.Println(coding.ErrInvalidMsg)
					return nil, coding.ErrInvalidMsg
				}
				if contentLengthHdr != nil && len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
					return sipMessage, nil
				}
				sipMessage.BodyContent = append(sipMessage.BodyContent, b)
			} //inner for
		} //switch
	} //for

	err = fmt.Errorf("unexpected")
	panic(err)
	return nil, err
}

//should be called in a go routine, since it is blocking
func handleNewData(conn *Connection) {
	var buf []byte = make([]byte, coding.MaxUdpPacketLen)

	for {
		switch conn.TransportType {
		case TCP:
			laddr := conn.Conn.LocalAddr()
			raddr := conn.Conn.RemoteAddr()
			msg, err := fetchSipMessageFromReader(conn.Conn, TCP)
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
			msg, err := fetchSipMessageFromReader(udpReader, UDP)
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
