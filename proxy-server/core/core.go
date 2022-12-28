package core

import (
	"fmt"
	"net"
	"sync"
	"tcp-proxy/modules/log"
	"tcp-proxy/modules/message"
	"tcp-proxy/modules/tcp"
)

var lock = sync.Mutex{}

// 客户端注册
func registerClient(conn net.Conn, name string, port uint16) {
	defer conn.Close()
	md5Name, bMd5Name := GetMD5(name)
	var respData message.DataInfo
	// 进行服务注册
	if _, ok := ProxyPoolMap[md5Name]; ok {
		respData.Msg = fmt.Sprintf("服务名-%v 已注册，不能重复注册", name)
		respData.Method = message.TypeResponse
		respData.Write(conn)
		log.Println(fmt.Sprintf("[-]接入失败,服务名重复! 地址: %v 名称: %v 代理端口: %v", conn.RemoteAddr(), name, port))
		return
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil {
		respData.Msg = fmt.Sprintf("注册失败! %v", err)
		respData.Method = message.TypeResponse
		respData.Write(conn)
		log.Println(fmt.Sprintf("地址 %v 注册失败! %v", conn.RemoteAddr(), err))
		return
	}

	defer func() {
		log.Println(fmt.Sprintf("[-]断开连接! 地址: %v 名称: %v 代理端口: %v", conn.RemoteAddr(), name, port))
		ProxyPoolClose(md5Name)
		listener.Close()
	}()
	lock.Lock()
	ProxyPoolMap[md5Name] = NewProxyPool()
	lock.Unlock()
	respData = message.DataInfo{
		Method: message.TypeResponse,
		Status: 1,
		Msg:    fmt.Sprintf("服务-%v 注册成功", name),
	}
	respData.Write(conn)
	log.Println(fmt.Sprintf("[+]注册成功! 地址: %v 名称: %v 代理端口: %v", conn.RemoteAddr(), name, port))
	go func() {
		for {
			reqConn, err := listener.Accept()
			if err != nil {
				return
			}
			go func() {
				if pool, ok := ProxyPoolMap[md5Name]; ok {
					// 注册一个连接
					proxyConn := pool.CreateProxyConn(reqConn)
					err = message.ConnInfo{
						Method: message.TypeRequest,
						Name:   bMd5Name,
						ConnId: proxyConn.ID,
					}.Write(conn)
					if err != nil {
						return
					}
				}
			}()
		}
	}()
	// 检测是否断开连接

	buffer := make([]byte, 1024)
	for {
		_, err = conn.Read(buffer)
		if err != nil {
			return
		}
	}

}

// 建立数据通信
func handleClient(conn net.Conn, name string, id uint32) {
	pool, ok := ProxyPoolMap[name]
	if !ok {
		return
	}
	proxyConn := pool.GetProxyConn(id)
	if proxyConn.ID == 0 {
		return
	}
	defer pool.CloseProxyConn(id)
	// 返回建立成功
	message.DataInfo{
		Method: message.TypeResponse,
		Status: 1,
	}.Write(conn)
	tcp.DataExchange(proxyConn.Conn, conn)
}

func Start(conn net.Conn) {

	// 获取客户端返回的信息
	var recvData message.DataInfo
	err := recvData.Read(conn)
	if err != nil || recvData.Status == 0 {
		conn.Close()
		return
	}
	if recvData.Method == message.TypeRegister && recvData.Status != 0 {
		registerClient(conn, recvData.Name, recvData.Port)
	} else if recvData.Method == message.TypeRequest && recvData.Status != 0 {
		handleClient(conn, recvData.Name, recvData.ConnId)
	} else {
		conn.Close()
	}

}
