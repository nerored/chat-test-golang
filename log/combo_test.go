package log

import (
	"testing"
)

func TestTermLog(t *testing.T) {
	Info("对酒当歌，人生几何")
	Warn("譬如朝露，去日苦多")
	Debu("慨当以慷，忧思难忘")
	Erro("何以%v？唯%v杜康", NewCombo("解忧", FGC_BLUE), NewCombo("有", BGC_YELLOW, FGC_LIGHTMAGENTA))
	Trac("%v子衿，悠悠我心", NewCombo("青青", FMT_BLINK))
	Fata("但为君故，沉吟至今")
	Info("呦呦%v之苹", NewCombo("鹿鸣，食野", FGC_MAGENTA, FMT_UNDERLINED))
	Warn("我有嘉宾，鼓瑟吹笙")
	Debu("明明如月，何时可掇")
	Erro("忧从中来，不可断绝")
}
