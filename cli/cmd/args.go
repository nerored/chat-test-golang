/*
	参数解析:
	通过反射将cmd结构中的每个field的的类型，structtag解析
	生成cli所需要的自动补全的信息
*/
package cmd

import (
	"reflect"
	"strings"

	args "github.com/urfave/cli/v2"
)

func buildSupportForArgs(name string, refT reflect.Type) (command *ArgsAction) {
	if len(name) <= 0 || refT.Kind() == reflect.Invalid || refT.Kind() != reflect.Struct {
		return
	}

	command = NewArgsAction()
	command.Name = name

	for i := 0; i < refT.NumField(); i++ {
		field := refT.Field(i)

		fieldName := strings.ToLower(field.Name)

		var aliases []string

		if shortName := field.Tag.Get("short"); len(shortName) > 0 {
			aliases = append(aliases, shortName)
		}

		switch field.Type.Kind() {
		case reflect.Bool:
			command.Flags = append(command.Flags, &args.BoolFlag{
				Name:    fieldName,
				Aliases: aliases,
				Usage:   field.Tag.Get("usage"),
			})
		case reflect.String:
			command.Flags = append(command.Flags, &args.StringFlag{
				Name:    fieldName,
				Aliases: aliases,
				Usage:   field.Tag.Get("usage"),
			})
		}
	}

	return
}
