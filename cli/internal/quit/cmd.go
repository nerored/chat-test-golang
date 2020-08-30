/*
	退出指令，在cli中使用quit，可以终止程序运行
	此指令存在的意义是使用的第三方包catch了 C-D
*/
package quit

import (
	"errors"

	"github.com/nerored/chat-test-golang/cli/cmd"
)

type ExitCommand struct {
	cmd.BaseCommand
}

func NewCommand() cmd.Command {
	return &ExitCommand{}
}

func (ec *ExitCommand) Init(appCtrl cmd.Control) (err error) {
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

func (ec *ExitCommand) Name() string {
	return "quit"
}

func (ec *ExitCommand) Usage() string {
	return "退出"
}

func (ec *ExitCommand) Exec(c *cmd.ArgsContext) (err error) {
	ec.AppControl.Exit()
	return
}
