/*
	解藕封装，不对外暴露cli包
*/
package cmd

import args "github.com/urfave/cli/v2"

type ArgsAction struct {
	*args.Command
}

func NewArgsAction() *ArgsAction {
	return &ArgsAction{
		Command: new(args.Command),
	}
}

type ArgsContext struct {
	*args.Context
}
