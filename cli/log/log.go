/*
	自定义日志结构
*/
package log

import "github.com/nerored/chat-test-golang/log"

var (
	isDebugMod bool = false
)

func Init() {
	log.InitLog("")
}

func DebugMod() bool {
	return isDebugMod
}

func SetDebugMod(open bool) {
	isDebugMod = open
}

func Debu(format string, args ...interface{}) {
	if !isDebugMod {
		return
	}

	log.Ulog(log.LOG_LEVEL_DEBU, log.PRINT_LEVELAB, format, args...)
}

func DebuS(format string, args ...interface{}) {
	if !isDebugMod {
		return
	}

	log.Ulog(log.LOG_LEVEL_DEBU, log.PRINT_LEVELAB|log.PRINT_STACKIN|5, format, args...)
}

func Info(format string, args ...interface{}) {
	log.Ulog(log.LOG_LEVEL_INFO, log.PRINT_DEFAULT, format, args...)
}

func Warn(format string, args ...interface{}) {
	log.Ulog(log.LOG_LEVEL_WARN, log.PRINT_DEFAULT, format, args...)
}

func Erro(format string, args ...interface{}) {
	flag := log.PRINT_DEFAULT

	if isDebugMod {
		flag = log.PRINT_UTRACE
	}

	log.Ulog(log.LOG_LEVEL_ERRO, flag, format, args...)
}
