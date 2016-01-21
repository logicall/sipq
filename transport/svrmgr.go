package transport

import (
	"github.com/henryscala/sipq/trace"

	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/container/concurrent"
)

var allServers *concurrent.List = concurrent.NewList()

//start the server and keep the server in all server list
func StartServer(ip string, port int, transportType TransportType) error {
	addr := util.AddrStr(ip, port)
	trace.Trace.Println("starting server", transportType, addr)
	var server *Server
	var err error
	switch transportType {
	case TCP:
		server, err = createTcpServer(addr)
		go handleNewConn(server)

	case UDP:
		server, err = CreateUdpServer(addr)

	}

	if err != nil {
		return err
	}
	allServers.Add(server)
	return nil
}
