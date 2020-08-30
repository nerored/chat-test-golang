/*
	使用ringbuffer来循环保存N条最新的消息
*/
package main

import (
	"container/ring"
	"sync"

	"github.com/nerored/chat-test-golang/message"
)

type ringBuffer struct {
	lastNode *ring.Ring
	sync.RWMutex
}

func newRingBuffer(n int) *ringBuffer {
	return &ringBuffer{
		lastNode: ring.New(n),
	}
}

func (mb *ringBuffer) append(msg *message.API_NEW_MESSAGE_NOTIFY) {
	if msg == nil || mb.lastNode == nil {
		return
	}

	mb.Lock()
	defer mb.Unlock()

	mb.lastNode.Value = msg
	mb.lastNode = mb.lastNode.Next()
}

func (mb *ringBuffer) readAll() (messages []*message.API_NEW_MESSAGE_NOTIFY) {
	if mb.lastNode == nil {
		return
	}

	mb.RLock()
	defer mb.RUnlock()

	messages = make([]*message.API_NEW_MESSAGE_NOTIFY, 0, mb.lastNode.Len())

	mb.lastNode.Do(func(data interface{}) {
		msg, ok := data.(*message.API_NEW_MESSAGE_NOTIFY)

		if !ok || msg == nil {
			return
		}

		messages = append(messages, msg)
	})

	return
}
