package stacktrace

import (
	"fmt"
	"log/slog"
)

func ErrToSlogAttr(err error, opts ...TracesOpt) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	st, ok := Unwrap(err)
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
				slog.String("type", string(stack.Type)),
				slog.String("severity", string(stack.Severity)),
				slog.String("position", stack.LinePos),
				slog.String("message", stack.Message),
			)
			stackAttrs = append(stackAttrs, stackAttr)
		}
		tracebackAttr := slog.Group(fmt.Sprintf("%d", traceIndex), "stack", stackAttrs)
		tracebackAttrs = append(tracebackAttrs, tracebackAttr)
	}

	return slog.Group("tracebacks", "traces", tracebackAttrs)
}
