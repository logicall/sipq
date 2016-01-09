package util

import (
	"fmt"

	"github.com/henryscala/sipq/trace"

	"bytes"
	"strings"
)

var (
	uuidNum int
)

func ErrorPanic(err error) {
	if err != nil {
		trace.Trace.Fatalln(err)
	}
}

//ignore case compare
func StrEq(s1, s2 string) bool {
	if strings.ToLower(s1) == strings.ToLower(s2) {
		return true
	}
	return false
}

func StrTrim(s string) string {
	return strings.Trim(s, " \t\r\n")
}

//replace \n with \r\n.
//if it is already "\r\n", then takes no effect.
func CookSipMsg(s string) string {
	var buf bytes.Buffer
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		buf.WriteString(line)
		if i == len(lines)-1 {
			break //on last line
		}
		if strings.HasSuffix(line, "\r") {
			buf.WriteString("\n")
		} else {
			buf.WriteString("\r\n")
		}
	}

	return buf.String()
}

//temporarily solution
func UUID() string {
	uuidNum++
	return fmt.Sprint(uuidNum)
}
