/*
	输出器：
	1.解析日志参数并填写对应信息
	2.多源输出,以及终端格式输出可选
*/
package log

import (
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type repeater struct {
	colored, nocolor *writerChan
}

func (r *repeater) Write(p []byte) (n int, err error) {
	if r.colored != nil {
		n, err = r.colored.Write(p)
	}

	if r.nocolor != nil {
		n, err = r.nocolor.Write(p)
	}

	return
}

type Writenew func() Writer

type printer struct {
	timeFormat string
	debugprint bool
	bufcreator []Writenew
}

var (
	sharedPrinter printer
)

func InitLog(timeFormat string) {
	if len(timeFormat) > 0 {
		SetTimeFormat(timeFormat)
	} else {
		SetTimeFormat("2006-01-02 15:04:05")
	}

	SetDebugPrint(true)
	AddBufCreator(newTermWriter)
}

func IsDebugMod() bool {
	return sharedPrinter.debugprint
}

func SetTimeFormat(format string) {
	sharedPrinter.timeFormat = format
}

func SetDebugPrint(isOpen bool) {
	sharedPrinter.debugprint = isOpen
}

func AddBufCreator(f ...Writenew) {
	if len(f) == 0 {
		return
	}

	sharedPrinter.bufcreator = append(sharedPrinter.bufcreator, f...)
}

func (p *printer) printHeadInfo(writer io.Writer, logLevel LogLv, flags PrintFlag) {
	if writer == nil {
		return
	}

	p.printTimeInfo(writer, flags)

	var prefix, suffix string

	if flags&PRINT_LEVELAB > 0 {
		switch logLevel {
		case LOG_LEVEL_DEBU:
			prefix = "[debu"
			suffix = "] "
		case LOG_LEVEL_INFO:
			prefix = "[info"
			suffix = "] "
		case LOG_LEVEL_TRAC:
			prefix = "[trac"
			suffix = "] "
		case LOG_LEVEL_WARN:
			prefix = "[warn"
			suffix = "] "
		case LOG_LEVEL_ERRO:
			prefix = "[erro"
			suffix = "] "
		case LOG_LEVEL_FATA:
			prefix = "[fata"
			suffix = "] "
		default:
			if flags&PRINT_DEFINE > 0 {
				prefix = "["
				suffix = "] "
			}
		}
	}

	if len(prefix) > 0 {
		fmt.Fprint(writer, prefix)
	}
	p.printFileInfo(writer, flags)
	p.printFuncName(writer, flags)
	if len(suffix) > 0 {
		fmt.Fprint(writer, suffix)
	}
}

func (p *printer) printTimeInfo(writer io.Writer, flags PrintFlag) {
	if writer == nil || flags&PRINT_TIMELAB <= 0 {
		return
	}

	fmt.Fprint(writer, time.Now().Format(p.timeFormat))
}

const LOC_STACK_DEPTH = 4

func (p *printer) printFileInfo(writer io.Writer, flags PrintFlag) {
	if writer == nil || flags&PRINT_FILELOC <= 0 {
		return
	}

	_, fileName, line, ok := runtime.Caller(LOC_STACK_DEPTH)

	if !ok {
		return
	}

	fmt.Fprint(writer, " ")
	fmt.Fprint(writer, fileName)
	fmt.Fprint(writer, ":")
	fmt.Fprint(writer, strconv.Itoa(line))
}

func (p *printer) printFuncName(writer io.Writer, flags PrintFlag) {
	if writer == nil || flags&PRINT_FUNCNAM <= 0 {
		return
	}

	pc, _, _, ok := runtime.Caller(LOC_STACK_DEPTH)

	if !ok {
		return
	}

	if funcPC := runtime.FuncForPC(pc); funcPC != nil {
		fmt.Fprint(writer, " ")
		fmt.Fprint(writer, funcPC.Name())
	}
}

const STD_STACK_DEPTH = 3

func (p *printer) printStackDep(writer io.Writer, flags PrintFlag) {
	if writer == nil || flags&PRINT_STACKIN <= 0 {
		return
	}

	depLen := int(flags & PRINT_STACKDP)

	if depLen <= 0 {
		return
	}

	fmt.Fprintf(writer, "-------- stack info\n")
	for add := 0; add < depLen; add++ {
		pc, fileName, line, ok := runtime.Caller(add + STD_STACK_DEPTH)

		if !ok {
			continue
		}

		fmt.Fprintf(writer, "\t%d: %s:%d ", add, fileName, line)

		funcPC := runtime.FuncForPC(pc)

		if funcPC == nil {
			fmt.Fprint(writer, "unknow")
		} else {
			fmt.Fprint(writer, funcPC.Name())
		}

		fmt.Fprintf(writer, "\n")
	}
}

func (p *printer) reflush(coloredIO, nocolorIO []Writer) {
	for _, writer := range coloredIO {
		if writer == nil {
			continue
		}

		writer.reflush()
	}

	for _, writer := range nocolorIO {
		if writer == nil {
			continue
		}

		writer.reflush()
	}
}

func (p *printer) newIOList() (coloredIO, nocolorIO []Writer) {
	for _, creator := range p.bufcreator {
		if creator == nil {
			continue
		}

		writer := creator()

		if writer == nil {
			continue
		}

		if writer.needColor() {
			coloredIO = append(coloredIO, writer)
		} else {
			nocolorIO = append(nocolorIO, writer)
		}
	}

	return
}

func (p *printer) print(logLevel LogLv, flags PrintFlag, format string, args ...interface{}) {
	if len(format) <= 0 && len(args) <= 0 {
		return
	}

	coloredIO, nocolorIO := p.newIOList()
	defer p.reflush(coloredIO, nocolorIO)

	majorW := getMajorChan(logLevel, coloredIO, nocolorIO)
	defer freeMajorChan(&majorW)

	root := newRoot(logLevel)

	if len(args) > 0 {
		root.makeChan(!majorW.colored.isEmpty(), args)
	}

	root.begin(majorW.colored)
	defer root.end(majorW.colored)

	p.printHeadInfo(&majorW, logLevel, flags)

	fmt.Fprintf(majorW.colored, format, args...)
	fmt.Fprintf(majorW.nocolor, format, args...)

	freeCombos(args)

	if !strings.HasSuffix(format, "\n") {
		fmt.Fprintf(&majorW, "\n")
	}

	p.printStackDep(&majorW, flags)
}

func (p *printer) printNoEnter(logLevel LogLv, flags PrintFlag, format string, args ...interface{}) {
	if len(format) <= 0 && len(args) <= 0 {
		return
	}

	coloredIO, nocolorIO := p.newIOList()
	defer p.reflush(coloredIO, nocolorIO)

	majorW := getMajorChan(logLevel, coloredIO, nocolorIO)
	defer freeMajorChan(&majorW)

	root := newRoot(logLevel)

	if len(args) > 0 {
		root.makeChan(!majorW.colored.isEmpty(), args)
	}

	root.begin(majorW.colored)
	defer root.end(majorW.colored)

	p.printHeadInfo(&majorW, logLevel, flags)

	fmt.Fprintf(majorW.colored, format, args...)
	fmt.Fprintf(majorW.nocolor, format, args...)

	freeCombos(args)

	p.printStackDep(&majorW, flags)
}
