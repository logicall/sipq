// RFCs:
// SIP https://tools.ietf.org/html/rfc3261
// Internet Message Format https://tools.ietf.org/html/rfc2822
// ABNF https://tools.ietf.org/html/rfc5234
// HTTP/1.1 https://tools.ietf.org/html/rfc2616
// URI Gneric Syntax https://tools.ietf.org/html/rfc2396
package coding

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/henryscala/sipq/trace"
)

//Grammar
const (
	TAB = "\t"
	SP  = " "

	CR     = "\r"
	LF     = "\n"
	COLON  = ":"
	CRLF   = CR + LF
	ESCAPE = "%"

	UriSchemeSIP  = "sip"
	UriSchemeSIPS = "sips"
	UriSchemeTEL  = "tel"

	SIPVersion1 = "SIP/1.0"
	SIPVersion2 = "SIP/2.0"
)

var (
	//Request-Line  =  Method SP Request-URI SP SIP-Version CRLF
	reRequestLine *regexp.Regexp = regexp.MustCompile(`(?i)(?P<method>\w+) (?P<uri>[:@.\w]+) (?P<version>SIP/[12]\.0)`)

	//Status-Line  =   SIP-Version SP Status-Code SP Reason-Phrase CRLF
	reStatusLine *regexp.Regexp = regexp.MustCompile(`(?i)(?P<version>SIP/[12]\.0) (?P<status>\d+) (?P<reason>[ \w]+)`)
)

var (
	ErrInvalidLine error = fmt.Errorf("invalid line")
	ErrNotMatch    error = fmt.Errorf("not match")
	ErrNotFound    error = fmt.Errorf("not found")
	ErrInvalidMsg  error = fmt.Errorf("invalid message")
)

type StartLine interface {
	Stringify() string
}

type RequestLine struct {
	Method, Uri, Version string
}

type StatusLine struct {
	Version, Reason string
	Status          int
}

func (reql *RequestLine) Stringify() string {
	s := fmt.Sprintf("%s %s %s", reql.Method, reql.Uri, reql.Version)
	return s
}

func (stl *StatusLine) Stringify() string {
	s := fmt.Sprintf("%s %d %s", stl.Version, stl.Status, stl.Reason)
	return s
}

func ParseStartLine(line string) (startLine StartLine, msgType int, err error) {
	trace.Trace.Println("enter ParseStartLine", line)
	defer trace.Trace.Println("exit ParseStartLine", line)
	var result []string

	result = reRequestLine.FindStringSubmatch(line)

	if result == nil {
		result = reStatusLine.FindStringSubmatch(line)

		if result == nil {
			return nil, MsgTypeInvalid, ErrInvalidLine
		} else {
			statusLine := &StatusLine{}
			statusLine.Version = result[1]
			statusLine.Reason = result[3]
			n, _ := strconv.Atoi(result[2])
			statusLine.Status = n
			return statusLine, MsgTypeResponse, nil
		}
	} else {
		requestLine := &RequestLine{}
		requestLine.Method = result[1]
		requestLine.Uri = result[2]
		requestLine.Version = result[3]
		return requestLine, MsgTypeRequest, nil
	}
}

//implementations should avoid spaces between the field
//name and the colon and use a single space (SP) between the colon and
//the field-value.

//Implementations processing SIP messages over stream-oriented
//transports MUST ignore any CRLF appearing before the start-line

//SIP follows the requirements and guidelines of RFC 2396 [5] when
//defining the set of characters that must be escaped in a SIP URI, and
//uses its ""%" HEX HEX" mechanism for escaping.

//If a request is within 200 bytes of the path MTU, or if it is larger
//than 1300 bytes and the path MTU is unknown, the request MUST be sent
//using an RFC 2914 [43] congestion controlled transport protocol, such
//as TCP.
