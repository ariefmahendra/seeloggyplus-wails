package logger

import (
	"github.com/rs/zerolog"
	"path"
	"runtime"
)

const FunctionFieldName = "function"

type FunctionHook struct {
}

func (h FunctionHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	// 3 is skip value, for caller function name
	pc, _, _, ok := runtime.Caller(3)
	if ok {
		fullFuncName := runtime.FuncForPC(pc).Name()
		shortFuncName := path.Base(fullFuncName)
		e.Str(FunctionFieldName, shortFuncName)
	}
}
