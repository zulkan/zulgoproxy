package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"log"
	"net/http"
)

func filterIP(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	log.Println("host: ", req.RemoteAddr)

	return req, nil
}

type handleConnect struct {
}

func (h *handleConnect) HandleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
	log.Println("CONNECT REQ", host, " FROM ", ctx.Req.RemoteAddr)

	return goproxy.OkConnect, host
}

func getHandleConnect() goproxy.HttpsHandler {
	return &handleConnect{}
}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	auth.ProxyBasic(proxy, "Bearer", func(user, passwd string) bool {
		if user == "zulkan" || passwd == "zulkan" {
			return true
		}
		return false
	})
	proxy.OnRequest().DoFunc(filterIP)
	proxy.OnRequest().HandleConnect(getHandleConnect())

	log.Fatal(http.ListenAndServe(":8181", proxy))
}
