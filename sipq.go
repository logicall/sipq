// sipq.go
package main

import (
	"bufio"
	"fmt"
	t "sipq/transport"
	"strconv"
	"sync"
	_ "time"
	"io"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

var serverAddr string = "127.0.0.1:5060"

func server() *t.Server {

	svr, err := t.CreateTcpServer(serverAddr)
	ErrorPanic(err)
	return svr
}

func serverHandleConnection(svr *t.Server) {
	fmt.Println("serverHandleConnection")
	conn, err := svr.Accept()
	
	ErrorPanic(err)
	
	fmt.Println("begin read")
	var reader *bufio.Reader
	reader = conn.Reader()
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break 
		}
		fmt.Println("server received:", str,"and err:",err)
		
	}
	fmt.Println("exit serverHandleConnection")
	wg.Done()
}

func client() *t.Connection {

	clt, err := t.CreateTcpConnection(serverAddr)
	ErrorPanic(err)
	return clt

}

var wg sync.WaitGroup

func main() {
	wg.Add(2)
	fmt.Println("create server")
	svr := server()
	fmt.Println("call handle server connection")
	go serverHandleConnection(svr)
	fmt.Println("call client")
	clt := client()

	fmt.Println("call client writer")
	writer := clt.Writer()
	

	for i := 0; i < 10; i++ {
		n,err:=writer.WriteString("client sent " + strconv.Itoa(i)+"\n")
		writer.Flush()
		fmt.Println("client sent " + strconv.Itoa(i),"n:",n,"err:",err)
	}
	clt.Close()
	
	wg.Done()
	wg.Wait()
	fmt.Println("call close")

	
	svr.Close()
	
	
	
}
