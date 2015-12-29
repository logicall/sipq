package config

import (
	"testing"
)

func TestDefault(t *testing.T) {
	defaultIP := "127.0.0.1"
	defaultPort := 50600
	expectedType := []string{"udp", "tcp"}

	for i := 0; i < 2; i++ {
		if TheExeConfig.Server[i].Ip != defaultIP {
			t.Errorf("%s!=%s", TheExeConfig.Server[i].Ip, defaultIP)
		}
		if TheExeConfig.Server[i].Port != defaultPort {
			t.Errorf("%d!=%d", TheExeConfig.Server[i].Port, defaultPort)
		}
		if TheExeConfig.Server[i].Type != expectedType[i] {
			t.Errorf("%s!=%s", TheExeConfig.Server[i].Type, expectedType[i])
		}
	}
}
