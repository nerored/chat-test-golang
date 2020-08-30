package session

import (
	"bytes"
	"context"
	"io"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/socket"
)

func BenchmarkTCPPerformance(b *testing.B) {
	PerformanceTest("tcp", false, make([]byte, 100000), b)
}

func BenchmarkTCPPerformanceWithEncryption(b *testing.B) {
	PerformanceTest("tcp", true, make([]byte, 100000), b)
}

func PerformanceTest(network string, needEncryption bool, transData []byte, b *testing.B) {
	log.InitLog("")

	systemGoroutineNum := runtime.NumGoroutine()

	server := socket.NewServer(network, "localhost:9527")

	if server == nil {
		b.Fatalf("launching server failed")
	}

	processStep := new(sync.WaitGroup)

	server.RegisterEvent(socket.SERVER_EVENT_ON_START, func() {
		processStep.Done()
	})

	const (
		wait_before_stop = 5 * time.Second
	)

	processStep.Add(1)

	go server.StartListen(wait_before_stop, func(ctx context.Context) {
		if ctx == nil {
			b.Error("accept nil ctx")
			return
		}

		session := NewSession(ctx)

		if session == nil {
			b.Error("accept session failed")
			return
		}

		defer session.Disconnect()

		if !session.WaitHandshake() {
			b.Error("server side session handshake failed")
			return
		}

		for session.IsConnected() {
			data, err := session.Recv()

			if err != nil {
				if err != io.ErrClosedPipe {
					b.Errorf("server side session recv data failed %v", err)
				}

				break
			}

			if !bytes.Equal(transData, data) {
				b.Errorf("server side session recv data failed %v msg broken", string(data))
			}
		}

		_, _, rcvPkgCnt, rcvPkgLen := session.Profile()
		log.Info("client session profile rcvPkgCnt %v rcvPkgLen %v", rcvPkgCnt, rcvPkgLen)
	})

	processStep.Wait() // wait server start

	session := NewSession(context.Background())

	if session == nil {
		b.Fatal("create client session failed")
	}

	if !session.TryConnect(network, "localhost:9527", needEncryption) {
		b.Fatal("client side session handshake failed")
	}

	b.ResetTimer()
	for i := 0; i < 10000; i++ {
		if !session.Send(transData) {
			b.Fatal("client side session send data faild")
		}
	}

	sndPkgCnt, sndPkgLen, _, _ := session.Profile()
	log.Info("server session profile sndPkgCnt %v sndPkgLen %v", sndPkgCnt, sndPkgLen)

	server.Close()
	session.Disconnect()

	log.Info("wait %v * 2 for server stop", 2*wait_before_stop)
	time.Sleep(2 * wait_before_stop)

	currentGoroutineNum := runtime.NumGoroutine()
	if systemGoroutineNum != currentGoroutineNum {
		b.Fatalf("goroutines before test %v,after test %v", systemGoroutineNum, currentGoroutineNum)
	}
}
