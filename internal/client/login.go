/*
	1.登陆协议，及其回包处理
	a.注册名字
	b.输出历史消息
*/
package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) loginAck(data []byte) {
	if u.sess == nil || !u.sess.IsConnected() || len(data) == 0 {
		return
	}

	var ack message.API_LOGIN_ACK
	if !utils.UnpackMsg(data, &ack) {
		return
	}

	if ack.Result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Erro("login failed")
		u.sess.Disconnect()
		return
	}

	u.id = ack.UserID
	u.name = ack.Name

	log.Info("%v user id %v name %v", log.NewCombo("login success", log.FGC_GREEN), u.id, u.name)

	for _, msg := range ack.CachedMsg {
		printMsg(msg)
	}
}

func (u *user) loginReq(name string) {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	if u.isVerified() {
		log.Warn("user is login at id %v name %v", u.id, u.name)
		return
	}

	u.send(message.ChatServiceAPI_SERVICE_API_LOGIN_REQ, &message.API_LOGIN_REQ{
		Name: name,
	})
}
