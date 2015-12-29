package coding

import (
	"net"
)

//Message
const (
	MsgTypeRequest  = 0
	MsgTypeResponse = 1

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
	LocalAddr  net.Addr
	RemoteAddr net.Addr
	Method     string //INVITE,ACK,etc
	Raw        string //all the contents
}
