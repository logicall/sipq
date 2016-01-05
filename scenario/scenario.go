package scenario

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/robertkrimen/otto"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/rawstr"
)

type message struct {
	raw    string
	cooked *coding.SipMessage
	isSent bool
}

type messageErr struct {
	msg string
	err error
}

var (
	gError   *messageErr
	messages []message
	vm       *otto.Otto
)

func (mErr messageErr) String() string {
	return fmt.Sprintf("%v:%s", mErr.err, mErr.msg)
}

func handleUserFunc(call otto.FunctionCall, isSent bool) otto.Value {
	msgRaw, err := call.Argument(0).ToString()
	if err != nil {
		gError.err = err
		gError.msg = msgRaw
		return otto.FalseValue()
	}

	msgRaw = util.CookSipMsg(msgRaw)

	if isSent {
		msgCooked, err := coding.FetchSipMessageFromReader(bytes.NewReader([]byte(msgRaw)), config.IsStreamTransport())
		if err != nil {
			gError.err = err
			gError.msg = msgRaw
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
	trace.Trace.Println("enter RunText")
	defer trace.Trace.Println("exit RunText")
	compiled, err := rawstr.PreCompile(scenarioText)
	if err != nil {
		return err
	}

	_, err = vm.Run(compiled)
	if err != nil {
		return err
	}
	return nil
}

func LoadFile(scenarioFile string) error {
	bs, err := ioutil.ReadFile(scenarioFile)
	if err != nil {
		return err
	}

	return LoadText(string(bs))
}

func init() {
	vm = otto.New()
	vm.Set("send", send)
	vm.Set("recv", recv)
}
