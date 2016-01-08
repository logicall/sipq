package transaction

import (
	"github.com/henryscala/sipq/coding"
)

type Transaction struct {
	transType int //ICT IST NICT NIST
	state     int //ICTCalling, ICTProceeding, ...
}

type HandleFunc func(t *Transaction, msg *coding.SipMessage)

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
