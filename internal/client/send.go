/*
	1.发送消息协议实现：
	a.可对用户（id,name索引均可，id匹配优先级最高）发送消息
	b.进行消息广播

	2.service 新消息通知处理
*/

package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) newMsgNotify(data []byte) {
	if u.sess == nil || !u.sess.IsConnected() || len(data) == 0 {
		return
	}

	var ntf message.API_NEW_MESSAGE_NOTIFY
	if !utils.UnpackMsg(data, &ntf) {
		return
	}

	printMsg(&ntf)
}

func (u *user) sendAck(data []byte) {
	if u.sess == nil || !u.sess.IsConnected() || len(data) == 0 {
		return
	}

	var ack message.API_SEND_ACK
	if !utils.UnpackMsg(data, &ack) {
		return
	}

	if ack.Result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Erro("send failed")
		return
	}
}

func (u *user) sendReq(toID int64, toName string, msg string) {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	if !u.isVerified() {
		log.Warn("please login first")
		return
	}

	u.send(message.ChatServiceAPI_SERVICE_API_SEND_REQ, &message.API_SEND_REQ{
		Name:    toName,
		UserID:  toID,
		Message: msg,
	})
}
