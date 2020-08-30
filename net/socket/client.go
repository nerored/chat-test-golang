/*
	client
	1.实现对网络协议的封装（tcp/ TODO::UDP)
	2.包编码，组包和分包
	3.单帧大小2^15,单包最大32M
*/
package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"sync"
	"time"

	"github.com/nerored/chat-test-golang/log"
)

type Client struct {
	// conn
	socket   *Socket
	rd_cache []byte

	//write lock
	sync.Mutex

	// settings
	headerEncoder   binary.ByteOrder
	rTimeout        time.Duration
	wTimeout        time.Duration
	addrSetByRemote string
}

//-----------  new

func DialClient(network string, remote string) *Client {
	socket := dial(network, remote)

	if socket == nil {
		return nil
	}

	return newClient(socket)
}

func ConnClient(conn net.Conn) *Client {
	return newClient(&Socket{
		conn: conn,
	})
}

func newClient(socket *Socket) *Client {
	if socket == nil {
		log.Erro("[socket-client] new client without socket obj")
		return nil
	}

	return &Client{
		socket:        socket,
		rd_cache:      make([]byte, 1024*4),
		headerEncoder: binary.BigEndian,
		rTimeout:      15 * time.Minute,
		wTimeout:      15 * time.Minute,
	}
}

//-----------  aciton

func (c *Client) Close() {
	if !c.IsAlive() {
		return
	}

	c.socket.Close("active disconnect")
}

func (c *Client) IsAlive() (ok bool) {
	return c.socket != nil && c.socket.IsAlive()
}

func (c *Client) LocalAddr() string {
	if !c.IsAlive() {
		return ""
	}

	return c.socket.LocalAddr().String()
}

func (c *Client) RemoteAddr() string {
	if !c.IsAlive() {
		return ""
	}

	if len(c.addrSetByRemote) > 0 {
		return c.addrSetByRemote
	}

	return c.socket.RemoteAddr().String()
}

func (c *Client) Send(data []byte) (ok bool) {
	if c.socket == nil || !c.socket.IsAlive() {
		log.Erro("[session-client] can't send data on a disconnect client")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			if err != nil {
				log.Ulog(log.LOG_LEVEL_WARN, log.PRINT_DEBUG, "%v", err)
			}

			c.socket.Close("send failed")
		}
	}()

	//always make for goroutine safe
	c.Lock()
	defer c.Unlock()

	header := []byte{0, 0}
	for buffer, sndCnt := bytes.NewBuffer(data), 0; buffer.Len() > 0 || sndCnt <= 0; sndCnt++ {
		splitPacket := buffer.Next(PACKET_LIMIT_SIZE)

		err := c.encodeHeader(len(splitPacket), buffer.Len() <= 0, header)

		if err != nil {
			panic(err)
		}

		err = c.socket.SetWriteDeadline(time.Now().Add(c.getWriteTimeout()))

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v write header set deadline failed %v",
				c.RemoteAddr(), err))
		}

		n, err := c.socket.Write(header)

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v send header failed %v",
				c.RemoteAddr(), err))
		}

		if n != len(header) {
			panic(fmt.Errorf("[socket-client] remote %v send header failed, write process %v/%v",
				c.RemoteAddr(), n, len(header)))
		}

		err = c.socket.SetWriteDeadline(time.Now().Add(c.getWriteTimeout()))

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v write body set deadline failed %v",
				c.RemoteAddr(), err))
		}

		n, err = c.socket.Write(splitPacket)

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v send packet failed %v",
				c.RemoteAddr(), err))
		}

		if n != len(splitPacket) {
			panic(fmt.Errorf("[socket-client] remote %v send packet failed, write process %v/%v",
				c.RemoteAddr(), n, len(splitPacket)))
		}
	}

	return true
}

