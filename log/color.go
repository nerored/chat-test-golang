/*
	终端输出格式枚举
	1.通过将数据包含在color中，则在print时会将
	color指定的终端格式一起输出,而无需通过手动字符拼接
*/
package log

import (
	"fmt"
	"io"
)

const (
	FMT_RESET = 0

	FMT_BOLD       = 1
	FMT_DIM        = 2
	FMT_UNDERLINED = 4
	FMT_BLINK      = 5
	FMT_MINVERTED  = 7
	FMT_HIDDEN     = 8

	FGC_DEFAULT      = 39
	FGC_BLACK        = 30
	FGC_RED          = 31
	FGC_GREEN        = 32
	FGC_YELLOW       = 33
	FGC_BLUE         = 34
	FGC_MAGENTA      = 35
	FGC_CYAN         = 36
	FGC_LIGHTGREY    = 37
	FGC_DARKGREY     = 90
	FGC_LIGHTRED     = 91
	FGC_LIGHTGREEN   = 92
	FGC_LIGHTYELLOW  = 93
	FGC_LIGHTBLUE    = 94
	FGC_LIGHTMAGENTA = 95
	FGC_LIGHTCYAN    = 96
	FGC_LIGHTWHITE   = 97

	BGC_DEFAULT      = 49
	BGC_BLACK        = 40
	BGC_RED          = 41
	BGC_GREEN        = 42
	BGC_YELLOW       = 43
	BGC_BLUE         = 44
	BGC_MAGENTA      = 45
	BGC_CYAN         = 46
	BGC_LIGHTGREY    = 47
	BGC_DARKGREY     = 100
	BGC_LIGHTRED     = 101
	BGC_LIGHTGREEN   = 102
	BGC_LIGHTYELLOW  = 103
	BGC_LIGHTBLUE    = 104
	BGC_LIGHTMAGENTA = 105
	BGC_LIGHTCYAN    = 106
	BGC_LIGHTWHITE   = 107
)

type color struct {
	attributes []int
	colorChan  *color
}

func (c *color) isEmpty() bool {
	return len(c.attributes) == 0 &&
		(c.colorChan == nil || c.colorChan.isEmpty())
}

func (c *color) setAttrs(attrs ...int) {
	c.attributes = attrs
}

func (c *color) linkTo(preChan *color) {
	c.colorChan = preChan
}

func (c *color) writeHead(writer io.Writer) {
	if writer == nil || len(c.attributes) == 0 {
		return
	}

	fmt.Fprintf(writer, "\x1b[")
	defer fmt.Fprintf(writer, "m")

	fmt.Fprintf(writer, "%d", c.attributes[0])

	for i := 1; i < len(c.attributes); i++ {
		fmt.Fprintf(writer, ";%d", c.attributes[i])
	}
}

func (c *color) writeTail(writer io.Writer) {
	if writer == nil || len(c.attributes) == 0 {
		return
	}

	fmt.Fprintf(writer, "\x1b[%dm", FMT_RESET)
}

func (c *color) begin(writer io.Writer) {
	if writer == nil || c.isEmpty() {
		return
	}

	//clear format before
	if c.colorChan != nil {
		c.colorChan.end(writer)
	}

	c.writeHead(writer)
}

func (c *color) end(writer io.Writer) {
	if writer == nil || c.isEmpty() {
		return
	}

	c.writeTail(writer)

	//recover format before
	if c.colorChan != nil {
		c.colorChan.begin(writer)
	}
}
