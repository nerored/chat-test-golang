/*
	聊天消息的友好打印
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
)

func printMsg(ntf *message.API_NEW_MESSAGE_NOTIFY) {
	if ntf == nil {
		return
	}

	sendTime := time.Unix(ntf.SendUnix, 0).Format("2006-01-02 15:04:05")

	var format int

	switch ntf.Type {
	case message.MESSAGE_TYPE_MESSAGE_TYPE_UNKNOWN:
		fallthrough
	case message.MESSAGE_TYPE_MESSAGE_TYPE_PRIVATE:
		format = log.FGC_LIGHTMAGENTA
	case message.MESSAGE_TYPE_MESSAGE_TYPE_BROADCAST:
		format = log.FGC_BLUE
	}

	log.Info("%v [%v]:%v",
		log.NewCombo(sendTime, format),
		log.NewCombo(ntf.From, format),
		log.NewCombo(ntf.Msg, format))
}
