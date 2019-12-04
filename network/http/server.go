package http

import (
	"github.com/yaice-rx/yaice/network"
	"github.com/yaice-rx/yaice/resource"
	router_ "github.com/yaice-rx/yaice/router"
	"net/http"
	"strconv"
	"sync"
)

type httpServer struct {
	sync.Mutex
	server  *http.Server
	network string
}

var HttpServerMgr = newHttpServer()

func newHttpServer() network.IServer {
	srv := &httpServer{
		network: "http",
	}
	srv.server = &http.Server{}
	return nil
}

func (this *httpServer) GetNetwork() string {
	return this.network
}

func (this *httpServer) Start() (int, error) {
	mux := http.NewServeMux()
	for router, handler := range router_.RouterMgr.GetHttpHandlerMap() {
		mux.HandleFunc(router, handler)
	}
	this.server.Handler = mux
	port := resource.ServiceResMgr.HttpPort
	//开启监听
	this.server.Addr = ":" + strconv.Itoa(port)
	//开启http网络
	go this.server.ListenAndServe()
	return port, nil
}

func (this *httpServer) Close() {
	this.server.Close()
}
