package coding

const (
	//rfc2822: Each line of characters MUST be no more than
	//998 characters, and SHOULD be no more than 78 characters, excluding
	//the CRLF.
	MaxLineLen      = 1000 //can be used when decoding receiving message
	MaxPageLineLen  = 80   //can be used when encoding sending message together with folding
	MaxUdpPacketLen = 1300 //If larger than that, should switch to TCP
)
