/*
	导出调用的接口
	note:Ulog 为用户自定格式日志，自定参数参考参数列表
	Prefix 为终端输入提供支持（不换行)
*/
package log

func Debu(format string, args ...interface{}) {
	if !sharedPrinter.debugprint {
		return
	}

	sharedPrinter.print(LOG_LEVEL_DEBU, PRINT_DEBUG, format, args...)
}

func Trac(format string, args ...interface{}) {
	sharedPrinter.print(LOG_LEVEL_TRAC, PRINT_UTRACE, format, args...)
}

func Info(format string, args ...interface{}) {
	sharedPrinter.print(LOG_LEVEL_INFO, PRINT_DEFAULT, format, args...)
}

func Warn(format string, args ...interface{}) {
	sharedPrinter.print(LOG_LEVEL_WARN, PRINT_DEFINE, format, args...)
}

func Erro(format string, args ...interface{}) {
	sharedPrinter.print(LOG_LEVEL_ERRO, PRINT_DEFINE, format, args...)
}

func Fata(format string, args ...interface{}) {
	sharedPrinter.print(LOG_LEVEL_FATA, PRINT_DEFINE, format, args...)
}

func Ulog(level LogLv, flags PrintFlag, format string, args ...interface{}) {
	sharedPrinter.print(level, flags, format, args...)
}

func Prefix(format string, args ...interface{}) {
	sharedPrinter.printNoEnter(LOG_LEVEL_INFO, PRINT_DEFAULT, format, args...)
}
