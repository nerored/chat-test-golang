/*
	最大堆，用于词频统计
	每5秒刷新一次
	reflush 后 可以在任何时候调用top 获取最受欢迎的词
*/
package main

import (
	"sync"
)

type heapNode struct {
	num  int64
	word string
}

func (hn *heapNode) greater(other *heapNode) bool {
	if other == nil {
		return true
	}

	return hn.num > other.num
}

type maxHeap struct {
	n    int
	heap []*heapNode
	menu map[string]*heapNode

	sync.RWMutex
}

func newMaxHeap() *maxHeap {
	return &maxHeap{
		menu: make(map[string]*heapNode),
	}
}

func (mh *maxHeap) record(word string) {
	mh.Lock()
	defer mh.Unlock()

	node := mh.menu[word]

	if node != nil {
		node.num++
		return
	}

	mh.menu[word] = &heapNode{
		num:  1,
		word: word,
	}
}

func (mh *maxHeap) top() (word string) {
	mh.RLock()
	defer mh.RUnlock()

	if mh.len() <= 0 {
		return
	}

	return mh.heap[0].word
}

func (mh *maxHeap) reflush() {
	mh.Lock()
	defer mh.Unlock()

	mh.n = 0
	mh.heap = make([]*heapNode, len(mh.menu))

	for _, node := range mh.menu {
		if node == nil {
			continue
		}

		mh.push(node)
	}

	mh.menu = make(map[string]*heapNode)
}

//----------------------------------------

func (mh *maxHeap) cap() int {
	return len(mh.heap)
}

func (mh *maxHeap) len() int {
	return mh.n
}

func (mh *maxHeap) isFull() bool {
	return mh.len() >= mh.cap()
}

func (mh *maxHeap) push(node *heapNode) {
	if node == nil || mh.isFull() {
		return
	}

	i := mh.len()
	for ; i != 0 && node.greater(mh.heap[i/2]); i /= 2 {
		mh.heap[i] = mh.heap[i/2]
	}

	mh.n++
	mh.heap[i] = node
}
