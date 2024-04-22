package logger

import (
	"github.com/rs/zerolog"
)

type ContextHook struct{}

func (h ContextHook) Run(e *zerolog.Event, level zerolog.Level, _ string) {
	ctx := e.GetCtx()

	serviceCtx, _ := ctx.Value(ServiceContextInfo).(ServiceContext)

	e.Interface("serviceContext", serviceCtx)
}
