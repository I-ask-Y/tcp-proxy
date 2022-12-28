package message

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
)

const TypeRegister = 1 // 注册请求
const TypeRequest = 2  // 请求连接
const TypeResponse = 3 // 回应

const ()

type Message interface {
	Write(conn net.Conn) error
	Read(conn net.Conn) error
}

type DataInfo struct {
	Method uint8 `json:"method"` // 数据报文注册还是请求新的连接 request or register
	// 请求状态
	Status uint8  `json:"status"` // 1 代表正确
	Name   string `json:"name"`
	Port   uint16 `json:"port"`
	ConnId uint32 `json:"conn_id"` // 请求的连接id
	Msg    string `json:"msg"`
}

func (m DataInfo) Write(conn net.Conn) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = conn.Write(data)
	return err
}

func (m *DataInfo) Read(conn net.Conn) error {
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data[:n], m)
	if err != nil {
		return err
	}
	return err
}

// ConnInfo 请求建立连接  共21字节
type ConnInfo struct {
	Method uint8    // 占1字节
	Name   [16]byte // 占16 字节
	ConnId uint32   // 占4字节
}

func (c ConnInfo) Write(conn net.Conn) error {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, &c); err != nil {
		return err
	}
	_, err := conn.Write(buf.Bytes())
	return err
}

func (c *ConnInfo) Read(conn net.Conn) error {
	data := make([]byte, 21)
	n, err := conn.Read(data)
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	buf.Write(data[:n])
	err = binary.Read(buf, binary.BigEndian, c)
	return err
}
