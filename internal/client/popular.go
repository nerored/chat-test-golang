package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) popularWordAck(data []byte) {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	var ack message.API_POPULAR_WORD_ACK
	if !utils.UnpackMsg(data, &ack) {
		return
	}

	if ack.Result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Erro("get popularWord failed")
		return
	}

	log.Info("most popularWord is [%v]", log.NewCombo(ack.Word, log.FGC_GREEN))
}

func (u *user) popularWordReq() {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	if !u.isVerified() {
		log.Warn("please login first")
		return
	}

	u.send(message.ChatServiceAPI_SERVICE_API_POPULAR_WORD_REQ, &message.API_POPULAR_WORD_REQ{})
}
