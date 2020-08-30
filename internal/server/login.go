/*
	登陆协议实现：
	功能：
	1.登陆时名字注册
	2.注册成功则返回最近50条世界消息
*/
package main

import (
	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) loginAck(result message.ErrCode, cachedMsg []string) {
	if result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Ulog(log.LOG_LEVEL_ERRO, log.PRINT_DEFINE|log.PRINT_STACKIN|2, "login failed %v", result)
	}

	ack := &message.API_LOGIN_ACK{
		Result: result,
		UserID: u.id,
		Name:   u.name,
	}

	if result == message.ErrCode_ERR_CODE_SUCCESS {
		ack.CachedMsg = sharedMsgCache.readAll()
	}

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_LOGIN_ACK, ack)

	u.send(data)
}

func (u *user) loginReq(msg []byte) {
	if len(msg) == 0 || !u.sess.IsConnected() {
		return
	}

	var req message.API_LOGIN_REQ
	if !utils.UnpackMsg(msg, &req) {
		return
	}

	if req.Name == "" || req.Name == u.name {
		u.loginAck(message.ErrCode_ERR_CODE_FAILED, nil)
		return
	}

	if !sharedUserMgr.registerName(u.id, req.Name) {
		u.loginAck(message.ErrCode_ERR_CODE_DUPLICATE_NAME, nil)
		return
	}

	u.name = req.Name
	log.Info("user id %v name %v login success", u.id, u.name)
	u.loginAck(message.ErrCode_ERR_CODE_SUCCESS, nil)
}
