package transaction

const (
	ICT int = iota
	IST
	NICT
	NIST
)

const (
	ICTCalling int = iota
	ICTProceeding
	ICTCompleted
	ICTTerminated
)

const (
	ISTProceeding int = iota
	ISTCompleted
	ISTConfirmed
	ISTTerminated
)

const (
	NICTTrying int = iota
	NICTProceeding
	NICTCompleted
	NICTTerminated
)

const (
	NISTTrying int = iota
	NISTProceeding
	NISTCompleted
	NISTTerminated
)
