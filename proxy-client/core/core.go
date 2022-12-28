package core

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"tcp-proxy/modules/log"
	"tcp-proxy/modules/message"
	"tcp-proxy/modules/tcp"
	"time"
)

type Client struct {
	Name       string // 名称
	ServerAddr string // 远程地址
	ProxyAddr  string // 本地需要代理的地址
	RemotePort uint16 // 远程开放的地址
}

func (c *Client) register() error {
	// 向服务端进行注册服务
	conn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	//   服务端写入注册信息
	err = message.DataInfo{
		Method: message.TypeRegister,
		Name:   c.Name,
		Port:   c.RemotePort,
		Status: 1,
	}.Write(conn)
	if err != nil {
		return err
	}

	// 获取服务端返回的信息，判断是否注册成功
	var respMessage message.DataInfo
	respMessage.Read(conn)

	if respMessage.Status != 0 && respMessage.Method == message.TypeResponse {
		log.Println(respMessage.Msg) //
	} else {
		return errors.New(respMessage.Msg)
	}

	for {
		var sMessage message.ConnInfo
		err := sMessage.Read(conn)
		if err != nil {
			return err
		}
		if sMessage.Method == message.TypeRequest { // 进行注册
			go c.clientHandle(sMessage.Name, sMessage.ConnId)
		}
		runtime.GC()
	}
}

func (c *Client) clientHandle(name [16]byte, connId uint32) {
	if connId == 0 {
		return
	}
	proxyConn, err := net.Dial("tcp", c.ProxyAddr)
	if err != nil {
		return
	}

	// 连接服务端
	sConn, err := net.Dial("tcp", c.ServerAddr)
	if err != nil {
		return
	}
	// 数据请求建立连接通道
	message.DataInfo{
		Method: message.TypeRequest,
		ConnId: connId,
		Name:   fmt.Sprintf("%x", name),
		Status: 1,
	}.Write(sConn)
	// 是否请求建立成功
	var respMessage message.DataInfo
	err = respMessage.Read(sConn)
	if err != nil || respMessage.Status == 0 || respMessage.Method != message.TypeResponse {
		return
	}
	tcp.DataExchange(sConn, proxyConn)
	runtime.GC()
}

func (c *Client) Register() {
	go func() {
		for {
			err := c.register()
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Second * 10) // 连接断开后 隔10s重新注册
		}
	}()
}
