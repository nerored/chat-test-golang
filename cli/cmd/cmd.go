/*
	command 运行支持
	扩展cmd需要组合BaseCommand 并实现自己的运行接口
	即可和cli进行联动，来扩展指令
*/
package cmd

import (
	"errors"
	"reflect"
	"strings"

	"github.com/nerored/chat-test-golang/cli/log"
	args "github.com/urfave/cli/v2"
)

type Command interface {
	Init(appCtrl Control) error
	Exec(c *ArgsContext) error
	Name() string
	Usage() string
	Action() *ArgsAction
}

type BaseCommand struct {
	AppControl Control
	ArgsAction *ArgsAction
}

func (bc *BaseCommand) Init(appCtrl Control) (err error) {
	if appCtrl == nil {
		return errors.New("no app control interface")
	}

	bc.AppControl = appCtrl
	return nil
}

func (bc *BaseCommand) Name() string {
	return ""
}

func (bc *BaseCommand) Usage() string {
	return ""
}

func (bc *BaseCommand) Action() *ArgsAction {
	return bc.ArgsAction
}

func (bc *BaseCommand) Exec(c *ArgsContext) (err error) {
	return errors.New("no action,please implement action method")
}

func BuildSupport(command Command) (argsAction *ArgsAction) {
	if command == nil {
		log.DebuS("build cmd is a nil object")
		return
	}

	refV := reflect.ValueOf(command)

	if refV.Kind() != reflect.Ptr {
		log.DebuS("build cmd is not a ptr")
		return
	}

	elem := refV.Elem()

	if elem.Kind() != reflect.Struct {
		log.DebuS("build cmd is not a strcut ptr")
		return
	}

	refT := elem.Type()

	argsAction = buildSupportForArgs(command.Name(), refT)

	if argsAction == nil {
		return
	}

	err := sharedAutoComple.register(command.Name(), command.Usage(), refT)

	if err != nil {
		log.DebuS("build cmd register auto comple failed")
		return
	}

	argsAction.Action = func(c *args.Context) (err error) {
		return command.Exec(&ArgsContext{
			Context: c,
		})
	}

	return
}

func SetArgs(o interface{}, c *ArgsContext) {
	if o == nil || c == nil {
		return
	}

	refV := reflect.ValueOf(o)

	if refV.Kind() != reflect.Ptr {
		return
	}

	elem := refV.Elem()

	if elem.Kind() != reflect.Struct {
		return
	}

	refT := elem.Type()

	for i := 0; i < refT.NumField(); i++ {
		field := refT.Field(i)
		elemF := elem.Field(i)

		if field.Anonymous {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			elemF.SetString(c.String(strings.ToLower(field.Name)))
		case reflect.Bool:
			elemF.SetBool(c.Bool(strings.ToLower(field.Name)))
		}
	}
}
