/*
	网络会话层：
	1.不关心底层通信实现
	2.统计网络流量
	3.TODO::限流(令牌桶)
*/
package session

import (
	"bytes"
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/session/internal/secert"
	"github.com/nerored/chat-test-golang/net/socket"
)

type SessionState int

const (
	SESSION_STAT_INIT SessionState = iota
	SESSION_STAT_HANDSHAKING
	SESSION_STAT_ESTABLISHED
	SESSION_STAT_DISCONNECTED
)

type Session struct {
	context.Context

	SessionState

	recvChan   chan []byte
	blockCrypt secert.BlockCrypt

	sndPkgCnt int64
	rcvPkgCnt int64
	sndPkgLen int64
	rcvPkgLen int64
	ConnectAt time.Time
}

func NewSession(ctx context.Context) *Session {
	if ctx == nil {
		log.Erro("[session] new session failed, context is nil")
		return nil
	}

	return &Session{
		Context:      ctx,
		recvChan:     make(chan []byte),
		SessionState: SESSION_STAT_INIT,
	}
}

func (s *Session) State() SessionState {
	return s.SessionState
}

func (s *Session) client() *socket.Client {
	if s.Context == nil {
		return nil
	}

	param := s.Context.Value(socket.SERVER_CTX_PARAM_CONN)

	if param == nil {
		return nil
	}

	return param.(*socket.Client)
}

func (s *Session) IsConnected() bool {
	if s.Context == nil || s.SessionState != SESSION_STAT_ESTABLISHED {
		return false
	}

	client := s.client()

	return client != nil && client.IsAlive()
}

func (s *Session) IsEncrypted() bool {
	return s.blockCrypt != nil
}

func (s *Session) Disconnect() {
	if !s.IsConnected() {
		return
	}

	client := s.client()

	if client == nil {
		return
	}

	client.Close()
	s.SessionState = SESSION_STAT_DISCONNECTED
}

func (s *Session) WaitHandshake() (ok bool) {
	if s.SessionState != SESSION_STAT_INIT {
		log.Erro("[session] wait for hands shake failed,session state err %v",
			s.SessionState)
		return
	}

	client := s.client()

	if client == nil || !client.IsAlive() {
		log.Erro("[session] wait for hands shake failed,no set alived client")
		return
	}

	go client.Recv(s.recvChan)

	ok = s.handshakeServerSide() && s.IsConnected()

	if !ok {
		s.Disconnect()
	}

	return
}

func (s *Session) TryConnect(network, addr string, needEncryption bool) (ok bool) {
	if network == "" || addr == "" {
		log.Erro("[session] try connect to %v failed, illegal connection info", addr)
		return
	}

	if s.SessionState != SESSION_STAT_INIT {
		log.Erro("[session] try connect to %v failed,state can't be %v",
			addr, s.SessionState)
		return
	}

	client := s.client()

	if client != nil {
		log.Erro("[session] try connect to %v failed,is connected to %v",
			addr, client.RemoteAddr())
		return
	}

	client = socket.DialClient(network, addr)

	if client == nil || !client.IsAlive() {
		log.Erro("[session] try connect to %v failed, dial failed", addr)
		return
	}

	s.Context = context.WithValue(s.Context, socket.SERVER_CTX_PARAM_CONN, client)
	go client.Recv(s.recvChan)

	switch needEncryption {
	case true:
		ok = s.handshakeClientSide1() && s.IsConnected()
	default:
		ok = s.handshakeClientSide2() && s.IsConnected()
	}

	if !ok {
		s.Disconnect()
	}

	return
}

func (s *Session) Send(data []byte) (ok bool) {
	return s.sendWithCompleteConn(data, true)
}

func (s *Session) Recv() ([]byte, error) {
	return s.recvWithCompleteConn(true)
}

func (s *Session) SetWTimeout(delta time.Duration) {
	client := s.client()

	if client == nil {
		return
	}

	client.SetWriteTimeout(delta)
}

func (s *Session) SetRTimeout(delta time.Duration) {
	client := s.client()

	if client == nil {
		return
	}

	client.SetReadTimeout(delta)
}

