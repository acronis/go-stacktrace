package stacktrace

type TracesOpt interface {
	Apply(o *TracesOptions)
}

type TracesOptions struct {
	// EnsureDuplicates ensures that duplicates are not printed
	EnsureDuplicates bool
	dupLocs          map[string]struct{}
}

func NewTracesOptions() *TracesOptions {
	opts := &TracesOptions{
		EnsureDuplicates: false,
		dupLocs:          make(map[string]struct{}),
	}
	return opts
}

type ensureDuplicatesOpt struct{}

func (ensureDuplicatesOpt) Apply(o *TracesOptions) {
	o.EnsureDuplicates = true
}

func WithEnsureDuplicates() TracesOpt {
	return &ensureDuplicatesOpt{}
}

type Stack struct {
	LinePos  *string
	Severity *Severity
	Message  string
	Type     *Type
}

func NewStack() *Stack {
	return &Stack{}
}

type Trace struct {
	Stack []Stack
}

func NewTrace() *Trace {
	return &Trace{Stack: make([]Stack, 0)}
}

func (st *StackTrace) getTraces(opts *TracesOptions) []Trace {
	traces := make([]Trace, 0)

	tracesWithList := func() []Trace {
		for _, elem := range st.List {
			elemTraces := elem.getTraces(opts)
			traces = append(traces, elemTraces...)
		}
		return traces
	}

	trace := NewTrace()
	stack := NewStack()
	stack.LinePos = st.GetLocWithPosPtr()
	stack.Severity = st.Severity
	stack.Message = st.MessageWithInfo()
	stack.Type = st.Type

	if stack.LinePos != nil {
		if _, ok := opts.dupLocs[*stack.LinePos]; ok {
			return tracesWithList()
		}
	}

	trace.Stack = append(trace.Stack, *stack)
	if st.Wrapped != nil {
		wrappedTraces := st.Wrapped.getTraces(opts)
		if len(wrappedTraces) == 0 {
			return tracesWithList()
		}
		for i := range wrappedTraces {
			trace.Stack = append(trace.Stack, wrappedTraces[i].Stack...)
		}
	} else if opts.EnsureDuplicates && stack.LinePos != nil {
		opts.dupLocs[*stack.LinePos] = struct{}{}
	}

	traces = append(traces, *trace)

	return tracesWithList()
}

func (st *StackTrace) GetTraces(opts ...TracesOpt) []Trace {
	o := NewTracesOptions()
	for _, opt := range opts {
		opt.Apply(o)
	}
	return st.getTraces(o)
}
