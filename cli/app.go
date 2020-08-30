/*
	主控app:
	功能：
	1.注册扩展指令，并将消息传递给对应的指令.
	2.输入循环保持，提供补全和历史记录
*/
package cli

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/nerored/chat-test-golang/cli/cmd"
	"github.com/nerored/chat-test-golang/cli/internal/quit"
	"github.com/nerored/chat-test-golang/cli/log"

	args "github.com/urfave/cli/v2"
)

const (
	APP_RUN_STATE_RUNNING = iota
	APP_RUN_STATE_EXIT
)

const (
	APP_HISTORY_FILENAME = ".clihistory"
)

type App struct {
	runOnce    bool
	argsApp    *args.App
	history    []string
	commandMap map[string]cmd.Command

	defaultPrefix string
	runningState  int
}

func NewApp(prefix string) *App {
	log.Init()

	app := &App{
		defaultPrefix: prefix,
		commandMap:    make(map[string]cmd.Command),
	}

	app.RegisterCmd(quit.NewCommand())

	return app
}

//-------------------  control interface

func (app *App) Exit() {
	app.runningState = APP_RUN_STATE_EXIT
}

func (app *App) ExitHard() {
	log.Info("bye :)")
	os.Exit(0)
}

func (app *App) ClearHistory() {
	app.history = nil
	app.saveHistoryFile()

	log.Info("[app] input history is clear")
}

func (app *App) isRunning() bool {
	return app.runningState == APP_RUN_STATE_RUNNING
}

//-------------------  control interface

func (app *App) Run() {
	app.init()

	if len(os.Args) > 1 {
		err := app.argsApp.Run(os.Args)

		if err != nil {
			return
		}
	}

	if app.runOnce {
		return
	}

	pt := prompt.New(func(_ string) {
	}, cmd.Completer)

	if pt == nil {
		return
	}

	options := []prompt.Option{
		prompt.OptionPrefix(app.defaultPrefix),
		prompt.OptionCompletionWordSeparator(" ="),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Cyan),
		prompt.OptionSelectedDescriptionBGColor(prompt.Cyan),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.White),
		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionScrollbarBGColor(prompt.DarkGray),
		prompt.OptionScrollbarThumbColor(prompt.LightGray),
		prompt.OptionSwitchKeyBindMode(prompt.CommonKeyBind),
	}

	for _, opt := range options {
		if err := opt(pt); err != nil {
			return
		}
	}

	app.loadHistoryFile()
	app.activeHistory(pt)

	for ; app.isRunning(); app.activeHistory(pt) {
		input := pt.Input()

		if input == "" {
			continue
		}

		var arguments = []string{"cli"}

		const (
			state_normal = 0
			state_block  = 1
		)

		currentState := state_normal
		var buffer []rune
		for _, r := range input {
			switch currentState {
			case state_normal:
				switch r {
				case ' ':
					arguments = append(arguments, string(buffer))
					buffer = buffer[:0]
				case '"':
					fallthrough
				case '\'':
					currentState = state_block
				default:
					buffer = append(buffer, r)
				}
			case state_block:
				switch r {
				case '"':
					fallthrough
				case '\'':
					currentState = state_normal
					arguments = append(arguments, string(buffer))
					buffer = buffer[:0]
				default:
					buffer = append(buffer, r)
				}
			}
		}

		if len(buffer) > 0 {
			arguments = append(arguments, string(buffer))
		}

		err := app.argsApp.Run(arguments)

		switch err {
		case cmd.ErrNotSaveThisToHistory:
			continue
		default:
			app.appendHistory(input)
		}
	}

	app.saveHistoryFile()
	log.Info("bye :)")
}

func (app *App) RegisterCmd(command cmd.Command) {
	if command == nil {
		log.Erro("register failed command obj is nil")
		return
	}

	if len(command.Name()) == 0 {
		log.Erro("register failed command must have a name")
		return
	}

	if _, ok := app.commandMap[command.Name()]; ok {
		log.Erro("register failed command %v is already exist,please check", command.Name())
		return
	}

	if err := command.Init(app); err != nil {
		log.Erro("register failed command %v init error : %v", command.Name(), err)
		app.Exit()
	}

	app.commandMap[command.Name()] = command
}

func (app *App) init() {
	app.argsApp = args.NewApp()
	app.argsApp.Name = SYSTEM_NAME
	app.argsApp.Usage = SYSTEM_NAME
	app.argsApp.Version = VERSION
	app.argsApp.Flags = []args.Flag{
		&args.BoolFlag{
			Name:  "once",
			Usage: "run once then exit",
		},
	}

	app.argsApp.After = func(c *args.Context) (err error) {
		app.runOnce = c.Bool("once")
		return
	}

	app.argsApp.CommandNotFound = func(c *args.Context, input string) {
		log.Erro("can't find command %v,please see 'help' for more information", input)
	}

	for _, command := range app.commandMap {
		if command == nil {
			continue
		}

		action := command.Action()

		if action == nil {
			continue
		}

		app.argsApp.Commands = append(app.argsApp.Commands, action.Command)
	}
}

func (app *App) loadHistoryFile() {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Debu("[app] load history faild,can't get user home dir")
		return
	}

	historyFileName := path.Join(homeDir, APP_HISTORY_FILENAME)

	file, err := os.OpenFile(historyFileName, os.O_RDONLY, 0600)

	if err != nil || file == nil {
		log.Debu("[app] load history file %v failed %v", historyFileName, err)
		return
	}

	defer file.Close()

	history, err := ioutil.ReadAll(file)

	if err != nil {
		log.Debu("[app] load history file %v failed %v", historyFileName, err)
		return
	}

	if len(history) == 0 {
		return
	}

	app.history = strings.Split(string(history), "\n")

	log.Debu("[app] load history success,record %v", len(app.history))
}

func (app *App) saveHistoryFile() {
	if len(app.history) == 0 {
		return
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		log.Debu("[app] save history faild,can't get user home dir")
		return
	}

	historyFileName := path.Join(homeDir, APP_HISTORY_FILENAME)

	file, err := os.OpenFile(historyFileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)

	if err != nil || file == nil {
		log.Debu("[app] save history file %v failed %v", historyFileName, err)
		return
	}

	defer file.Close()

	_, _ = file.WriteString(app.history[0])

	for i := 1; i < len(app.history); i++ {
		_, _ = file.WriteString("\n")
		_, _ = file.WriteString(app.history[i])
	}

	log.Debu("[app] save history success,record %v", len(app.history))
}

func (app *App) appendHistory(input string) {
	if len(input) == 0 {
		return
	}

	app.history = append(app.history, input)

	if len(app.history) > 1000 {
		app.history = app.history[1:]
	}
}

func (app *App) activeHistory(pt *prompt.Prompt) {
	if pt == nil {
		return
	}

	if fixhistory := prompt.OptionHistory(app.history); fixhistory != nil {
		_ = fixhistory(pt)
	}

	cmd.SetHistory(app.history)
}
