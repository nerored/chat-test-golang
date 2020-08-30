package main

import (
	"errors"

	"github.com/nerored/chat-test-golang/cli/cmd"
)

type WordCommand struct {
	cmd.BaseCommand
}

func NewWordCommand() cmd.Command {
	return &WordCommand{}
}

func (ec *WordCommand) Init(appCtrl cmd.Control) (err error) {
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

func (ec *WordCommand) Name() string {
	return "popular"
}

func (ec *WordCommand) Usage() string {
	return "查询5秒内最受欢迎的词汇"
}

func (ec *WordCommand) Exec(c *cmd.ArgsContext) (err error) {
	sharedUser.popularWordReq()
	return
}
