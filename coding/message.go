package coding

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
)

//Message
const (
	MsgTypeInvalid  = 0
	MsgTypeRequest  = 1
	MsgTypeResponse = 2

	MethodInvite   = "INVITE"
	MethodAck      = "ACK"
	MethodBye      = "BYE"
	MethodRegister = "REGISTER"
	MethodCancel   = "CANCEL"
	MethodOptions  = "OPTIONS"
	MethodPrack    = "PRACK"
)

var ReasonPhrase map[int]string = map[int]string{
	100: "Trying",
	180: "Ringing",
	183: "Session Progress",
	200: "OK",
	301: "Moved Permanently",
	302: "Moved Temporarily",
	400: "Bad Request",
	401: "Unauthorized",
	480: "Temporarily Unavailable",
	482: "Loop Detected",
	486: "Busy Here",
	500: "Server Internal Error",
	501: "Not Implemented",
	503: "Service Unavailable",
	600: "Busy Everywhere",
}

//stands for a whole SIP message
type SipMessage struct {
	MsgType    int //request or response
	StartLine  StartLine
	LocalAddr  net.Addr
	RemoteAddr net.Addr
	HeaderMap  map[string]SipHeader

	BodyContent []byte
}

func (msg *SipMessage) GetHeader(name string) (SipHeader, error) {
	name = strings.ToLower(name)
	hdr, ok := msg.HeaderMap[name]
	if ok {
		return hdr, nil
	}
	return nil, ErrNotFound
}

func (msg *SipMessage) AddHeader(hdr SipHeader) {
	if msg.HeaderMap == nil {
		msg.HeaderMap = make(map[string]SipHeader)
	}
	headerKey := strings.ToLower(hdr.Name())
	oldHdr, ok := msg.HeaderMap[headerKey]

	if ok {
		oldHdr.SetValue(oldHdr.Value() + HeaderValueSep + hdr.Value())
	} else {

		msg.HeaderMap[headerKey] = hdr
	}
}

func (msg *SipMessage) Add(headerName, headerValue string) {
	if msg.HeaderMap == nil {
		msg.HeaderMap = make(map[string]SipHeader)
	}

	headerKey := strings.ToLower(headerName)
	oldHdr, ok := msg.HeaderMap[headerKey]
	if ok {
		oldHdr.SetValue(oldHdr.Value() + HeaderValueSep + headerValue)
	} else {
		hdr := &SipHeaderCommon{StrName: headerName, StrValue: headerValue}
		msg.HeaderMap[headerKey] = hdr
	}
}

func (msg *SipMessage) String() string {
	var msgstr string
	msgstr += msg.StartLine.String() + CRLF
	for _, hdr := range msg.HeaderMap {
		msgstr += hdr.String()
	}
	if len(msg.BodyContent) > 0 {
		msgstr += CRLF
		msgstr += string(msg.BodyContent)
	}
	return msgstr
}

//if error is EOF, need to handle specially
func FetchSipMessageFromReader(reader io.Reader, isStreamTransport bool) (*SipMessage, error) {
	trace.Trace.Println("enter FetchSipMessageFromReader")
	defer trace.Trace.Println("exit FetchSipMessageFromReader")
	var bufReader *bufio.Reader

	bufReader = bufio.NewReader(reader)

	const (
		expectingStartLine int = iota
		expectingHeader
		expectingBody
	)

	var sipMessage *SipMessage = &SipMessage{}
	var lineCache string // catche a line to handle the folding case
	var line string      //currently handling line
	var lineLen int      //length of the currently handling line
	var err error
	var state int
	var contentLengthHdr *SipHeaderContentLength
	var hdr SipHeader

	parseAndAddHeader := func(line string) error {
		//regard the content in the cache as a complete line
		hdr, err = ParseHeader(line)
		if err != nil {
			trace.Trace.Println(ErrInvalidLine, line)
			return ErrInvalidLine
		}
		sipMessage.AddHeader(hdr)
		return nil
	}

	state = expectingStartLine
	for {
		switch state {
		case expectingStartLine:
			line, err = bufReader.ReadString(LF[0])

			if err != nil {
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				//tolerate empty line
				if util.StrTrim(line) == "" {
					continue
				} else {
					return nil, ErrInvalidLine
				}
			}
			if !strings.HasSuffix(line, CRLF) {
				trace.Trace.Println(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			line = util.StrTrim(line)
			if line == "" {
				continue
			}

			startLine, msgType, err := ParseStartLine(line)
			if err != nil {
				trace.Trace.Println(err)
				return nil, err
			}
			sipMessage.MsgType = msgType
			sipMessage.StartLine = startLine
			state = expectingHeader
		case expectingHeader:
			line, err = bufReader.ReadString(LF[0])

			if err != nil {
				if err == io.EOF {
					return sipMessage, err
				}
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				trace.Trace.Println(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			if !strings.HasSuffix(line, CRLF) {
				trace.Trace.Println(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			if lineLen == 2 {
				if lineCache != "" {
					err = parseAndAddHeader(lineCache)
					if err != nil {
						return nil, err
					}
				}

				state = expectingBody
				hdr, err = sipMessage.GetHeader(HdrContentLength)
				switch isStreamTransport {
				case true:
					if err != nil {
						continue
					}
				//content length header is mandatory for stream based transport
				default:

					if err != nil {
						return nil, ErrInvalidMsg
					}

				}
				contentLengthHdr = hdr.(*SipHeaderContentLength)
				continue
			}

			// cache the line
			if lineCache == "" {
				lineCache = line
				continue
			}

			// append the line to cache
			if strings.HasPrefix(line, SP) ||
				strings.HasPrefix(line, TAB) {
				lineCache = util.StrTrim(lineCache) + SP + line
				continue
			}

			//regard the content in the cache as a complete line
			err = parseAndAddHeader(lineCache)

			if err != nil {
				trace.Trace.Println(ErrInvalidLine, lineCache)
				return nil, ErrInvalidLine
			}

			// cache the new line
			lineCache = line

		case expectingBody:
			for {

				b, err := bufReader.ReadByte()
				if err != nil {
					if err == io.EOF {
						switch isStreamTransport {
						case true:
							return sipMessage, err
						default:
							if len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
								return sipMessage, err
							} else {
								trace.Trace.Println(ErrInvalidMsg)
								return nil, ErrInvalidMsg
							}
						}
					}
					trace.Trace.Println(ErrInvalidMsg)
					return nil, ErrInvalidMsg
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
