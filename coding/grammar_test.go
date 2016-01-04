package coding

import (
	"regexp"
	"testing"

	"github.com/henryscala/sipq/util"
)

//for testing my regexp skills
func TestRegExp(t *testing.T) {
	line := " sip:master@qinmishu.org SIP/2.0 "
	re := regexp.MustCompile(`(?i) (?P<uri>[:@.\w]+) (?P<version>SIP/2\.0|SIP/1\.0)`)
	result := re.FindStringSubmatch(line)

	if result == nil {
		t.Fatal(re, "not match", line)
	}
}

func TestParseStartLine(t *testing.T) {
	requstLineStr := " INVITE sip:henry@qinmishu.org SIP/2.0 "

	startLine, msgType, err := ParseStartLine(requstLineStr)
	if err != nil {
		t.Fatal(err)
	}
	if msgType != MsgTypeRequest {
		t.Fail()
	}
	requestLine := startLine.(*RequestLine)
	if requestLine.Method != "INVITE" {
		t.Fatal(requestLine.Method)
	}
	if requestLine.Uri != "sip:henry@qinmishu.org" {
		t.Fatal(requestLine.Uri)
	}
	if requestLine.Version != "SIP/2.0" {
		t.Fatal(requestLine.Version)
	}

	statusLineStr := "SIP/2.0 200 OK"
	startLine, msgType, err = ParseStartLine(statusLineStr)
	if err != nil {
		t.Fatal(err)
	}
	if msgType != MsgTypeResponse {
		t.Fatal(msgType)
	}

	statusLine := startLine.(*StatusLine)
	if statusLine.Reason != "OK" {
		t.Fatal(statusLine.Reason)
	}
	if statusLine.Status != 200 {
		t.Fatal(statusLine.Status)
	}
	if statusLine.Version != "SIP/2.0" {
		t.Fatal(statusLine.Version)
	}
}

func TestParseHeader(t *testing.T) {
	line := "Unknown: unknown123.qinmishu.org"

	hdr, err := ParseHeader(line)
	if err != nil {
		t.Fatal("failed parse", line)
	}

	if !util.StrEq(hdr.Name(), "unknown") {
		t.Fatal("failed parse", hdr.Name())
	}

	if hdr.Value() != "unknown123.qinmishu.org" {
		t.Fatal("failed parse", hdr.Value())
	}
	_, ok := hdr.(*SipHeaderCommon)
	if !ok {
		t.Fatal("invalid type")
	}
}

func TestParseHeaderContentLength(t *testing.T) {
	line := "Content-Length: 349"
	var hdrContentLength *SipHeaderContentLength
	hdr, err := ParseHeader(line)
	if err != nil {
		t.Fatal("failed parse", line)
	}

	t.Log(hdr)
	hdrContentLength = hdr.(*SipHeaderContentLength)

	if !util.StrEq(hdrContentLength.Name(), "Content-Length") {
		t.Fatal("failed parse", hdr.Name())
	}
	if hdrContentLength.StrValue != "349" {
		t.Fatal("failed parse", hdrContentLength.StrValue)
	}
	if hdrContentLength.Length() != 349 {
		t.Fatal("failed parse", hdrContentLength.Length())
	}
}
