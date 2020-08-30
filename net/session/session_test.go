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

func TestTCPSessionConnect1(t *testing.T) {
	SessionConnect("tcp", true, t)
}

func TestTCPSessionConnect2(t *testing.T) {
	SessionConnect("tcp", false, t)
}

func SessionConnect(network string, needEncryption bool, t *testing.T) {
	log.InitLog("")

	systemGoroutineNum := runtime.NumGoroutine()

	server := socket.NewServer(network, "localhost:9527")

	if server == nil {
		t.Fatalf("launching server failed")
	}

	processStep := new(sync.WaitGroup)

	server.RegisterEvent(socket.SERVER_EVENT_ON_START, func() {
		processStep.Done()
	})

	const (
		wait_before_stop = 5 * time.Second
	)

	msg := []byte("hello nestnet by session")

	processStep.Add(1)

	go server.StartListen(wait_before_stop, func(ctx context.Context) {
		if ctx == nil {
			t.Error("accept nil ctx")
			return
		}

		session := NewSession(ctx)

		if session == nil {
			t.Error("accept session failed")
			return
		}

		defer session.Disconnect()

		if !session.WaitHandshake() {
			t.Error("server side session handshake failed")
			return
		}

		for session.IsConnected() {
			data, err := session.Recv()

			if err != nil {
				if err != io.ErrClosedPipe {
					t.Errorf("server side session recv data failed %v", err)
				}

				break
			}

			if !bytes.Equal(msg, data) {
				t.Errorf("server side session recv data failed %v msg broken", string(data))
				break
			}

			if !session.Send(data) {
				t.Error("server side session send data failed")
			}
		}
	})

	processStep.Wait() // wait server start

	session := NewSession(context.Background())

	if session == nil {
		t.Fatal("create client session failed")
	}

	if !session.TryConnect(network, "localhost:9527", needEncryption) {
		t.Fatal("client side session handshake failed")
	}

	if !session.Send(msg) {
		t.Fatalf("client side session send data faild")
	}

	session.SetRTimeout(3 * time.Second)

	for session.IsConnected() {
		data, err := session.Recv()

		if err != nil {
			if err != io.ErrClosedPipe {
				t.Fatalf("client side session recv data failed %v", err)
			}
			break
		}

		if !bytes.Equal(msg, data) {
			t.Fatalf("client side session recv data failed %v msg broken", string(data))
		}
	}

	server.Close()
	session.Disconnect()

	t.Logf("wait %v * 2 for server stop", 2*wait_before_stop)
	time.Sleep(2 * wait_before_stop)

	currentGoroutineNum := runtime.NumGoroutine()
	if systemGoroutineNum != currentGoroutineNum {
		t.Fatalf("goroutines before test %v,after test %v", systemGoroutineNum, currentGoroutineNum)
	}
}
