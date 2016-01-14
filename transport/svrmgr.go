package transport

import (
	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/container/concurrent"
)

var allServers *concurrent.List = concurrent.NewList()

func StartServers(cfg *config.ExeConfig) {
	allServers = concurrent.NewList()

	for _, svrCfg := range cfg.Server {
		StartServer(svrCfg.Ip, svrCfg.Port, TransportType(svrCfg.Type))
	}
}

func StartServer(ip string, port int, transportType TransportType) {
	addr := util.AddrStr(ip, port)
	trace.Trace.Println("starting server", transportType, addr)
	var server *Server
	var err error
	switch transportType {
	case TCP:
		server, err = CreateTcpServer(addr)
		go handleNewConn(server)

	case UDP:
		server, err = CreateUdpServer(addr)

	}
	util.ErrorPanic(err)
	allServers.Add(server)
}
