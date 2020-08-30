/*
	服务端的消息分发器
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/cli/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func dispatchmsg(u *user, packet []byte) {
	if u == nil || len(packet) == 0 {
		return
	}

	api, msg := utils.TranslateMsg(packet)

	defer func(begin time.Time) {
		log.Info("process msg %v len %v to user id %v name %v cost %v",
			api, len(msg), u.id, u.name, time.Now().Sub(begin))
	}(time.Now())

	switch api {
	case message.ChatServiceAPI_SERVICE_API_LOGIN_REQ:
		u.loginReq(msg)
	case message.ChatServiceAPI_SERVICE_API_SEND_REQ:
		u.sendReq(msg)
	case message.ChatServiceAPI_SERVICE_API_STATS_REQ:
		u.statsReq(msg)
	case message.ChatServiceAPI_SERVICE_API_POPULAR_WORD_REQ:
		u.popularWordReq(msg)
	default:
		log.Warn("recv unknown api %v from user id %v name %v", api, u.id, u.name)

		if u.sess != nil {
			u.sess.Disconnect()
		}
	}
}
