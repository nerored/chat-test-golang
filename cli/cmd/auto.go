/*
	自动补全支持：
	1.解析注册过来的cmd结构体，生成其action以及相关参数的补全信息
	2.历史输入补全支持
*/
package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"path"
	"reflect"
	"strings"

	"github.com/c-bata/go-prompt"
)

type AutoComplete struct {
	preActionSuggests   []prompt.Suggest
	historiesSuggests   []string
	fieldNameSuggestMap map[string][]prompt.Suggest
	shortNameSuggestMap map[string][]prompt.Suggest
}

var (
	sharedAutoComple = AutoComplete{
		fieldNameSuggestMap: make(map[string][]prompt.Suggest),
		shortNameSuggestMap: make(map[string][]prompt.Suggest),
	}
)

func SetHistory(histories []string) {
	found := func(s []string, a string) bool {
		for _, i := range s {
			if i == a {
				return true
			}
		}

		return false
	}

	for _, history := range histories {
		if found(sharedAutoComple.historiesSuggests, history) {
			continue
		}

		sharedAutoComple.historiesSuggests = append(sharedAutoComple.historiesSuggests, history)
	}
}

func Completer(d prompt.Document) []prompt.Suggest {
	return append(sharedAutoComple.completer(d), sharedAutoComple.getHistoryMatch(d)...)
}

func (ac *AutoComplete) register(name, usage string, refT reflect.Type) (err error) {
	if len(name) <= 0 || refT.Kind() != reflect.Struct {
		return errors.New("autocomplete register failed name is nil or kind is not a struct")
	}

	for _, suggest := range ac.preActionSuggests {
		if suggest.Text == name {
			return fmt.Errorf("autocomplete command %v is already registered", name)
		}
	}

	var fieldSuggestList []prompt.Suggest
	var shortSuggestList []prompt.Suggest

	for i := 0; i < refT.NumField(); i++ {
		field := refT.Field(i)

		if field.Anonymous {
			continue
		}

		fieldSuggestList = append(fieldSuggestList, prompt.Suggest{
			Text:        fmt.Sprintf("--%v", strings.ToLower(field.Name)),
			Description: field.Tag.Get("usage"),
		})

		shortName := field.Tag.Get("short")

		if len(shortName) <= 0 {
			continue
		}

		shortSuggestList = append(shortSuggestList, prompt.Suggest{
			Text:        fmt.Sprintf("-%v", shortName),
			Description: field.Tag.Get("usage"),
		})
	}

	ac.preActionSuggests = append(ac.preActionSuggests, prompt.Suggest{
		Text:        name,
		Description: usage,
	})

	ac.fieldNameSuggestMap[name] = fieldSuggestList
	ac.shortNameSuggestMap[name] = shortSuggestList

	return nil
}

func (ac *AutoComplete) completer(d prompt.Document) (result []prompt.Suggest) {
	if len(d.Text) == 0 {
		return
	}

	commandTag := strings.IndexByte(d.Text, ' ')

	if commandTag < 0 {
		return prompt.FilterHasPrefix(ac.preActionSuggests, d.GetWordBeforeCursor(), true)
	}

	command := d.Text[:commandTag]

	paramWord := d.GetWordBeforeCursorUntilSeparator("= ")

	switch {
	case strings.HasPrefix(paramWord, "--"):
		return prompt.FilterHasPrefix(ac.fieldNameSuggestMap[command], paramWord, true)
	case strings.HasPrefix(paramWord, "-"):
		return prompt.FilterHasPrefix(ac.shortNameSuggestMap[command], paramWord, true)
	case strings.HasPrefix(paramWord, "."):
		fallthrough
	case strings.HasPrefix(paramWord, "/"):
		dir, _ := path.Split(paramWord)

		command := exec.Command("ls", dir)

		if command == nil {
			return
		}

		output, err := command.Output()

		if err != nil {
			return
		}

		var suggestList []prompt.Suggest

		for _, fileName := range strings.Split(string(output), "\n") {
			suggestList = append(suggestList, prompt.Suggest{
				Text: strings.Join([]string{dir, fileName}, ""),
			})
		}

		return prompt.FilterHasPrefix(suggestList, paramWord, true)
	}

	return
}

func (ac *AutoComplete) getHistoryMatch(d prompt.Document) (result []prompt.Suggest) {
	if len(d.Text) == 0 {
		return
	}

	for _, history := range ac.historiesSuggests {
		if !strings.HasPrefix(history, d.Text) {
			continue
		}

		splitIndex := strings.LastIndex(d.Text, " ") + 1

		if splitIndex < 0 {
			splitIndex = len(d.Text)
		}

		record := strings.TrimPrefix(history, d.Text[:splitIndex])

		if len(record) == 0 {
			continue
		}

		result = append(result, prompt.Suggest{
			Text:        record,
			Description: "history",
		})
	}

	return
}
