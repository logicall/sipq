package rawstr

import (
	"strings"
	"testing"
)

func TestPreCompile(t *testing.T) {
	program := `
		a=1;
		b=2;
		str1="abc";
	`
	program += "str2=`line1\nline2\rline3\r\nline4\n`;\n"
	program += "console.log(str2);"

	resultProgram, err := PreCompile(program)
	if err != nil {
		t.Error(err)
	}
	t.Log(program)
	t.Log(resultProgram)
	if strings.Contains(resultProgram, BackQuote) {
		t.Error("not expected")
	}
}
