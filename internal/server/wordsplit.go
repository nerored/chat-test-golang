package main

import (
	"strings"
	"unicode"
)

func splitWords(msg string) (words []string) {
	var buffer strings.Builder

	for _, r := range msg {
		switch {
		case unicode.IsLetter(r):
			buffer.WriteRune(r)
		default:
			if buffer.Len() > 0 {
				words = append(words, buffer.String())
				buffer.Reset()
			}
		}
	}

	if buffer.Len() > 0 {
		words = append(words, buffer.String())
	}

	return
}
