/*
	最受欢迎的词汇：
	功能：
	1.获取每5秒钟刷新的最受欢迎词汇
*/
package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) popularWordAck(result message.ErrCode, word string) {
	if result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Ulog(log.LOG_LEVEL_ERRO, log.PRINT_DEFINE|log.PRINT_STACKIN|2, "login failed %v", result)
	}

	u.send(utils.PackMsg(message.ChatServiceAPI_SERVICE_API_POPULAR_WORD_ACK, &message.API_POPULAR_WORD_ACK{
		Result: result,
		Word:   word,
	}))
}

func (u *user) popularWordReq(msg []byte) {
	if !u.sess.IsConnected() {
		return
	}

	var req message.API_POPULAR_WORD_REQ
	if !utils.UnpackMsg(msg, &req) {
		return
	}

	u.popularWordAck(message.ErrCode_ERR_CODE_SUCCESS, sharedMsgCache.heap.top())
}
