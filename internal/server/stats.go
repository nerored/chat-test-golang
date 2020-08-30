/*
	状态查询协议实现：
	功能：
	1.根据用户名字查询登陆时长
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) statsAck(result message.ErrCode, joinTime int64) {
	if result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Ulog(log.LOG_LEVEL_ERRO, log.PRINT_DEFINE|log.PRINT_STACKIN|2, "login failed %v", result)
	}

	ack := &message.API_STATS_ACK{
		Result:   result,
		JoinTime: joinTime,
	}

	data := utils.PackMsg(message.ChatServiceAPI_SERVICE_API_STATS_ACK, ack)

	u.send(data)
}

func (u *user) statsReq(msg []byte) {
	if len(msg) == 0 || !u.sess.IsConnected() {
		return
	}

	var req message.API_STATS_REQ
	if !utils.UnpackMsg(msg, &req) {
		return
	}

	user := sharedUserMgr.findUserByName(req.Name)

	if user == nil {
		u.statsAck(message.ErrCode_ERR_CODE_USER_IS_NOT_EXIST, 0)
		return
	}

	u.statsAck(message.ErrCode_ERR_CODE_SUCCESS, int64(time.Now().Sub(user.join)))
}
