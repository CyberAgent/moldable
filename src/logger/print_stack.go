package logger

import (
	"log"
	"runtime"
)

/** https://go.dev/play/p/DkGku2xlhSr */
func PrintStack() {
	var pc [100]uintptr
	n := runtime.Callers(0, pc[:])
	frames := runtime.CallersFrames(pc[:n])
	var (
		fr runtime.Frame
		ok bool
	)
	if _, ok = frames.Next(); !ok {
		return
	}
	for ok {
		fr, ok = frames.Next()
		if !ok {
			return
		}
		log.Default().Println(fr.Function, fr.File, fr.Line)
	}
}
