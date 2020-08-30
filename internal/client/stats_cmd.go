/*
	stats 协议的cli衔接指令
*/
package main

import (
	"errors"

	"github.com/nerored/chat-test-golang/cli/cmd"
	"github.com/nerored/chat-test-golang/cli/log"
)

type StatsCommand struct {
	UserName string `short:"n" usage:"需要查询的用户Name"`
	cmd.BaseCommand
}

func NewStatsCommand() cmd.Command {
	return &StatsCommand{}
}

func (ec *StatsCommand) Init(appCtrl cmd.Control) (err error) {
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

func (ec *StatsCommand) Name() string {
	return "stats"
}

func (ec *StatsCommand) Usage() string {
	return `查询用户在线时长`
}

func (ec *StatsCommand) Exec(c *cmd.ArgsContext) (err error) {
	cmd.SetArgs(ec, c)

	if ec.UserName == "" {
		log.Info("empty name")
		return
	}

	sharedUser.statsReq(ec.UserName)
	return
}
