/*
	用户对象
	1.持有与service的回话状态
	2.网络通信
*/
package main

import (
	"context"
	"os"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/message"
	"github.com/nerored/chat-test-golang/net/session"
	"github.com/nerored/chat-test-golang/utils"
	"google.golang.org/protobuf/proto"
)

type user struct {
	id   int64
	name string
	join time.Time

	sess *session.Session
}

var (
	sharedUser user
)

func (u *user) init() (ok bool) {
	if u.sess != nil && u.sess.IsConnected() {
		return
	}

	sess := session.NewSession(context.Background())

	if sess == nil {
		log.Erro("create session failed")
		return
	}

	if !sess.TryConnect("tcp", ":9527", true) {
		log.Erro("handshake failed")
		return
	}

	sess.SetRTimeout(15 * time.Minute)
	u.sess = sess
	return true
}

func (u *user) isVerified() bool {
	return u.id > 0 && u.name != ""
}

func (u *user) send(api message.ChatServiceAPI, msg proto.Message) (ok bool) {
	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	data := utils.PackMsg(api, msg)

	if len(data) == 0 {
		return
	}

	return u.sess.Send(data)
}

func (u *user) recvmsg() {
	defer func() {
		if err := recover(); err != nil {
			log.Trac("err %v", err)
		}

		log.Info("disconnected")
		os.Exit(1)
	}()

	if u.sess == nil || !u.sess.IsConnected() {
		return
	}

	for u.sess.IsConnected() {
		msg, err := u.sess.Recv()

		if err != nil {
			log.Erro("user id %v join %v recv err %v", u.id, u.join, err)
			return
		}

		dispatchmsg(u, msg)
	}
}
