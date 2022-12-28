package main

import (
	"tcp-proxy/modules/log"
	. "tcp-proxy/proxy-client/config"
	"tcp-proxy/proxy-client/core"
)

func main() {
	// 对当前需要代理的端口进行遍历
	for _, proxyItem := range Config.Proxy {
		client := &core.Client{
			Name:       proxyItem.Name,
			ServerAddr: Config.ServerAddr,
			ProxyAddr:  proxyItem.ProxyAddr,
			RemotePort: proxyItem.RemotePort,
		}
		client.Register()
	}

	log.Println("代理客户端启动成功！")
	c := make(chan bool)
	<-c
}
