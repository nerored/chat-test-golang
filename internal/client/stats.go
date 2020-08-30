/*
	用户状态查询协议
	1.输出查询到的指定用户的在线时长
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/utils"
)

func (u *user) statsAck(data []byte) {
	if u.sess == nil || !u.sess.IsConnected() || len(data) == 0 {
		return
	}

	var ack message.API_STATS_ACK
	if !utils.UnpackMsg(data, &ack) {
		return
	}

	if ack.Result != message.ErrCode_ERR_CODE_SUCCESS {
		log.Erro("stats query failed %v", ack.Result)
		return
	}

	loginTime := time.Duration(ack.JoinTime)

	d := loginTime / (24 * time.Hour)
	loginTime %= 24 * time.Hour
	h := loginTime / time.Hour
	loginTime %= time.Hour
	m := loginTime / time.Minute
	loginTime %= time.Minute
	s := loginTime / time.Second

	log.Info("%.2dd %.2dh %.2dm %.2ds", d, h, m, s)
}

func (u *user) statsReq(name string) {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	if !u.isVerified() {
		log.Warn("please login first")
		return
	}

	u.send(message.ChatServiceAPI_SERVICE_API_STATS_REQ, &message.API_STATS_REQ{
		Name: name,
	})
}