func (s *Session) Profile() (sndPkgCnt, sndPkgLen, rcvPkgCnt, rcvPkgLen int64) {
	sndPkgCnt = atomic.LoadInt64(&s.sndPkgCnt)
	sndPkgLen = atomic.LoadInt64(&s.sndPkgLen)
	rcvPkgCnt = atomic.LoadInt64(&s.rcvPkgCnt)
	rcvPkgLen = atomic.LoadInt64(&s.rcvPkgLen)

	return
}

// -------------------------------------------------------- inner methods

func (s *Session) sendWithCompleteConn(data []byte, needCompleteConn bool) (ok bool) {
	if needCompleteConn && !s.IsConnected() {
		return
	}

	client := s.client()

	if client == nil || !client.IsAlive() {
		return
	}

	var transportData []byte

	switch s.SessionState {
	case SESSION_STAT_ESTABLISHED:
		if s.blockCrypt != nil {
			transportData = s.encryptData(data)
		} else {
			transportData = data
		}
	case SESSION_STAT_HANDSHAKING:
		transportData = data
	default:
		log.Erro("[session] remote %v session state err %v can't send data",
			client.RemoteAddr(), s.SessionState)
		return
	}

	ok = client.Send(transportData)

	if !ok {
		return
	}

	atomic.AddInt64(&s.sndPkgCnt, 1)
	atomic.AddInt64(&s.sndPkgLen, int64(len(transportData)))
	return
}

func (s *Session) recvWithCompleteConn(needCompleteConn bool) ([]byte, error) {
	if s.recvChan == nil || s.Context == nil {
		return nil, io.ErrClosedPipe
	}

	if needCompleteConn && !s.IsConnected() {
		return nil, io.ErrClosedPipe
	}

	client := s.client()

	if client == nil || !client.IsAlive() {
		return nil, io.ErrClosedPipe
	}

	var recvData []byte

	select {
	case <-s.Context.Done():
		return recvData, s.Context.Err()
	case chanData, ok := <-s.recvChan:
		if !ok {
			return recvData, io.ErrClosedPipe
		}

		recvData = chanData
	}

	switch s.SessionState {
	case SESSION_STAT_ESTABLISHED:
		if s.blockCrypt != nil {
			recvData = s.decryptData(recvData)
		}
	case SESSION_STAT_HANDSHAKING:
	default:
		log.Warn("[session] remote %v recv data on %v error state",
			client.RemoteAddr(), s.SessionState)
	}

	atomic.AddInt64(&s.rcvPkgCnt, 1)
	atomic.AddInt64(&s.rcvPkgLen, int64(len(recvData)))
	return recvData, nil
}

func (s *Session) encryptData(data []byte) (result []byte) {
	if s.SessionState != SESSION_STAT_ESTABLISHED || s.blockCrypt == nil {
		return data
	}

	result = make([]byte, len(data))
	s.blockCrypt.Encrypt(result, data)
	return
}

func (s *Session) decryptData(data []byte) (result []byte) {
	if s.SessionState != SESSION_STAT_ESTABLISHED || s.blockCrypt == nil {
		return data
	}

	result = make([]byte, len(data))
	s.blockCrypt.Decrypt(result, data)
	return
}

func (s *Session) encryptChan(key []byte) (err error) {
	blockCrypt, err := secert.NewAESBlockCrypt(key)

	if err != nil {
		return
	}

	s.blockCrypt = blockCrypt
	return
}

func (s *Session) encodeHandshakeConfimData(switchKey []byte) []byte {
	secertKey := s.encryptData(switchKey)

	var buffer bytes.Buffer
	buffer.Write(switchKey)
	buffer.Write(secertKey)

	return buffer.Bytes()
}

func (s *Session) decodeHandshakeConfimData(data []byte) (switchKey, secertKey []byte) {
	buffer := bytes.NewBuffer(data)
	return buffer.Next(secert.ECDH_KEY_LEN), buffer.Next(secert.ECDH_KEY_LEN)
}

