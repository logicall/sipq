package scenario

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"os"

	"github.com/robertkrimen/otto"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/transport"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/rawstr"
)

var (
	gError   *messageErr
	messages []message
	vm       *otto.Otto
)

type message struct {
	raw    string
	cooked *coding.SipMessage
	isSent bool //sending from sipq
}

type messageErr struct {
	msg string
	err error
}

type Scenario struct {
	messages []message
	index    int //used to check to which step the scenario has run
}

func New() *Scenario {
	s := &Scenario{}
	s.messages = messages
	return s
}

func (mErr messageErr) String() string {
	return fmt.Sprintf("%v:%s", mErr.err, mErr.msg)
}

func (self *Scenario) Run(success chan<- bool) {
	trace.Trace("enter Run")
	defer trace.Trace("exit Run")
	remoteAddr, err := util.Addr(config.RemoteIP, config.RemotePort, config.TransportType)
	if err != nil {
		trace.Debug("parse remote addr failed")
		success <- false
		return
	}
	localAddr, err := util.Addr(config.LocalIP, config.LocalPort, config.TransportType)
	if err != nil {
		trace.Debug("parse local addr failed")
		success <- false
		return
	}

	transportType := transport.Type(config.TransportType)

	for ; self.index < len(self.messages); self.index++ {
		msg := self.messages[self.index]

		if msg.isSent {
			msg.cooked.LocalAddr = localAddr
			msg.cooked.RemoteAddr = remoteAddr
			err = transport.Send(msg.cooked, transportType)
			if err != nil {
				trace.Debug("send message failed", err)
				success <- false
				return
			}
			trace.Debug("<<<sent message", msg.cooked.StartLine)
		} else {
			msgReceived := transport.FetchSipMessage()
			if err = verifyMessage(msg.cooked, msgReceived); err != nil {
				trace.Debug("verify message failed")
				success <- false
				return
			}
			trace.Debug("<<<received message", msgReceived.StartLine)
		}
	}

	success <- true
}

func verifyMessage(msgExpected, msgReceived *coding.SipMessage) error {
	return nil //TODO
}

func handleUserFunc(call otto.FunctionCall, isSent bool) otto.Value {
	msgRaw, err := call.Argument(0).ToString()
	if err != nil {
		gError.err = err
		gError.msg = msgRaw
		return otto.FalseValue()
	}
	trace.Debug("handleUserFunc", "isSent", isSent)
	msgRaw = util.CookSipMsg(msgRaw)

	if isSent {
		msgCooked, err := coding.FetchSipMessageFromReader(bytes.NewReader([]byte(msgRaw)), config.IsStreamTransport())
		if err != nil && err != io.EOF {
			gError.err = err
			gError.msg = msgRaw
			trace.Debug("FetchSipMessageFromReader failed", err, msgRaw)
			return otto.FalseValue()
		}
		msg := message{raw: msgRaw, cooked: msgCooked, isSent: isSent}
		messages = append(messages, msg)
	} else {
		msg := message{raw: msgRaw, isSent: isSent}
		messages = append(messages, msg)
	}

	return otto.TrueValue()
}

func send(call otto.FunctionCall) otto.Value {
	return handleUserFunc(call, true)
}

func recv(call otto.FunctionCall) otto.Value {
	return handleUserFunc(call, false)
}

func LoadText(scenarioText string) error {
	trace.Trace("enter LoadText")
	defer trace.Trace("exit LoadText")
	compiled, err := rawstr.PreCompile(scenarioText)
	if err != nil {
		return err
	}

	_, err = vm.Run(compiled)
	if err != nil {
		trace.Debug("run scenario failed", err, compiled)
		return err
	}
	return nil
}

func LoadFile(scenarioFile string) error {

	bs, err := ioutil.ReadFile(scenarioFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Please specify a scenario file")
		}
		return err
	}

	return LoadText(string(bs))
}

func init() {
	vm = otto.New()
	vm.Set("send", send)
	vm.Set("recv", recv)
}
