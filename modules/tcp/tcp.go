package tcp

import (
	"context"
	"net"
)

// DataExchange 数据交换
func DataExchange(conn1, conn2 net.Conn) {
	defer func() {
		conn1.Close()
		conn2.Close()
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()
		buffer := make([]byte, 4096)
		for {
			n, err := conn1.Read(buffer)
			if err != nil {
				return
			}
			n, err = conn2.Write(buffer[:n])
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer cancel()
		buffer := make([]byte, 4096)
		for {
			n, err := conn2.Read(buffer)
			if err != nil {
				return
			}
			n, err = conn1.Write(buffer[:n])
			if err != nil {
				return
			}
		}
	}()
	<-ctx.Done()
}