func (s *Session) handshakeServerSide() (success bool) {
	client := s.client()

	if client == nil || !client.IsAlive() {
		log.Erro(`[session-handshake] server side:
step0: remote %v no client obj`,
			client.RemoteAddr())
	}

	if s.SessionState != SESSION_STAT_INIT {
		log.Erro(`[session-handshake] server side:
step0: remote %v session state error %v`,
			client.RemoteAddr(), s.SessionState)
		return
	}

	s.SessionState = SESSION_STAT_HANDSHAKING
	// hands shake process for server side
	//#1 wait secert
	log.Debu(`[session-handshake] server side:
step1: remote %v wait handshake ack`,
		client.RemoteAddr())

	hs_ack, err := s.recvWithCompleteConn(false)

	if err != nil {
		log.Erro(`[session-handshake] server side:
step1: remote %v recv handshake ack failed %v`,
			client.RemoteAddr(), err)
		return
	}

	switch len(hs_ack) {
	//------------------------------------------------ unencrypted channel handshake
	case 0:
		log.Debu(`[session-handshake] server side:
step1: remote %v while running on unencrypted channel`,
			client.RemoteAddr())

		if !s.sendWithCompleteConn(nil, false) {
			log.Erro(`[session-handshake] server side:
step1: remote %v send confirm ack failed`,
				client.RemoteAddr())
			return
		}

		s.SessionState = SESSION_STAT_ESTABLISHED

		log.Debu(`[session-handshake] server side:
step2: remote %v wait addr setting`,
			client.RemoteAddr())

		remoteAddr, err := s.recvWithCompleteConn(false)

		if err != nil {
			log.Erro(`[session-handshake] server side:
step2: remote %v wait addr setting failed %v`,
				client.RemoteAddr(), err)
			return
		}

		_, _, err = net.SplitHostPort(string(remoteAddr))

		if err != nil {
			log.Erro(`[session-handshake] server side:
step3: remote %v parse ip setting failed %v`,
				client.RemoteAddr(), err)
			return
		}

		client.SetRemoteAddr(string(remoteAddr))

		log.Debu(`[session-handshake] server side:
step final: remote %v finished handshake,success on server side`,
			client.RemoteAddr())

	//------------------------------------------------ encryption channel handshake
	case secert.ECDH_KEY_LEN:
		log.Debu(`[session-handshake] server side:
step1: remote %v while running on encryption channel`,
			client.RemoteAddr())

		switchKey, selfKey, _ := secert.ECDHKeyNew()

		publickKey, err := secert.ECDHKeyGen(hs_ack[:32], selfKey)

		if err != nil {
			log.Erro(`[session-handshake] server side:
step1: remote %v secert gen failed`,
				client.RemoteAddr())
			return
		}

		if err := s.encryptChan(publickKey[:]); err != nil {
			log.Erro(`[session-handshake] server side:
step1: remote %v encrypt chann failed`,
				client.RemoteAddr())
			return
		}

		//#2 send server confirm data
		log.Debu(`[session-handshake] server side:
step2: remote %v send confirm ack`,
			client.RemoteAddr())

		if !s.sendWithCompleteConn(s.encodeHandshakeConfimData(switchKey), false) {
			log.Erro(`[session-handshake] server side:
step2: remote %v send confirm ack`,
				client.RemoteAddr())
			return
		}

		s.SessionState = SESSION_STAT_ESTABLISHED

		log.Debu(`[session-handshake] server side:
step3: remote %v wait addr setting`,
			client.RemoteAddr())

		remoteAddr, err := s.recvWithCompleteConn(false)

		if err != nil {
			log.Erro(`[session-handshake] server side:
step3: remote %v wait addr setting failed %v`,
				client.RemoteAddr(), err)
			return
		}

		_, _, err = net.SplitHostPort(string(remoteAddr))

		if err != nil {
			log.Erro(`[session-handshake] server side:
step3: remote %v parse ip setting failed %v`,
				client.RemoteAddr(), err)
			return
		}

		client.SetRemoteAddr(string(remoteAddr))

		log.Debu(`[session-handshake] server side:
step final: remote %v finished handshake,success on server side`,
			client.RemoteAddr())

	//------------------------------------------------ unknown handshake
	default:
		log.Erro(`[session-handshake] server side:
step1: remote %v recv handshake ack failed,unknown len %v`,
			client.RemoteAddr(), len(hs_ack))
		return
	}

	return true
}

