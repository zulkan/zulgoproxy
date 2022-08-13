package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
	"log"
	"net/http"
	"strings"
)

func filterIP(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

	log.Println("host: ", req.RemoteAddr)

	return req, nil
}

func zulUserPass(user, passwd string) bool {
	if user == "zulkan" || passwd == "zulkan" {
		return true
	}
	return false
}

type handleConnect struct {
}

func getHandleConnect() goproxy.HttpsHandler {
	return goproxy.FuncHttpsHandler(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		log.Println("CONNECT REQ", host, " FROM ", ctx.Req.RemoteAddr)

		if !strings.HasPrefix(ctx.Req.RemoteAddr, "140.213") {
			return auth.BasicConnect("realm", zulUserPass).HandleConnect(host, ctx)
		}
		return goproxy.OkConnect, host
	})
}

//func getHandleConnect() goproxy.HttpsHandler {
//	return &handleConnect{}
//}

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(filterIP)
	proxy.OnRequest().HandleConnect(getHandleConnect())

	log.Fatal(http.ListenAndServe(":8181", proxy))
}
