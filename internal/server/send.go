/*
	发送消息协议实现：
	功能：
	1.当Name不为空且userID为0时，向name的用户发送消息
	2.userID不为0时，向id指向的用户发送消息
	3.当name与userID均未指定，则进行消息广播
*/
package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) sendAck(result message.ErrCode) {
	if result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Ulog(log.LOG_LEVEL_ERRO, log.PRINT_DEFINE|log.PRINT_STACKIN|2, "login failed %v", result)
	}

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_SEND_ACK, &message.API_LOGIN_ACK{
		Result: result,
	})

	u.send(data)
}

func (u *user) sendReq(msg []byte) {
	if len(msg) == 0 || !u.sess.IsConnected() {
		return
	}

	var req message.API_SEND_REQ
	if !utils.UnpackMsg(msg, &req) {
		return
	}

	isSuccess := false
	switch {
	case req.UserID > 0:
		isSuccess = sharedMsgCache.sendMsgByID(u.name, req.UserID, req.Message)
	case req.Name != "":
		isSuccess = sharedMsgCache.sendMsgByName(u.name, req.Name, req.Message)
	default:
		sharedMsgCache.broadcast(u.name, req.Message)
		isSuccess = true
	}

	if !isSuccess {
		u.sendAck(message.ErrCode_ERR_CODE_FAILED)
		return
	}

	u.sendAck(message.ErrCode_ERR_CODE_SUCCESS)
}
