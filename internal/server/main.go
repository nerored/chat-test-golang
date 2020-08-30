/*
	service 入口
*/
package main

import (
	"os"

	"github.com/nerored/chat-test-golang/log"
	"github.com/urfave/cli/v2"
)

func main() {
	log.InitLog("")
	log.SetDebugPrint(false)

	serviceApp := cli.NewApp()
	serviceApp.Name = "chat-test-service"

	serviceApp.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "listen,l",
			Value: ":9527",
			Usage: "listen address",
		},
		&cli.StringFlag{
			Name:  "profanitywords,w",
			Value: "./list.txt",
			Usage: "load profanitywords",
		},
	}

	serviceApp.Action = func(c *cli.Context) error {
		return serviceRun(c.String("listen"), c.String("profanitywords"))
	}

	serviceApp.Run(os.Args)
}
