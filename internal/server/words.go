/*
   加载脏词文件
*/
package main

import (
	"io/ioutil"
	"strings"

	"github.com/nerored/chat-test-golang/log"
)

func loadProfanityWordsFromFile(path string) *trieNode {
	words, err := ioutil.ReadFile(path)

	if err != nil {
		log.Erro("load profanitywords failed %v", err)
		return nil
	}

	return buildTrie(strings.Split(string(words), "\n")...)
}
