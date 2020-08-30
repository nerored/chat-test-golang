package main

import (
	"testing"
)

var testwords = []string{
	"abc",
	"hello",
	"你好",
	"口吐芬芳",
}

var tryreplace = []string{
	"abc123321123",
	"acbacb2121acbabc",
	"你123123好",
	"你好你很好真的你好lol",
	"口来吐了口吐芬芳",
}

var results = []string{
	"***123321123",
	"acbacb2121acb***",
	"你123123好",
	"**你很好真的**lol",
	"口来吐了****",
}

func TestTrieTest(t *testing.T) {
	root := buildTrie(testwords...)

	for i, word := range tryreplace {
		repStr := root.replace(word, '*')

		if repStr != results[i] {
			t.Fatalf("faild at %v %v", word, i)
		}
	}
}

func BenchmarkTrieTest(b *testing.B) {
	root := buildTrie(testwords...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tIndex := i % len(tryreplace)

		src := tryreplace[tIndex]
		des := results[tIndex]

		if root.replace(src, '*') != des {
			b.Fatalf("faild at %v", i)
		}
	}
}
