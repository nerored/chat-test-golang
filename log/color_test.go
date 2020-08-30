package log

import (
	"fmt"
	"strings"
	"testing"
)

func TestColorSet(t *testing.T) {
	for _, colorAttr := range []int{
		FGC_DEFAULT,
		FGC_BLACK,
		FGC_RED,
		FGC_GREEN,
		FGC_YELLOW,
		FGC_BLUE,
		FGC_MAGENTA,
		FGC_CYAN,
		FGC_LIGHTGREY,
		FGC_DARKGREY,
		FGC_LIGHTRED,
		FGC_LIGHTGREEN,
		FGC_LIGHTYELLOW,
		FGC_LIGHTBLUE,
		FGC_LIGHTMAGENTA,
		FGC_LIGHTCYAN,
		FGC_LIGHTWHITE,
	} {
		buffer := new(strings.Builder)

		if buffer == nil {
			t.Fatalf("can't create buffer")
		}

		obj := new(color)

		if obj == nil {
			t.Fatalf("can't create color obj")
		}

		obj.setAttrs(colorAttr)
		obj.begin(buffer)
		buffer.WriteString("吾有一言。曰「問天地好在」。")
		obj.end(buffer)
		fmt.Println(buffer.String())
	}
}

func TestFormatSet(t *testing.T) {
	for _, colorAttr := range []int{
		FMT_BOLD,
		FMT_DIM,
		FMT_UNDERLINED,
		FMT_BLINK,
		FMT_MINVERTED,
		FMT_HIDDEN,
	} {
		buffer := new(strings.Builder)

		if buffer == nil {
			t.Fatalf("can't create buffer")
		}

		obj := new(color)

		if obj == nil {
			t.Fatalf("can't create color obj")
		}

		obj.setAttrs(colorAttr)
		obj.begin(buffer)
		buffer.WriteString("落霞与孤鹜齐飞，秋水共长天一色")
		obj.end(buffer)
		fmt.Println(buffer.String())
	}
}

func TestBGColorSet(t *testing.T) {
	for _, colorAttr := range []int{
		BGC_DEFAULT,
		BGC_BLACK,
		BGC_RED,
		BGC_GREEN,
		BGC_YELLOW,
		BGC_BLUE,
		BGC_MAGENTA,
		BGC_CYAN,
		BGC_LIGHTGREY,
		BGC_DARKGREY,
		BGC_LIGHTRED,
		BGC_LIGHTGREEN,
		BGC_LIGHTYELLOW,
		BGC_LIGHTBLUE,
		BGC_LIGHTMAGENTA,
		BGC_LIGHTCYAN,
		BGC_LIGHTWHITE,
	} {
		buffer := new(strings.Builder)

		if buffer == nil {
			t.Fatalf("can't create buffer")
		}

		obj := new(color)

		if obj == nil {
			t.Fatalf("can't create color obj")
		}

		obj.setAttrs(colorAttr)
		obj.begin(buffer)
		buffer.WriteString("仿佛兮若轻云之蔽月,飘飘兮若流风之回雪")
		obj.end(buffer)
		fmt.Println(buffer.String())
	}
}
