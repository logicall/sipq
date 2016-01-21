package util

import (
	"strings"
	"testing"
)

func TestCookSipMsg(t *testing.T) {
	rawStr := `
line1
line2
`
	if strings.Contains(rawStr, "\r") {
		t.Error("not expected")
	}
	if !strings.Contains(rawStr, "\n") {
		t.Error("not expected")
	}
	cookedStr := CookSipMsg(rawStr)
	if !strings.Contains(cookedStr, "\r\n") {
		t.Error("not expected")
	}

	rawStr = "line3\nline4\n"
	strs := strings.Split(rawStr, "\n")
	if len(strs) != 3 {
		t.Errorf("%#v", strs)
	}

	rawStr = "line5\r\nline6"
	cookedStr = CookSipMsg(rawStr)
	if rawStr != cookedStr {
		t.Error("not expected", cookedStr)
	}
}

func TestFindFreePort(t *testing.T) {
	//for udp
	t.Log("find free port for udp")
	ports, err := FindFreePort("udp", "127.0.0.1", 2)
	if err != nil {
		t.Error(err)
	}
	if ports[0] == 0 || ports[1] == 0 {
		t.Error(ports)
	}
	if ports[0] == ports[1] {
		t.Error(ports)
	}

	//fot tcp
	t.Log("find free port for tcp")
	ports, err = FindFreePort("tcp", "127.0.0.1", 2)
	if err != nil {
		t.Error(err)
	}
	if ports[0] == 0 || ports[1] == 0 {
		t.Error(ports)
	}
	if ports[0] == ports[1] {
		t.Error(ports)
	}

}
