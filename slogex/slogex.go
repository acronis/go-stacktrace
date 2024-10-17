package slogex

import (
	"fmt"
	"log/slog"

	"github.com/acronis/go-stacktrace"
)

func ErrToSlogAttr(err error, opts ...stacktrace.TracesOpt) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	st, ok := stacktrace.Unwrap(err)
	if !ok {
		return slog.String("error", err.Error())
	}

	if st == nil {
		return slog.String("error", "nil stacktrace")
	}

	tracebacks := st.GetTraces(opts...)

	tracebackAttrs := make([]slog.Attr, 0, len(tracebacks))
	for traceIndex := range tracebacks {
		traceback := tracebacks[traceIndex]
		stackAttrs := make([]slog.Attr, 0, len(traceback.Stack))
		for stackIndex := range traceback.Stack {
			stack := traceback.Stack[stackIndex]
			key := fmt.Sprintf("%d", stackIndex)
			stackAttr := slog.Group(
				key,
				func() []any {
					attrs := []any{}
					if stack.Type != nil {
						attrs = append(attrs, slog.String("type", stack.Type.String()))
					}
					if stack.Severity != nil {
						attrs = append(attrs, slog.String("severity", stack.Severity.String()))
					}
					if stack.LinePos != nil {
						attrs = append(attrs, slog.String("position", *stack.LinePos))
					}
					attrs = append(attrs, slog.String("message", stack.Message))
					return attrs
				}()...,
			)
			stackAttrs = append(stackAttrs, stackAttr)
		}
		tracebackAttr := slog.Group(fmt.Sprintf("%d", traceIndex), "stack", stackAttrs)
		tracebackAttrs = append(tracebackAttrs, tracebackAttr)
	}

	return slog.Group("tracebacks", "traces", tracebackAttrs)
}
