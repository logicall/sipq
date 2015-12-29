package transport

import (
	"fmt"
	"sipq/config"
	"sipq/trace"
	"sipq/util"
	"sync"
)

type Servers struct {
	serversLock sync.Mutex
	servers     []*Server
}

func (svrs *Servers) Add(svr *Server) {
	svrs.serversLock.Lock()
	svrs.servers = append(svrs.servers, svr)
	svrs.serversLock.Unlock()
}

var AllServers *Servers = &Servers{}

func StartServers(cfg *config.ExeConfig) *Servers {
	var servers []*Server
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
		servers = append(servers, server)
	}
	return &Servers{servers: servers}
}