func (c *Client) Recv(receiver chan []byte) {
	if c.socket == nil || !c.socket.IsAlive() {
		log.Erro("[session-client] can't recv data on a disconnect client")
		return
	}

	defer func() {
		if err := recover(); err != nil {
			log.Ulog(log.LOG_LEVEL_WARN, log.PRINT_DEBUG, "%v", err)
		}

		if receiver != nil {
			close(receiver) // stop the recv process outside
		}
		c.socket.Close("recv failed")
	}()

	if receiver == nil {
		log.Erro("[session-client] remote %v no receiver chan", c.RemoteAddr())
		return
	}

	msgBuffer := new(bytes.Buffer)
	for {
		err := c.socket.SetReadDeadline(time.Now().Add(c.getReadTimeout()))

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v read head set deadline failed %v", c.RemoteAddr(), err))
		}

		n, err := io.ReadFull(c.socket, c.rd_cache[:PACKET_HEADER_LEN])

		if err != nil {
			panic(fmt.Errorf("[socket-client] remote %v recv header failed %v", c.RemoteAddr(), err))
		}

		if n != PACKET_HEADER_LEN {
			panic(fmt.Errorf("[socket-client] remote %v recv header failed, recv process %v/%v",
				c.RemoteAddr(), n, PACKET_HEADER_LEN))
		}

		packetLen, isEOF, err := c.decodeHeader(c.rd_cache[:PACKET_HEADER_LEN])

		if err != nil {
			panic(err)
		}

		if packetLen < 0 || packetLen > PACKET_LIMIT_SIZE {
			panic(fmt.Errorf("[socket-client] remote %v recv header failed, illegal package size %v",
				c.RemoteAddr(), packetLen))
		}

		for unreadCnt := packetLen; unreadCnt > 0; {
			readOnce := int(math.Min(float64(unreadCnt), float64(len(c.rd_cache))))

			err := c.socket.SetReadDeadline(time.Now().Add(c.getReadTimeout()))

			if err != nil {
				panic(fmt.Errorf("[socket-client] remote %v read body set deadline failed %v", c.RemoteAddr(), err))
			}

			n, err := io.ReadFull(c.socket, c.rd_cache[:readOnce])

			if err != nil {
				panic(fmt.Errorf("[socket-client] remote %v recv packet failed %v", c.RemoteAddr(), err))
			}

			if n != readOnce {
				panic(fmt.Errorf("[socket-client] remote %v recv packet failed, recv process %v/%v",
					c.RemoteAddr(), n, readOnce))
			}

			if msgBuffer.Len()+n > MESSAGE_LIMIT_SIZE {
				panic(fmt.Errorf("[socket-client] remote %v recv packet failed,cached more than %v bytes",
					c.RemoteAddr(), MESSAGE_LIMIT_SIZE))
			}

			unreadCnt -= n

			msgBuffer.Write(c.rd_cache[:n])
		}

		if !isEOF {
			continue
		}

		receiver <- msgBuffer.Bytes()
		msgBuffer = new(bytes.Buffer)
	}
}

//-----------  settings

func (s *Client) SetRemoteAddr(addr string) {
	if !s.IsAlive() {
		return
	}

	if s.RemoteAddr() == addr {
		return
	}

	s.addrSetByRemote = addr
	log.Info("[session-client] remote addr change %v->%v", s.socket.RemoteAddr(), addr)
}

func (c *Client) getReadTimeout() time.Duration {
	return c.rTimeout
}

func (c *Client) getWriteTimeout() time.Duration {
	return c.wTimeout
}

func (c *Client) SetReadTimeout(t time.Duration) {
	c.rTimeout = t
}

func (c *Client) SetWriteTimeout(t time.Duration) {
	c.wTimeout = t
}

func (c *Client) SetHeaderEncoder(f binary.ByteOrder) {
	switch f {
	case binary.BigEndian, binary.LittleEndian:
	default:
		log.Erro("[socket-client] header encode unsupported byteorder type %v", f)
		return
	}

	c.headerEncoder = f
}

//-----------  helpers

const (
	PACKET_HEADER_LEN  = 2      // 2 bytes header
	PACKET_LIMIT_SIZE  = 0x7FFF // max 32k data per packet
	PACKET_HEADER_EOF  = 0x8000
	MESSAGE_LIMIT_SIZE = 16 * 1024 * 1024 // msg limit 16M
)

func (c *Client) encodeHeader(n int, isEOF bool, header []byte) error {
	if len(header) != PACKET_HEADER_LEN {
		return fmt.Errorf("[socket-client] remote %v pack header faild, header len must be %v",
			c.RemoteAddr(), PACKET_HEADER_LEN)
	}

	if n > PACKET_LIMIT_SIZE || n < 0 {
		return fmt.Errorf("[socket-client] remote %v pack header faild, n is between [0~%v]",
			c.RemoteAddr(), PACKET_LIMIT_SIZE)
	}

	if c.headerEncoder == nil {
		return fmt.Errorf("[socket-client] remote %v pack header faild, headerEncoder not setted",
			c.RemoteAddr())
	}

	cache := uint16(n)

	if isEOF {
		cache |= PACKET_HEADER_EOF
	}

	c.headerEncoder.PutUint16(header[:], cache)
	return nil
}

func (c *Client) decodeHeader(header []byte) (n int, isEOF bool, err error) {
	if len(header) != PACKET_HEADER_LEN {
		err = fmt.Errorf("[socket-client] remote %v unpack header faild, header len must be %v",
			c.RemoteAddr(), PACKET_HEADER_LEN)
		return
	}

	if c.headerEncoder == nil {
		panic(fmt.Errorf("[socket-client] remote %v unpack header faild,  headerEncoder not setted",
			c.RemoteAddr()))
	}

	cache := c.headerEncoder.Uint16(header)

	n = int(cache & ^uint16(PACKET_HEADER_EOF))

	if n < 0 || n > PACKET_LIMIT_SIZE {
		err = fmt.Errorf("[socket-client] remote %v unpack header faild, n is between [0~%v]",
			c.RemoteAddr(), PACKET_LIMIT_SIZE)
		return
	}

	isEOF = cache&PACKET_HEADER_EOF != 0
	return
}
