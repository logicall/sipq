// transport_test.go
package transport

import (
	"bufio"
	"net"
	"testing"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/config"
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

var SipMessageInviteWithLineFolding string = `
INVITE sip:bob@biloxi.com SIP/2.0
Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds
Max-Forwards: 70
To: Bob <sip:bob@biloxi.com>
From: Alice <sip:alice@atlanta.com>;tag=1928301774
Call-ID: a84b4c76e66710@pc33.atlanta.com
CSeq: 314159
    INVITE
Contact: <sip:alice@pc33.atlanta.com>
Content-Type: application/sdp
Content-Length: 5

hello
`

var SipMessage200OK string = `
SIP/2.0 200 OK
Via: SIP/2.0/UDP server10.biloxi.com;branch=z9hG4bKnashds8;received=192.0.2.3
Via: SIP/2.0/UDP bigbox3.site3.atlanta.com;branch=z9hG4bK77ef4c2312983.1;received=192.0.2.2
Via: SIP/2.0/UDP pc33.atlanta.com;branch=z9hG4bK776asdhds ;received=192.0.2.1
To: Bob <sip:bob@biloxi.com>;tag=a6c85cf
From: Alice <sip:alice@atlanta.com>;tag=1928301774
Call-ID: a84b4c76e66710@pc33.atlanta.com
CSeq: 314159 INVITE
Contact: <sip:bob@192.0.2.4>
Content-Type: application/sdp
Content-Length: 5

hello
`

func init() {
	AllServers = StartServers(config.TheExeConfig)
}

func TestUdpConn(t *testing.T) {

	raddr, _ := net.ResolveUDPAddr(UDP.String(), ServerAddress)
	client, err := CreateUdpServer(ClientAddress)
	if err != nil {
		t.Fatal("create udp client side endpoint failed")
	}
	_, err = client.UdpConn.WriteTo([]byte(util.CookSipMsg(SipMessageInvite)), raddr)
	if err != nil {
		t.Fatal("sending data failed")
	}
	msg := FetchSipMessage()
	if msg.MsgType != coding.MsgTypeRequest {
		t.Fatal("getting message failed")
	}
	if string(msg.BodyContent) != "hello" {
		t.Log("body len", len(msg.BodyContent))
		t.Fatal("getting message body failed", string(msg.BodyContent))
	}

}

func TestTcpConn(t *testing.T) {

	client, err := CreateTcpConnection(ServerAddress)

	if err != nil {
		t.Fatal("create TCP client side endpoint failed")
	}

	buff := bufio.NewWriter(client.Conn)
	_, err = buff.WriteString(util.CookSipMsg(SipMessage200OK))

	if err != nil {
		t.Fatal("sending data failed", err)
	}
	err = buff.Flush()
	if err != nil {
		t.Fatal("flush data failed", err)
	}
	msg := FetchSipMessage()
	if msg.MsgType != coding.MsgTypeResponse {
		t.Fatal("getting message failed")
	}
	if string(msg.BodyContent) != "hello" {
		t.Log("body len", len(msg.BodyContent))
		t.Fatal("getting message body failed", string(msg.BodyContent))
	}

}

func TestTcpConnWithFoldingHeader(t *testing.T) {

	client, err := CreateTcpConnection(ServerAddress)

	if err != nil {
		t.Fatal("create TCP client side endpoint failed")
	}

	buff := bufio.NewWriter(client.Conn)
	_, err = buff.WriteString(util.CookSipMsg(SipMessageInviteWithLineFolding))

	if err != nil {
		t.Fatal("sending data failed", err)
	}
	err = buff.Flush()
	if err != nil {
		t.Fatal("flush data failed", err)
	}
	msg := FetchSipMessage()
	if msg.MsgType != coding.MsgTypeRequest {
		t.Fatal("getting message failed")
	}
	if string(msg.BodyContent) != "hello" {
		t.Log("body len", len(msg.BodyContent))
		t.Fatal("getting message body failed", string(msg.BodyContent))
	}

}
