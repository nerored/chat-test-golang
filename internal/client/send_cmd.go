/*
	send 协议的cli衔接指令
*/
package main

import (
	"errors"

	"github.com/nerored/chat-test-golang/cli/cmd"
)

type SendCommand struct {
	ToID    int64  `short:"t" usage:"接收方UserID"`
	ToName  string `short:"n" usage:"接收方Name"`
	Message string `short:"m" usage:"要发送的消息"`
	cmd.BaseCommand
}

func NewSendCommand() cmd.Command {
	return &SendCommand{}
}

func (ec *SendCommand) Init(appCtrl cmd.Control) (err error) {
	err = ec.BaseCommand.Init(appCtrl)

	if err != nil {
		return
	}

	argsAction := cmd.BuildSupport(ec)

	if argsAction == nil {
		return errors.New("can't build argsAction")
	}

	ec.ArgsAction = argsAction
	return nil
}

func (ec *SendCommand) Name() string {
	return "send"
}

func (ec *SendCommand) Usage() string {
	return `发送消息`
}

func (ec *SendCommand) Exec(c *cmd.ArgsContext) (err error) {
	cmd.SetArgs(ec, c)
	sharedUser.sendReq(ec.ToID, ec.ToName, ec.Message)
	return
}
