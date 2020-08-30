/*
	添加了close状态判断
*/
package socket

import (
	"io"
	"net"
	"time"

	"github.com/nerored/chat-test-golang/log"
)

type Socket struct {
	conn net.Conn
	stop DoOnce
}

func dial(network string, remote string) (s *Socket) {
	switch network {
	case "tcp":
		conn, err := net.Dial(network, remote)

		if err != nil {
			log.Erro("[socket tcp] dial %v failed err %v", remote, err)
			return
		}

		log.Info("[socket tcp] connected to remote %v", conn.RemoteAddr())

		return &Socket{
			conn: conn,
		}
	default:
		log.Erro("[socket] dial %v unknow network %v", remote, network)
		return nil
	}
}

//----------------   net.Conn interface

func (s *Socket) Read(b []byte) (n int, err error) {
	if !s.IsAlive() {
		err = io.EOF
		return
	}

	return s.conn.Read(b)
}

func (s *Socket) Write(b []byte) (n int, err error) {
	if !s.IsAlive() {
		err = io.EOF
		return
	}

	return s.conn.Write(b)
}

func (s *Socket) Close(reason string) (err error) {
	if !s.IsAlive() {
		return
	}

	s.stop.Do(func() {
		err = s.conn.Close()
		log.Info("[socket] disconnected to remote %v,reason %v", s.RemoteAddr(), reason)
	})

	return
}

func (s *Socket) IsAlive() (ok bool) {
	if s.conn == nil || s.stop.IsDone() {
		return
	}

	return true
}

func (s *Socket) LocalAddr() net.Addr {
	if !s.IsAlive() {
		log.Erro("[socket] get local addr failed, is not connected")
		return nil
	}

	return s.conn.LocalAddr()
}

func (s *Socket) RemoteAddr() net.Addr {
	if !s.IsAlive() {
		log.Erro("[socket] get remote addr failed, is not connected")
		return nil
	}

	return s.conn.RemoteAddr()
}

func (s *Socket) SetDeadline(t time.Time) (err error) {
	if !s.IsAlive() {
		err = io.EOF
		return
	}

	return s.conn.SetDeadline(t)
}

func (s *Socket) SetReadDeadline(t time.Time) (err error) {
	if !s.IsAlive() {
		err = io.EOF
		return
	}

	return s.conn.SetReadDeadline(t)
}

func (s *Socket) SetWriteDeadline(t time.Time) (err error) {
	if !s.IsAlive() {
		err = io.EOF
		return
	}

	return s.conn.SetWriteDeadline(t)
}
