package coding

import (
	"bytes"
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
	msgstr += CRLF //even if there is not body, a empty CRLF shall be included
	if len(msg.BodyContent) > 0 {
		msgstr += string(msg.BodyContent)
	}
	return msgstr
}

//for debugging whether bufio.NewReader can be called on a reader multiple times.
//it showed that it can be called multiple times.
//in future, consider using bufio.NewReader in FetchSipMessageFromReader
func readString(reader io.Reader, delim byte) (line string, err error) {
	var buf bytes.Buffer
	var b []byte = make([]byte, 1) //one byte slice
	for {
		n, err := reader.Read(b)
		if err != nil {
			return "", err
		}
		if n != len(b) {
			return "", ErrInvalidLine
		}

		n, err = buf.Write(b)
		if err != nil {
			return "", err
		}
		if n != len(b) {
			return "", ErrInvalidLine
		}
		if b[0] == delim {
			break
		}
	}
	return buf.String(), nil
}

//for debugging whether bufio.NewReader can be called on a reader multiple times.
//it showed that it can be called multiple times.
//in future, consider using bufio.NewReader in FetchSipMessageFromReader
func readByte(reader io.Reader) (byte, error) {

	var b []byte = make([]byte, 1) //one byte slice

	n, err := reader.Read(b)
	if err != nil {
		return b[0], err
	}
	if n != len(b) {
		return b[0], ErrInvalidLine
	}

	return b[0], nil

}

//if error is EOF, need to handle specially
func FetchSipMessageFromReader(reader io.Reader, isStreamTransport bool) (*SipMessage, error) {
	trace.Trace("enter FetchSipMessageFromReader")
	defer trace.Trace("exit FetchSipMessageFromReader")

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
			trace.Debug(ErrInvalidLine, line)
			return ErrInvalidLine
		}
		sipMessage.AddHeader(hdr)
		return nil
	}

	state = expectingStartLine
	for {
		switch state {
		case expectingStartLine:
			trace.Debug("FetchSipMessageFromReader state", expectingStartLine)
			line, err = readString(reader, LF[0])
			trace.Debug("FetchSipMessageFromReader line", line)
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
				trace.Debug(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			line = util.StrTrim(line)
			if line == "" {
				continue
			}

			startLine, msgType, err := ParseStartLine(line)
			if err != nil {
				trace.Debug(err)
				return nil, err
			}
			sipMessage.MsgType = msgType
			sipMessage.StartLine = startLine
			trace.Debug("FetchSipMessageFromReader state", state, "->", expectingHeader)
			state = expectingHeader

		case expectingHeader:
			line, err = readString(reader, LF[0])

			if err != nil {
				if err == io.EOF {
					return sipMessage, err
				}
				return nil, err
			}
			lineLen = len(line)
			if lineLen < 2 {
				trace.Debug(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			if !strings.HasSuffix(line, CRLF) {
				trace.Debug(ErrInvalidLine)
				return nil, ErrInvalidLine
			}
			if lineLen == 2 {
				if lineCache != "" {
					err = parseAndAddHeader(lineCache)
					if err != nil {
						return nil, err
					}
				}
				trace.Debug("FetchSipMessageFromReader state", state, "->", expectingBody)
				state = expectingBody
				hdr, err = sipMessage.GetHeader(HdrContentLength)
				trace.Debug("FetchSipMessageFromReader state", "hdr", hdr, "err", err)
				switch isStreamTransport {
				default:
					if err != nil {
						continue
					}
				//content length header is mandatory for stream based transport
				case true:

					if err != nil {
						return nil, ErrInvalidMsg
					}

				}
				contentLengthHdr = hdr.(*SipHeaderContentLength)
				//in case content-length exist, but the value is 0
				if contentLengthHdr.Length() <= 0 {
					return sipMessage, nil
				}
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
				trace.Debug(ErrInvalidLine, lineCache)
				return nil, ErrInvalidLine
			}

			// cache the new line
			lineCache = line

		case expectingBody:
			for {

				b, err := readByte(reader)
				trace.Debug("FetchSipMessageFromReader read a byte", string(b), "err", err, "contentLengthHdr", contentLengthHdr)
				if err != nil {
					if err == io.EOF {
						switch isStreamTransport {
						default:
							return sipMessage, err
						case true: //for streamed based connection, there must be a content length header
							if len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
								return sipMessage, err
							} else {
								trace.Debug(ErrInvalidMsg)
								return nil, ErrInvalidMsg
							}
						}
					}
					trace.Debug(ErrInvalidMsg)
					return nil, ErrInvalidMsg
				}

				sipMessage.BodyContent = append(sipMessage.BodyContent, b)
				if contentLengthHdr != nil && len(sipMessage.BodyContent) >= contentLengthHdr.Length() {
					return sipMessage, nil
				}
			} //inner for
		} //switch
	} //for

	err = fmt.Errorf("unexpected")
	panic(err)
	return nil, err
}
