package transaction

const (
	ICT int = iota
	IST
	NICT
	NIST
)

const (
	ICTIdle int = iota
	ICTCalling
	ICTProceeding
	ICTCompleted
	ICTTerminated
)

const (
	ISTIdle int = iota
	ISTProceeding
	ISTCompleted
	ISTConfirmed
	ISTTerminated
)

const (
	NICTIdle int = iota
	NICTTrying
	NICTProceeding
	NICTCompleted
	NICTTerminated
)

const (
	NISTIdle int = iota
	NISTTrying
	NISTProceeding
	NISTCompleted
	NISTTerminated
)
