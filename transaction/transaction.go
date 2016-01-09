package transaction

import (
	"fmt"

	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
)

type Transaction struct {
	transType int    //ICT IST NICT NIST
	state     int    //ICTCalling, ICTProceeding, ...
	branchID  string //It serves as the ID that solely identify a transaction
}

const (
	IDPrefix = "z9hG4bK"
)

func newID() string {
	return fmt.Sprintf("%s-%s", IDPrefix, util.UUID())
}

type HandleFunc func(t *Transaction, msg *coding.SipMessage)

func (t *Transaction) changeState(destState int) {
	trace.Trace.Println("transaction", t.transType, t.branchID, t.state, "->", destState)
	t.state = destState
}

func (t *Transaction) mICTCallingHandleMessage(msg *coding.SipMessage) {

}

func (t *Transaction) mICTProceedingHandleMessage(msg *coding.SipMessage) {

}

func (t *Transaction) mICTCompletedHandleMessage(msg *coding.SipMessage) {

}

func (t *Transaction) mICTTerminatedHandleMessage(msg *coding.SipMessage) {

}

func (t *Transaction) handleMessage(msg *coding.SipMessage) {
	matrix[t.transType][t.state](t, msg)
}

var matrix map[int]map[int]HandleFunc = map[int]map[int]HandleFunc{
	ICT: map[int]HandleFunc{
		ICTCalling:    (*Transaction).mICTCallingHandleMessage, //method expression
		ICTProceeding: (*Transaction).mICTProceedingHandleMessage,
		ICTCompleted:  (*Transaction).mICTCompletedHandleMessage,
		ICTTerminated: (*Transaction).mICTTerminatedHandleMessage,
	},
	IST:  map[int]HandleFunc{},
	NICT: map[int]HandleFunc{},
	NIST: map[int]HandleFunc{},
}
