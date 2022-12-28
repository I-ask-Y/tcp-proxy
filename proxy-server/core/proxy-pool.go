package core

import (
	"crypto/md5"
	"fmt"
	"net"
	"sync"
	"time"
)

type ProxyConn struct {
	ID       uint32
	Conn     net.Conn
	TimeUnix time.Time
}

type ProxyPool struct {
	Seq  *Sequencer
	data map[uint32]ProxyConn
	m    sync.Mutex
}

var ProxyPoolMap = map[string]ProxyPool{}

// NewProxyPool 创建新的代理池
func NewProxyPool() ProxyPool {
	return ProxyPool{
		Seq:  NewSequencer(),
		data: map[uint32]ProxyConn{},
	}
}

func (p *ProxyPool) CreateProxyConn(conn net.Conn) ProxyConn {
	p.m.Lock()
	defer p.m.Unlock()
	proxyConn := ProxyConn{
		ID:       p.Seq.Next(),
		Conn:     conn,
		TimeUnix: time.Now(),
	}
	if proxyConn.ID == 0 {
		proxyConn.ID = p.Seq.Next()
	}
	p.data[proxyConn.ID] = proxyConn
	return proxyConn
}

func (p *ProxyPool) GetProxyConn(id uint32) ProxyConn {
	p.m.Lock()
	defer p.m.Unlock()

	if pool, ok := p.data[id]; ok {
		return pool
	}
	return ProxyConn{}
}

func (p *ProxyPool) CloseProxyConn(id uint32) {
	p.m.Lock()
	defer p.m.Unlock()
	if pool, ok := p.data[id]; ok {
		pool.Conn.Close()
		delete(p.data, id)
	}
}

func ProxyPoolClose(name string) {
	if _, ok := ProxyPoolMap[name]; !ok {
		return
	}
	for id, item := range ProxyPoolMap[name].data {
		_ = item.Conn.Close()
		delete(ProxyPoolMap[name].data, id)
	}
	delete(ProxyPoolMap, name)
}

func GetMD5(s string) (string, [16]byte) {
	data := []byte(s)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has), has
}
