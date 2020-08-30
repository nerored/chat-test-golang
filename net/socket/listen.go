/*
	监听器封装
*/
package socket

import (
	"net"
)

type Listener struct {
	net.Listener
	acceptHandle func(conn net.Conn)
}

func NewTCPListener(listener net.Listener) *Listener {
	return &Listener{
		Listener: listener,
	}
}
