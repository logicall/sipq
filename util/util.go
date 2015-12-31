package util

import (
	"sipq/trace"
	"strings"
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

//replace \n with \r\n
//might need to be refactored for messages with "\r\n"
func CookSipMsg(s string) string {
	return strings.Replace(s, "\n", "\r\n", -1)
}
