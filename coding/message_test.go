package coding

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/henryscala/sipq/util"
)

var SipMessageInvite string = `
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

func TestSipMsgStringify(t *testing.T) {
	rawmsg := util.CookSipMsg(SipMessageInvite)
	msgObj, _ := FetchSipMessageFromReader(bytes.NewReader([]byte(rawmsg)), true)
	msgStrOnTheWire := msgObj.Stringify()
	fmt.Println(msgStrOnTheWire)
	//try to feed the stringified msg again, to see if it's right
	newMsgObj, _ := FetchSipMessageFromReader(bytes.NewReader([]byte(msgStrOnTheWire)), true)
	if string(newMsgObj.BodyContent) != "hello" {
		t.Log("body len", len(newMsgObj.BodyContent))
		t.Fatal("getting message body failed", string(newMsgObj.BodyContent))
	}
}
