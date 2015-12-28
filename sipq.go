// sipq.go
package main

import (
	"flag"
	"fmt"
	"net"
	"sipq/transport"
	"sipq/util"
)

func udpServer() {
	fmt.Println("udp server start")
	conn, err := transport.CreateUdpConnection(transport.ServerAddress)
	util.ErrorPanic(err)

	fmt.Println("local addr:", conn.Conn.LocalAddr())
	var buffer []byte = make([]byte, 1024)
	for {

		n, addr, err := conn.ReadFrom(buffer)

		util.ErrorPanic(err)

		fmt.Println("udp server received:", string(buffer[0:n]))

		fmt.Println("remote addr:", addr)

		_, err = conn.WriteTo(buffer[0:n], addr)
		util.ErrorPanic(err)

	}
	fmt.Println("udp server end")
}

func udpClient() {
	fmt.Println("udp client start")
	conn, err := transport.CreateUdpConnection(transport.ClientAddress)
	util.ErrorPanic(err)

	serverAddr, err := net.ResolveUDPAddr(transport.TransportTypeUDP.String(), transport.ServerAddress)
	util.ErrorPanic(err)
	var buf []byte = make([]byte, 1024)
	for i := 0; i < 3; i++ {
		_, err := conn.WriteTo([]byte(fmt.Sprintf("sent %d\n", i)), serverAddr)

		util.ErrorPanic(err)

		n, _, err := conn.ReadFrom(buf)
		util.ErrorPanic(err)
		fmt.Println("udp client received:", string(buf[0:n]))
		util.ErrorPanic(err)

	}
	fmt.Println("udp client end")
}

func main() {
	var startServer *bool = flag.Bool("server", true, "start a server")
	flag.Parse()

	if *startServer {
		udpServer()
	} else {
		udpClient()
	}

}
