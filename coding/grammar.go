// RFCs:
// SIP https://tools.ietf.org/html/rfc3261
// Internet Message Format https://tools.ietf.org/html/rfc2822
// ABNF https://tools.ietf.org/html/rfc5234
// HTTP/1.1 https://tools.ietf.org/html/rfc2616
// URI Gneric Syntax https://tools.ietf.org/html/rfc2396
package coding

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

//Request-Line  =  Method SP Request-URI SP SIP-Version CRLF
//Status-Line  =   SIP-Version SP Status-Code SP Reason-Phrase CRLF

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
