// transport_test.go
package transport

import (
	"bufio"
	
	
	
	"testing"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var serverAddr string = "127.0.0.1:5060"

func tcpServer() *Server {

	svr, err := CreateTcpServer(serverAddr)

	ErrorPanic(err)
	return svr
}

func tcpserverHandleConnection(svr *Server, outstr chan string) {

	conn, err := svr.Accept()

	ErrorPanic(err)

	var reader *bufio.Reader
	reader = conn.Reader()

	str, _ := reader.ReadString('\n')

	outstr <- str
}

func tcpClient() *Connection {

	clt, err := CreateTcpConnection(serverAddr)
	ErrorPanic(err)
	return clt

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
	outputStr = <- output

	clt.Close()

	svr.Close()

	if expectedStr != outputStr {
		t.Log("expectedStr:(",expectedStr,")")
		t.Log("outputStr:(",outputStr,")")
		t.Fail()
	}
}
