package coding

import (
	"strings"
)

const (
	HdrVia           = "VIA"
	HdrMaxForwards   = "MAX-FORWARDS"
	HdrTo            = "To"
	HdrFrom          = "From"
	HdrCallId        = "call-id"
	HdrCSeq          = "CSeq"
	HdrContact       = "Contact"
	HdrContentType   = "Content-Type"
	HdrContentLength = "Content-Length"
	HdrSupported     = "Supported"
	HdrRequire       = "Require"
	HdrProxyRequire  = "Proxy-Require"
	HdrExpires       = "Expires"
	HdrServer        = "Server"
	HdrSubject       = "Subject"
)

var CompactHdr map[string]string = map[string]string{
	strings.ToLower(HdrVia):       "v",
	strings.ToLower(HdrSubject):   "s",
	strings.ToLower(HdrSupported): "k",
	strings.ToLower(HdrTo):        "t",
}
