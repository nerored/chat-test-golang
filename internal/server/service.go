/*
	服务实现：
	功能：
	1.accpet 新的回话，并尝试进行握手
	2.握手成功后开始协议派发
*/
package main

import (
	"context"
	"errors"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/session"
	"github.com/nerored/chat-test-golang/net/socket"
)

func serviceRun(listenAddr, profanityfilePath string) error {
	sharedMsgCache.trieRoot = loadProfanityWordsFromFile(profanityfilePath)

	server := socket.NewServer("tcp", listenAddr)

	if server == nil {
		return errors.New("create server failed")
	}

	go sharedMsgCache.updatePopularWord()

	server.StartListen(5*time.Second, func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Trac("err %v", err)
			}

			log.Info("disconnected")
		}()

		sess := session.NewSession(ctx)

		if sess == nil {
			log.Erro("accept failed")
			return
		}

		defer sess.Disconnect()

		if !sess.WaitHandshake() {
			log.Erro("handshake failed")
			return
		}

		user := sharedUserMgr.createUser(sess)

		if user == nil {
			log.Erro("createUser failed")
			return
		}

		user.recvmsg()
	})

	return nil
}
