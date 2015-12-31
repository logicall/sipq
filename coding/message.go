package coding

import (
	"net"
	"strings"
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
