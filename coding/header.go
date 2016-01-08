package coding

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/henryscala/sipq/util"
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

	HeaderValueSep = ","
)

var CompactHdr map[string]string = map[string]string{
	strings.ToLower(HdrVia):       "v",
	strings.ToLower(HdrSubject):   "s",
	strings.ToLower(HdrSupported): "k",
	strings.ToLower(HdrTo):        "t",
}

type SipHeader interface {
	Name() string
	Value() string
	SetName(s string)
	SetValue(s string)
	String() string
}

type SipHeaderCommon struct {
	StrName  string
	StrValue string
}

var (
	//name : value
	reHeader *regexp.Regexp = regexp.MustCompile(`(?i)(?P<name>[-\w]+)\s*:\s*(?P<value>.+)`)
	reNumber *regexp.Regexp = regexp.MustCompile(`(?P<number>\d+)`)

	//Content-Length: 349
	reHeaderContentLength *regexp.Regexp = regexp.MustCompile(`(?i)(?P<name>Content-Length)\s*:\s*(?P<number>\d+)`)
)

func (hdr *SipHeaderCommon) SetName(s string) {
	hdr.StrName = s
}

func (hdr *SipHeaderCommon) SetValue(s string) {
	hdr.StrValue = s
}

func (hdr *SipHeaderCommon) Name() string {
	return hdr.StrName
}

func (hdr *SipHeaderCommon) Value() string {
	return hdr.StrValue
}

func (hdr *SipHeaderCommon) String() string {
	values := strings.Split(hdr.StrValue, HeaderValueSep)
	var hdrstr string
	for _, value := range values {
		hdrstr += hdr.StrName + ": " + value + CRLF
	}
	return hdrstr
}

func ParseHeader(line string) (SipHeader, error) {
	result := reHeader.FindStringSubmatch(line)
	if result == nil {
		return nil, ErrNotMatch
	}
	headerName := result[1]
	headerValue := result[2]
	var err error
	var hdr SipHeader
	switch {
	case util.StrEq(headerName, HdrContentLength):
		hdr, err = parseHeaderContentLength(headerValue)
	default:
		hdr = &SipHeaderCommon{StrName: headerName, StrValue: headerValue}
	}
	if err != nil {
		return nil, err
	}
	return hdr, nil
}

type SipHeaderContentLength struct {
	SipHeaderCommon
}

func (hdr *SipHeaderContentLength) Name() string {
	return HdrContentLength
}

func (hdr *SipHeaderContentLength) Length() int {
	l, _ := strconv.Atoi(hdr.StrValue)
	return l
}

func parseHeaderContentLength(value string) (*SipHeaderContentLength, error) {

	result := reNumber.FindStringSubmatch(value)
	if result == nil {
		return nil, ErrNotMatch
	}
	hdr := &SipHeaderContentLength{SipHeaderCommon{StrName: HdrContentLength, StrValue: result[1]}}
	return hdr, nil
}
