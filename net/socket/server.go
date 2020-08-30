/*
	server
	1.实现监听循环
	2.屏蔽协议细节(tcp ,TODO::udp or ws wss)
*/
package socket

import (
	"context"
	"net"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/socket/utils"
)

type ServerCtxParamKey int

const (
	SERVER_CTX_PARAM_CONN ServerCtxParamKey = 1
)

type ServerEvent int

const (
	SERVER_EVENT_UNKNOW = iota
	SERVER_EVENT_ON_START
	SERVER_EVENT_ON_CLOSE
)

type Server struct {
	network   string
	localAddr string
	utils.WhiteList

	stop     DoOnce
	listener *Listener

	eventHandlers map[ServerEvent][]func()
}

func NewServer(network, localAddr string) *Server {
	return &Server{
		network:       network,
		localAddr:     localAddr,
		eventHandlers: make(map[ServerEvent][]func()),
	}
}

func (s *Server) Close() {
	if s.listener == nil {
		return
	}

	s.stop.Do(func() {
		s.listener.Close()
	})
}

func (s *Server) IsClosed() bool {
	return s.listener == nil || s.stop.IsDone()
}

// not goroutine safe,only use it before listen started
func (s *Server) RegisterEvent(event ServerEvent, handle func()) (ok bool) {
	if s.eventHandlers == nil || handle == nil {
		log.Erro("[socket-server unknown] register event failed")
		return
	}

	s.eventHandlers[event] = append(s.eventHandlers[event], handle)
	return true
}

func (s *Server) FireEvent(event ServerEvent) {
	if s.eventHandlers == nil {
		return
	}

	for _, handle := range s.eventHandlers[event] {
		if handle != nil {
			handle()
		}
	}
}

func (s *Server) StartListen(waitWhenStop time.Duration, recvService func(ctx context.Context)) {
	switch s.network {
	case "tcp":
		tcpListener, err := net.Listen("tcp", s.localAddr)

		if err != nil {
			log.Erro("[socket-server tcp] listen at %v failed,err %v", s.localAddr, err)
			return
		}

		s.listener = NewTCPListener(tcpListener)
	}

	if s.listener == nil {
		log.Erro("[socket-server %v] localAddr %v listen failed,unknow err",
			s.network, s.localAddr)
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Trac("[socket-server %v] at %v,err %v", s.network, s.localAddr, err)
		}

		s.Close()
	}()

	log.Info("[socket-server %v] listen start %v", s.network, s.localAddr)

	serverCtx, cancel := context.WithCancel(context.Background())

	s.FireEvent(SERVER_EVENT_ON_START)

begin:
	for !s.IsClosed() {
		conn, err := s.listener.Accept()

		switch errcheck(err) {
		case error_type_othererro:
			log.Warn("[socket-server %v] at %v,accept err :%v", s.network, s.localAddr, err)
			break begin
		case error_type_temporary:
			pauseService := 5 * time.Second
			log.Warn("[socket-server %v] at %v,accept err :%v, service while recover after %v",
				s.network, s.localAddr, err, pauseService)
			time.Sleep(pauseService)
			continue begin
		}

		if conn == nil {
			continue
		}

		if !s.AccessCheck(conn.RemoteAddr().String()) {
			log.Warn("[socket-server %v] at %v,remote [%v] access deine",
				conn.RemoteAddr(), s.network, s.localAddr)
			conn.Close()
			continue
		}

		if recvService == nil {
			log.Info("[socket-server %v] at %v,accept remote %v recvService nil,so disconnected",
				s.network, s.localAddr, conn.RemoteAddr())
			conn.Close()
			continue
		}

		if s.listener.acceptHandle != nil {
			s.listener.acceptHandle(conn)
		}

		log.Info("[socket-server %v] at %v,accept conn from %v",
			s.network, s.localAddr, conn.RemoteAddr())

		go recvService(context.WithValue(serverCtx, SERVER_CTX_PARAM_CONN, ConnClient(conn)))
	}

	log.Info("[socket-server %v] at %v listen close", s.network, s.localAddr)

	//wait xxx second,then close all connections

	alertTimer := time.NewTimer(time.Second)
	defer alertTimer.Stop()

	for waitSecod := waitWhenStop; waitSecod > 0; waitSecod -= time.Second {
		_, ok := <-alertTimer.C

		if !ok {
			break
		}

		log.Info("[socket-server %v] at %v will close after %v", s.network, s.localAddr, waitSecod)
		alertTimer.Reset(time.Second)
	}

	cancel()
	s.FireEvent(SERVER_EVENT_ON_CLOSE)

	log.Info("[socket-server %v] at %v service are terminated", s.network, s.localAddr)
}
