/*
	数据与颜色的组合器,负责将数据与指定的终端格式组合
	可以通过link来继承父节点的格式
*/
package log

import (
	"fmt"
	"strings"
	"sync"
)

type Combo struct {
	color
	coloredCount int
	data         interface{}
}

var comboPool = sync.Pool{
	New: func() interface{} {
		return new(Combo)
	},
}

func NewCombo(data interface{}, attrs ...int) (c *Combo) {
	c = comboPool.Get().(*Combo)

	c.data = data
	c.color.setAttrs(attrs...)
	return
}

func (c *Combo) linkTo(comboChan *Combo) {
	if comboChan == nil {
		return
	}

	c.color.linkTo(&comboChan.color)
	c.addColoredChange(comboChan.coloredCount)
}

func (c *Combo) free() {
	c.data = nil
	c.coloredCount = 0
	c.color.attributes = nil
	c.color.colorChan = nil

	comboPool.Put(c)
}

func (c *Combo) addColoredChange(delta int) {
	c.coloredCount += delta
}

func (c *Combo) makeChan(needColor bool, args []interface{}) {
	if len(args) == 0 {
		return
	}

	if needColor {
		c.addColoredChange(1)
	}

	for _, arg := range args {
		chanObj, ok := arg.(*Combo)

		if !ok || chanObj == nil {
			continue
		}

		chanObj.linkTo(c)
	}
}

func (c *Combo) String() string {
	writing := func() (buffer *strings.Builder) {
		buffer = new(strings.Builder)

		c.begin(buffer)
		defer c.end(buffer)

		fmt.Fprintf(buffer, "%v", c.data)
		return
	}

	return writing().String()
}

func freeCombos(args []interface{}) {
	for _, arg := range args {
		combo, ok := arg.(*Combo)

		if !ok || combo == nil {
			continue
		}

		combo.free()
	}
}

//------ combo root maker

func newRoot(logLevel LogLv) (root Combo) {
	switch logLevel {
	case LOG_LEVEL_DEBU:
		root.setAttrs(FGC_LIGHTCYAN)
	case LOG_LEVEL_INFO:
		root.setAttrs(FGC_DEFAULT)
	case LOG_LEVEL_TRAC:
		root.setAttrs(FGC_LIGHTYELLOW, FMT_UNDERLINED)
	case LOG_LEVEL_WARN:
		root.setAttrs(FGC_YELLOW)
	case LOG_LEVEL_ERRO:
		root.setAttrs(FGC_RED)
	case LOG_LEVEL_FATA:
		root.setAttrs(FGC_LIGHTWHITE, BGC_RED)
	}

	return
}
