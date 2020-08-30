/*
	chat client:
	1.启动时自动链接service，并要求输入名字，以完成登陆过程(service 地址写死)
	2.登陆成功后，启动cli app 接管接下来的操作
*/
package main

import (
	"fmt"
	"time"

	"github.com/nerored/chat-test-golang/cli"
	"github.com/nerored/chat-test-golang/log"
)

func main() {
	log.InitLog("")

	app := cli.NewApp("/")

	if app == nil {
		log.Erro("client init faield")
		return
	}

	log.SetDebugPrint(false)

	if !sharedUser.init() {
		log.Erro("client conn faild")
		return
	}

	go sharedUser.recvmsg()

	log.Prefix("%v",
		log.NewCombo("Enter your name:", log.FGC_LIGHTBLUE))

	var name string

	_, err := fmt.Scanf("%s", &name)

	if err != nil {
		log.Erro("%v", err)
		return
	}

	sharedUser.loginReq(name)

	//等待登陆成功后进行操作
	for !sharedUser.isVerified() {
		time.Sleep(100 * time.Millisecond)
	}

	app.RegisterCmd(NewSendCommand())
	app.RegisterCmd(NewWordCommand())
	app.RegisterCmd(NewStatsCommand())
	app.Run()
}
