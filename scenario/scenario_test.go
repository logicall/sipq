package scenario

import (
	"strings"
	"testing"
)

func TestScenario(t *testing.T) {
	scenario := "msg0=`"
	scenario += `
INVITE sip:bob@biloxi.com SIP/2.0
Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds
Max-Forwards: 70
To: Bob <sip:bob@biloxi.com>
From: Alice <sip:alice@atlanta.com>;tag=1928301774
Call-ID: a84b4c76e66710@pc33.atlanta.com
CSeq: 314159 INVITE
Contact: <sip:alice@pc33.atlanta.com>
Content-Type: application/sdp
Content-Length: 5

hello
`
	scenario += "`;\n"

	scenario += "msg2=`"
	scenario += `
ACK sip:bob@biloxi.com SIP/2.0
Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhdt
Max-Forwards: 70
To: Bob <sip:bob@biloxi.com>
From: Alice <sip:alice@atlanta.com>;tag=1928301774
Call-ID: a84b4c76e66710@pc33.atlanta.com
CSeq: 314160 ACK
Contact: <sip:alice@pc33.atlanta.com>
Content-Type: application/sdp
Content-Length: 5

hello
`

	scenario += "`;\n"

	scenario += `
		send(msg0);
		recv("200");
		send(msg2);
	`

	err := LoadText(scenario)
	if err != nil {
		t.Error("not expected", err)
	}
	if len(messages) != 3 {
		t.Error("not expected", len(messages))
	}

	if !(strings.Contains(messages[0].raw, "INVITE")) {
		t.Error(messages[0].raw)
	}

	if !(strings.Contains(messages[1].raw, "200")) {
		t.Error(messages[1].raw)
	}

	if !(strings.Contains(messages[2].raw, "ACK")) {
		t.Error(messages[2].raw)
	}
}
