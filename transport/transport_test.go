// transport_test.go
package transport

import (
	"bufio"

	"fmt"
	"sipq/util"
	"testing"
)

func udpServerConn() *Connection {
	conn, err := CreateUdpConnection(ServerAddress)

	util.ErrorPanic(err)
	return conn
}

func udpClientConn() *Connection {
	conn, err := CreateUdpConnection(ClientAddress)

	util.ErrorPanic(err)
	return conn
}

func tcpServer() *Server {

	svr, err := CreateTcpServer(ServerAddress)

	util.ErrorPanic(err)
	return svr
}

func tcpserverHandleConnection(svr *Server, outstr chan string) {

	conn, err := svr.Accept()

	util.ErrorPanic(err)

	var reader *bufio.Reader
	reader = conn.Reader()

	str, _ := reader.ReadString('\n')

	outstr <- str
}

func tcpClient() *Connection {

	clt, err := CreateTcpConnection(ServerAddress)
	util.ErrorPanic(err)
	return clt

}

func udpServerHandleData(conn *Connection, outstr chan string) {
	fmt.Println("udpServerHandleData")
	var buf []byte = make([]byte, 1024)

	n, _, err := conn.ReadFrom(buf)
	util.ErrorPanic(err)
	fmt.Println("udpServerHandleData received:", string(buf[0:n]))

	outstr <- string(buf[0:n])
}

func TestUdpConn(t *testing.T) {
	t.Log("TestUdpConn start ")
	svrConn := udpServerConn()
	output := make(chan string)

	go udpServerHandleData(svrConn, output)

	cltConn := udpClientConn()

	expectedStr := "hello\n"

	fmt.Println("server conn addr:", svrConn.Conn.LocalAddr())
	n, err := cltConn.WriteTo([]byte(expectedStr), svrConn.Conn.LocalAddr())
	util.ErrorPanic(err)
	fmt.Println("finished writeTo ", n)

	var outputStr string
	outputStr = <-output

	cltConn.Close()

	svrConn.Close()

	if expectedStr != outputStr {
		t.Log("expectedStr:(", expectedStr, ")")
		t.Log("outputStr:(", outputStr, ")")
		t.Fail()
	}
}

func TestTcpConn(t *testing.T) {

	svr := tcpServer()

	output := make(chan string)

	go tcpserverHandleConnection(svr, output)

	clt := tcpClient()

	writer := clt.Writer()

	expectedStr := "hello\n"

	writer.WriteString(expectedStr)
	writer.Flush() //flush is mandatory

	var outputStr string
	outputStr = <-output

	clt.Close()

	svr.Close()

	if expectedStr != outputStr {
		t.Log("expectedStr:(", expectedStr, ")")
		t.Log("outputStr:(", outputStr, ")")
		t.Fail()
	}
}
