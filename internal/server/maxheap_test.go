package main

import (
	"testing"

	"github.com/nerored/chat-test-golang/log"
)

var inputs = []string{
	"hello world",
	"game begin",
	"good game",
	"i wanna play this game",
	"log great game",
	"fuck",
	"how are you",
	"i'm fine,thank you,and you",
	"li lei and han meimei",
	"你好，吗",
}

func TestMaxHeap(t *testing.T) {
	log.InitLog("")
	heap := newMaxHeap()

	for _, msg := range inputs {
		words := splitWords(msg)

		for _, word := range words {
			heap.record(word)
		}
	}

	for _, node := range heap.menu {
		log.Info("node %+v", node)
	}

	heap.reflush()

	log.Info("heap size %v", heap.len())

	for i, node := range heap.heap {
		if i >= heap.len() {
			break
		}

		log.Info("node %+v", node)
	}

	mostpopular := heap.top()

	t.Logf("most popular %v", mostpopular)
}
