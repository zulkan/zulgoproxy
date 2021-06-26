package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func filterIP(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	log.Println("host: ", req.RemoteAddr)

	return req, nil
}

type handleConnect struct {
}

func (h handleConnect) HandleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	log.Println("CONNECT REQ", host, " FROM ", ctx.Req.RemoteAddr)

	return goproxy.OkConnect, host
}

func getHandleConnect() goproxy.HttpsHandler {
	return handleConnect{}
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(filterIP)
	proxy.OnRequest().HandleConnect(getHandleConnect())

	log.Fatal(http.ListenAndServe(":8181", proxy))
}