func (s *Session) handshakeClientSide1() (ok bool) {
	client := s.client()

	if client == nil || !client.IsAlive() {
		log.Erro(`[session-handshake] client side:
step0: remote %v no client obj`,
			client.RemoteAddr())
	}

	if s.SessionState != SESSION_STAT_INIT {
		log.Erro(`[session-handshake] client side:
step0: remote %v session state error %v`,
			client.RemoteAddr(), s.SessionState)
		return
	}

	s.SessionState = SESSION_STAT_HANDSHAKING
	// hands shake process for client side
	//#1 send swithcKey
	log.Debu(`[session-handshake] client side:
step1: remote %v send handshake ack for encryption channel`,
		client.RemoteAddr())

	switchKey, selfKey, _ := secert.ECDHKeyNew()

	if !s.sendWithCompleteConn(switchKey[:], false) {
		log.Erro(`[session-handshake] client side:
step1: remote %v send handshake ack failed`,
			client.RemoteAddr())
		return
	}

	//#2 wait server confirm data
	log.Debu(`[session-handshake] client side:
step2: remote %v wait confirm ack`,
		client.RemoteAddr())

	confirmData, err := s.recvWithCompleteConn(false)

	if err != nil {
		log.Erro(`[session-handshake] client side:
step2: remote %v wait confirm ack failed %v`,
			client.RemoteAddr(), err)
		return
	}

	remoteKey, secertkey := s.decodeHandshakeConfimData(confirmData)

	publickKey, err := secert.ECDHKeyGen(remoteKey, selfKey)

	if err != nil {
		log.Erro(`[session-handshake] client side:
step2: remote %v secert gen failed %v`,
			client.RemoteAddr(), err)
		return
	}

	if err := s.encryptChan(publickKey[:]); err != nil {
		log.Erro(`[session-handshake] client side:
step2: remote %v encrypt chann failed`,
			client.RemoteAddr())
		return
	}

	if !bytes.Equal(s.decryptData(secertkey[:]), remoteKey[:]) {
		log.Erro(`[session-handshake] client side:
step2: remote %v confirm failed`,
			client.RemoteAddr())
		return
	}

	s.SessionState = SESSION_STAT_ESTABLISHED

	//#3 send self ipaddr to server
	log.Debu(`[session-handshake] client side:
step3: remote %v send self ipaddr to server`,
		client.RemoteAddr())

	if !s.sendWithCompleteConn([]byte(client.LocalAddr()), false) {
		log.Erro(`[session-handshake] client side:
step3: remote %v send self ipaddr to server failed`,
			client.RemoteAddr())
		return
	}

	log.Debu(`[session-handshake] client side:
step final: remote %v finished handshake,success on client side`,
		client.RemoteAddr())

	return true
}

func (s *Session) handshakeClientSide2() (ok bool) {
	client := s.client()

	if client == nil || !client.IsAlive() {
		log.Erro(`[session-handshake] client side:
step0: remote %v no client obj`,
			client.RemoteAddr())
	}

	if s.SessionState != SESSION_STAT_INIT {
		log.Erro(`[session-handshake] client side:
step0: remote %v session state error %v`,
			client.RemoteAddr(), s.SessionState)
		return
	}

	s.SessionState = SESSION_STAT_HANDSHAKING
	// hands shake process for client side
	//#1 send empty packet
	log.Debu(`[session-handshake] client side:
step1: remote %v send handshake ack for unencrypted channel`,
		client.RemoteAddr())

	if !s.sendWithCompleteConn(nil, false) {
		log.Erro(`[session-handshake] client side:
step1: remote %v send handshake ack failed`,
			client.RemoteAddr())
		return
	}

	//#2 wait server confirm data(also empty)
	log.Debu(`[session-handshake] client side:
step2: remote %v wait confirm ack`,
		client.RemoteAddr())

	confirmData, err := s.recvWithCompleteConn(false)

	if err != nil || len(confirmData) != 0 {
		log.Erro(`[session-handshake] client side:
step2: remote %v wait confirm ack failed %v`,
			client.RemoteAddr(), err)
		return
	}

	s.SessionState = SESSION_STAT_ESTABLISHED

	//#3 send self ipaddr to server
	log.Debu(`[session-handshake] client side:
step3: remote %v send self ipaddr to server`,
		client.RemoteAddr())

	if !s.sendWithCompleteConn([]byte(client.LocalAddr()), false) {
		log.Erro(`[session-handshake] client side:
step3: remote %v send self ipaddr to server failed`,
			client.RemoteAddr())
		return
	}

	log.Debu(`[session-handshake] client side:
stepf: remote %v finished handshake,success on client side`,
		client.RemoteAddr())

	return true
}
