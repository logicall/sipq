package scenario

import (
	"fmt"
	"io/ioutil"

	"github.com/robertkrimen/otto"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/rawstr"
)

type message struct {
	raw    string
	cooked coding.SipMessage
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

func send(call otto.FunctionCall) otto.Value {
	msgRaw, err := call.Argument(0).ToString()
	if err != nil {
		gError.err = err
		gError.msg = msgRaw
		return otto.FalseValue()
	}
	msg := message{raw: msgRaw, isSent: true}
	messages = append(messages, msg)
	return otto.TrueValue()
}

func recv(call otto.FunctionCall) otto.Value {
	msgRaw, err := call.Argument(0).ToString()
	if err != nil {
		gError.err = err
		gError.msg = msgRaw
		return otto.FalseValue()
	}
	msg := message{raw: msgRaw, isSent: false}
	messages = append(messages, msg)
	return otto.TrueValue()
}

func RunText(scenarioText string) error {
	compiled, err := rawstr.PreCompile(scenarioText)
	if err != nil {
		return err
	}

	_, err = vm.Run(util.CookSipMsg(compiled))
	if err != nil {
		return err
	}
	return nil
}

func RunFile(scenarioFile string) error {
	bs, err := ioutil.ReadFile(scenarioFile)
	if err != nil {
		return err
	}

	return RunText(string(bs))
}

func init() {
	vm = otto.New()
	vm.Set("send", send)
	vm.Set("recv", recv)
}
