/*
	消息全局cache
	1.cache 历史N条记录(50)
	2.脏词过滤
	3.词频统计
	4.消息推送
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

const (
	MAX_CACHE_MSG_COUNT = 50
)

type msgCache struct {
	cachedMsg []string

	//用于循环cache最新的50条消息
	*ringBuffer
	trieRoot *trieNode
	heap     *maxHeap
}

var (
	sharedMsgCache = msgCache{
		heap:       newMaxHeap(),
		ringBuffer: newRingBuffer(MAX_CACHE_MSG_COUNT),
	}
)

func (mc *msgCache) updatePopularWord() {
	defer func() {
		if err := recover(); err != nil {
			log.Trac("err %v", err)
		}
	}()

	for {
		select {
		case <-time.Tick(5 * time.Second):
			mc.heap.reflush()

			word := mc.heap.top()
			log.Info("reflush most popular word is %v", word)
		}
	}
}

func (mc *msgCache) msgProcess(msg string) string {
	if mc.trieRoot == nil {
		return msg
	}

	clean := mc.trieRoot.replace(msg, '*')

	for _, word := range splitWords(clean) {
		mc.heap.record(word)
	}

	return clean
}

func (mc *msgCache) sendMsgByID(from string, to int64, msg string) (ok bool) {
	target := sharedUserMgr.findUserByID(to)

	if target == nil {
		return
	}

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_MSG_NOTIFY, &message.API_NEW_MESSAGE_NOTIFY{
		From:     from,
		Msg:      mc.msgProcess(msg),
		SendUnix: time.Now().Unix(),
		Type:     message.MESSAGE_TYPE_MESSAGE_TYPE_PRIVATE,
	})

	if len(data) == 0 {
		return
	}

	return target.send(data)
}

func (mc *msgCache) sendMsgByName(from, to string, msg string) (ok bool) {
	target := sharedUserMgr.findUserByName(to)

	if target == nil {
		return
	}

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_MSG_NOTIFY, &message.API_NEW_MESSAGE_NOTIFY{
		From:     from,
		Msg:      mc.msgProcess(msg),
		SendUnix: time.Now().Unix(),
		Type:     message.MESSAGE_TYPE_MESSAGE_TYPE_PRIVATE,
	})

	if len(data) == 0 {
		return
	}

	return target.send(data)
}

func (mc *msgCache) broadcast(fromName string, msg string) {
	ntf := &message.API_NEW_MESSAGE_NOTIFY{
		From:     fromName,
		Msg:      mc.msgProcess(msg),
		SendUnix: time.Now().Unix(),
		Type:     message.MESSAGE_TYPE_MESSAGE_TYPE_BROADCAST,
	}

	mc.append(ntf)

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_MSG_NOTIFY, ntf)

	if len(data) == 0 {
		return
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Trac("err %v", err)
			}
		}()

		for _, user := range sharedUserMgr.getAllUser() {
			if user == nil {
				continue
			}

			user.send(data)
		}
	}()
}
