package socket

import (
	"bytes"
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/nerored/chat-test-golang/log"
)

func TestTCPClientConnect(t *testing.T) {
	ClientConnect("tcp", t)
}

func ClientConnect(network string, t *testing.T) {
	log.InitLog("")

	systemGoroutineNum := runtime.NumGoroutine()

	server := NewServer(network, "localhost:9527")

	if server == nil {
		t.Fatalf("launching server failed")
	}

	processStep := new(sync.WaitGroup)

	server.RegisterEvent(SERVER_EVENT_ON_START, func() {
		processStep.Done()
	})

	const (
		wait_before_stop = 5 * time.Second
	)

	msg := []byte("hello nestnet")

	processStep.Add(1)

	go server.StartListen(wait_before_stop, func(ctx context.Context) {
		if ctx == nil {
			t.Error("accept nil ctx")
			return
		}

		client := ctx.Value(SERVER_CTX_PARAM_CONN).(*Client)

		if client == nil || !client.IsAlive() {
			t.Error("ctx not include alive client obj")
			return
		}

		defer client.Close()

		recvChan := make(chan []byte)

		go client.Recv(recvChan)

	recvTag:
		for client.IsAlive() {
			select {
			case <-ctx.Done():
				t.Logf("server side client stoped by %v", ctx.Err())
				break recvTag
			case data, ok := <-recvChan:
				if !ok {
					t.Log("server side client chan closed")
					break recvTag
				}

				if !bytes.Equal(data, msg) {
					t.Errorf("server side client recv data failed %v", data)
					break recvTag
				}

				client.Send(data)
			}
		}
	})

	processStep.Wait() // wait server start

	client := DialClient(network, "localhost:9527")

	if client == nil {
		t.Fatalf("launching client failed")
	}

	recvChan := make(chan []byte)

	go client.Recv(recvChan)

	if !client.Send(msg) {
		t.Fatalf("send data to server failed")
	}

end:
	for {
		select {
		case <-time.After(3 * time.Second):
			t.Fatalf("wait server ack timeout")
		case data, ok := <-recvChan:
			if !ok {
				t.Fatal("client side client chan closed")
			}

			if !bytes.Equal(data, msg) {
				t.Fatalf("client side client recv data failed %v", data)
			}

			break end
		}
	}

	server.Close()
	client.Close()

	t.Logf("wait %v * 2 for server stop", 2*wait_before_stop)
	time.Sleep(2 * wait_before_stop)

	currentGoroutineNum := runtime.NumGoroutine()
	if systemGoroutineNum != currentGoroutineNum {
		t.Fatalf("goroutines before test %v,after test %v", systemGoroutineNum, currentGoroutineNum)
	}
}
