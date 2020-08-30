/*
	日志格式指定
	1.打印时间lab
	2.打印文件定位信息
	3.打印日志分级lab
	4.打印堆栈信息以及指定堆栈深度（最多10）
*/
package log

type PrintFlag int16

//print flags define
// [using- stack depth
// -int16] ↑      ↑
//00000000 00000000
//   ||||↓
//   |||↓time lable
//   ||↓file location
//   |↓func name
//   ↓stack info
//   level lable
const (
	PRINT_DEFAULT PrintFlag = 0x0000
	PRINT_TIMELAB PrintFlag = 0x0100
	PRINT_FILELOC PrintFlag = 0x0200
	PRINT_FUNCNAM PrintFlag = 0x0400
	PRINT_STACKIN PrintFlag = 0x0800
	PRINT_LEVELAB PrintFlag = 0x1000

	PRINT_STACKDP PrintFlag = 0x00FF

	PRINT_DEBUG  = PRINT_DEFAULT
	PRINT_DEFINE = PRINT_TIMELAB | PRINT_LEVELAB | PRINT_FILELOC | PRINT_FUNCNAM
	PRINT_UTRACE = PRINT_DEFINE | PRINT_STACKIN | 10
)
