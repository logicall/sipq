package rawstr

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	SingleQuote     = "'"
	DoubleQuote     = "\""
	DoubleQuoteRepr = "\\\""
	BackQuote       = "`"
	Quotes          = SingleQuote + DoubleQuote + BackQuote

	CR      = "\r"
	CR_REPR = "\\r"
	LF      = "\n"
	LF_REPR = "\\n"
	CRLF    = "\r\n"
)

var ErrCompile error = fmt.Errorf("Compile error(Unbalanced back quotes)")

func isQuote(b byte) bool {
	if strings.Contains(Quotes, string(b)) {
		return true
	}
	return false
}

//convert the text that may contain multiples line quoted by BackQuote
//to single line.
//It treat the text as plain ascii
func PreCompile(text string) (string, error) {
	strs := strings.Split(text, BackQuote)
	//there should be even number of backquote, and odd number of segments separated by it
	if len(strs)%2 == 0 {
		return text, ErrCompile
	}
	var result bytes.Buffer

	for i, str := range strs {
		if i%2 == 0 {
			result.WriteString(str)
		} else {
			result.WriteString(toSingleLine(str))
		}
	}
	return result.String(), nil
}

//Convert a multiple line text to a single line
func toSingleLine(multiline string) string {
	var buffer bytes.Buffer
	buffer.WriteString(DoubleQuote)
	for i := 0; i < len(multiline); i++ {
		c := multiline[i]
		switch c {
		case CR[0]:
			buffer.WriteString(CR_REPR)
		case LF[0]:
			buffer.WriteString(LF_REPR)
		case DoubleQuote[0]:
			buffer.WriteString(DoubleQuoteRepr)
		default:
			buffer.WriteByte(c)
		}

	}
	buffer.WriteString(DoubleQuote)
	return buffer.String()
}
