package util

import (
	"bytes"
	"fmt"

	"net"
	"strconv"
	"strings"

	"github.com/henryscala/sipq/trace"
)

var (
	uuidNum int
)

func ErrorPanic(err error) {
	if err != nil {
		trace.Error(err)

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

func AddrStr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}

func AddrStrSplit(addr string) (ip string, port int, err error) {
	strs := strings.Split(addr, ":")
	//assume the last segment is port
	portStr := strs[len(strs)-1]
	port, err = strconv.Atoi(portStr)

	ip = strings.Join(strs[:len(strs)-1], ":")
	return ip, port, err
}

func Addr(ip string, port int, transportType string) (net.Addr, error) {
	transportType = strings.ToLower(transportType)
	switch {
	case strings.Contains(transportType, "udp"):
		return net.ResolveUDPAddr("udp", AddrStr(ip, port))
	case strings.Contains(transportType, "tcp"):
		return net.ResolveTCPAddr("tcp", AddrStr(ip, port))
	}

	return nil, fmt.Errorf("not implemented")
}

//find number of free port;
//user of this function should call this function once;
//multiple call may result in repeated value
func FindFreePort(transportType string, ip string, num int) ([]int, error) {
	trace.Trace("enter FindFreePort")
	defer trace.Trace("exit FindFreePort")
	var result []int = make([]int, num)
	transportType = strings.ToLower(transportType)

	switch {
	case strings.Contains(transportType, "udp"):

		for i := 0; i < num; i++ {
			addr, err := net.ResolveUDPAddr("udp", AddrStr(ip, 0))
			if err != nil {
				return nil, err
			}
			udpConn, err := net.ListenUDP("udp", addr)
			if err != nil {
				return nil, err
			}
			defer udpConn.Close()
			_, port, err := AddrStrSplit(udpConn.LocalAddr().String())

			if err != nil {
				return nil, err
			}
			trace.Debug("find port", transportType, i, "-", port)
			result[i] = port
		}

		return result, nil

	case strings.Contains(transportType, "tcp"):

		for i := 0; i < num; i++ {
			addr, err := net.ResolveTCPAddr("tcp", AddrStr(ip, 0))
			if err != nil {
				return nil, err
			}
			tcpConn, err := net.ListenTCP("tcp", addr)
			if err != nil {
				return nil, err
			}
			defer tcpConn.Close()

			_, port, err := AddrStrSplit(tcpConn.Addr().String())
			if err != nil {
				return nil, err
			}
			trace.Debug("find port", transportType, i, "-", port)
			result[i] = port
		}

		return result, nil
	default:

		return nil, fmt.Errorf("not implemented")
	}
}
