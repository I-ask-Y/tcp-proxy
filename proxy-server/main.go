package main

import (
	"fmt"
	"net"
	"tcp-proxy/modules/log"
	"tcp-proxy/proxy-server/config"
	"tcp-proxy/proxy-server/core"
)

func main() {
	// 分为客户端注册部分 / 端口发送部分
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", config.Config.Host, config.Config.Port))
	// 属于注册端口
	if err != nil {
		log.Fatalf("服务启动失败：%s", err.Error())
		return
	}

	log.Println("proxy 服务端启动成功，等待客户端连接")
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Println(fmt.Sprintf("地址%v, 连接失败:%v", clientConn.RemoteAddr(), err.Error()))
		}
		// 处理每一个端口l
		go core.Start(clientConn)
	}
}
