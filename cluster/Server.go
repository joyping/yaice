package cluster

import (
	"github.com/yaice-rx/yaice/cluster/ETCDManager"
	"github.com/yaice-rx/yaice/config"
	"github.com/yaice-rx/yaice/rpc"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

type IServer interface {
	Start(conns []string) error
	Set(config config.Config)
	Listen(startPort int, endPort int) int
	Close()
}

type Server struct {
	etcdManger ETCDManager.IEtcdManager
	server     *grpc.Server
}

var ServerMgr = _NewServerMgr()

func _NewServerMgr() IServer {
	return &Server{
		etcdManger: ETCDManager.ServerMgr,
	}
}

func (s *Server) Start(conns []string) error {
	return s.etcdManger.Listen(conns)
}

func (s *Server) Set(config config.Config) {
	s.etcdManger.Set(config.ServerGroup+"+"+config.TypeId+"/"+strconv.Itoa(config.Pid), config)
}

func (s *Server) Listen(startPort int, endPort int) int {
	port := make(chan int)
	defer close(port)
	for i := startPort; i <= endPort; i++ {
		go func() {
			lis, err := net.Listen("tcp", ":"+strconv.Itoa(i))
			if err != nil {
				port <- -1
				return
			}
			port <- i
			server := grpc.NewServer()
			s.server = server
			rpc.RPCMgr.CallServerRPCFunc(server)
			server.Serve(lis)
		}()
		data := <-port
		if data > 0 {
			return data
		}
	}
	return -1
}

func (s *Server) Close() {
	s.server.Stop()
	s.etcdManger.Close()
}
