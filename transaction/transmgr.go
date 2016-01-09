package transaction

import (
	"github.com/henryscala/sipq/coding"
	"github.com/henryscala/sipq/util/container/concurrent"
)

var (
	allTrans *concurrent.Map = concurrent.NewMap()
)

//message from the transport
func handleRemoteMessage(msg *coding.SipMessage) {

}

//message from scenario file
func handleLocalMessage(msg *coding.SipMessage) {

}
