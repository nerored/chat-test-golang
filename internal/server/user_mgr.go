/*
	负责登陆用户的管理
	功能:
	1.接受sess建立新用户
	2.注册用户名字，并去重
	3.进行用户的查找
*/
package main

import (
	"sync"
	"time"

	"github.com/nerored/chat-test-golang/log"
	"github.com/nerored/chat-test-golang/net/session"
)

type idGen struct {
	freeIDIdx int64
	sync.Mutex
}

func (i *idGen) freeID() int64 {
	i.Lock()
	defer i.Unlock()

	i.freeIDIdx++
	return i.freeIDIdx
}

type userMgr struct {
	idGen
	id2User map[int64]*user
	name2ID map[string]int64
	sync.RWMutex
}

var (
	sharedUserMgr = userMgr{
		id2User: make(map[int64]*user),
		name2ID: make(map[string]int64),
	}
)

func (um *userMgr) createUser(sess *session.Session) *user {
	if sess == nil || !sess.IsConnected() {
		return nil
	}

	newUser := &user{
		id:   um.freeID(),
		sess: sess,
		name: "Anonymous",
		join: time.Now(),
	}

	return um.add(newUser)
}

func (um *userMgr) add(user *user) *user {
	if user == nil {
		return nil
	}

	um.Lock()
	defer um.Unlock()

	um.id2User[user.id] = user

	log.Info("user id %v name %v joined", user.id, user.name)
	return user
}

func (um *userMgr) del(user *user) *user {
	if user == nil {
		return nil
	}

	um.Lock()
	defer um.Unlock()

	delete(um.id2User, user.id)
	delete(um.name2ID, user.name)

	log.Info("user id %v name %v exited", user.id, user.name)
	return user
}

func (um *userMgr) findUserByID(id int64) *user {
	um.RLock()
	defer um.RUnlock()

	return um.id2User[id]
}

func (um *userMgr) findUserByName(name string) *user {
	um.RLock()
	defer um.RUnlock()

	id, ok := um.name2ID[name]

	if !ok {
		return nil
	}

	return um.id2User[id]
}

func (um *userMgr) getAllUser() []*user {
	um.RLock()
	defer um.RUnlock()

	users := make([]*user, 0, len(um.id2User))

	for _, user := range um.id2User {
		if user == nil {
			continue
		}

		users = append(users, user)
	}

	return users
}

func (um *userMgr) registerName(id int64, name string) (ok bool) {
	um.Lock()
	defer um.Unlock()

	if owner, ok := um.name2ID[name]; ok {
		return owner == id
	}

	um.name2ID[name] = id
	return true
}
