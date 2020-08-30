/*
	用户对象：
	1.维持会话
	2.与其他用户进行交互

	note:由于数据基本只读，所以无需加锁.
*/
package main

import (
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/session"
)

type user struct {
	id   int64
	name string
	join time.Time
	sess *session.Session
}

func (u *user) recvmsg() {
	if u.sess == nil {
		return
	}

	defer sharedUserMgr.del(u)

	for u.sess.IsConnected() {
		msg, err := u.sess.Recv()

		if err != nil {
			log.Erro("user id %v join %v recv err %v", u.id, u.join, err)
			return
		}

		dispatchmsg(u, msg)
	}
}

func (u *user) send(data []byte) (ok bool) {
	if len(data) == 0 || u.sess == nil || !u.sess.IsConnected() {
		return
	}

	if !u.sess.Send(data) {
		log.Ulog(log.LOG_LEVEL_ERRO, log.PRINT_DEFINE|log.PRINT_STACKIN|3, "send failed")
		return
	}

	return true
}
