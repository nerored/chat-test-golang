/*
	客户端协议分发器
*/
package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func dispatchmsg(u *user, packet []byte) {
	if u == nil || len(packet) == 0 {
		return
	}

	api, msg := utils.TranslateMsg(packet)

	switch api {
	case message.ChatServiceAPI_SERVICE_API_LOGIN_ACK:
		u.loginAck(msg)
	case message.ChatServiceAPI_SERVICE_API_SEND_ACK:
		u.sendAck(msg)
	case message.ChatServiceAPI_SERVICE_API_STATS_ACK:
		u.statsAck(msg)
	case message.ChatServiceAPI_SERVICE_API_MSG_NOTIFY:
		u.newMsgNotify(msg)
	case message.ChatServiceAPI_SERVICE_API_POPULAR_WORD_ACK:
		u.popularWordAck(msg)
	default:
		log.Warn("recv unknown api %v from user id %v name %v", api, u.id, u.name)
	}
}
