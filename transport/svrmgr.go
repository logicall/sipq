package transport

import (
	"fmt"

	"github.com/henryscala/sipq/config"
	"github.com/henryscala/sipq/trace"
	"github.com/henryscala/sipq/util"
	"github.com/henryscala/sipq/util/concurrent"
)

var AllServers *concurrent.List = concurrent.NewList()

func StartServers(cfg *config.ExeConfig) *concurrent.List {
	list := concurrent.NewList()
	var server *Server
	var err error
	for _, svrCfg := range cfg.Server {
		addr := fmt.Sprintf("%s:%d", svrCfg.Ip, svrCfg.Port)
		trace.Trace.Println("starting server", svrCfg.Type, addr)
		switch svrCfg.Type {
		case TCP.String():
			server, err = CreateTcpServer(addr)
			go handleNewConn(server)

		case UDP.String():
			server, err = CreateUdpServer(addr)

		}
		util.ErrorPanic(err)
		list.Add(server)
	}
	return list
}
