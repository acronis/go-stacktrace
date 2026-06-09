package stacktrace

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Type is the type of the error.
type Type string

// String is a fmt.Stringer implementation.
func (t *Type) String() string {
	if t == nil {
		return ""
	}
	return string(*t)
}

// Severity is the severity of the error.
type Severity string

// String is a fmt.Stringer implementation.
func (s *Severity) String() string {
	if s == nil {
		return ""
	}
	return string(*s)
}

// stringer is a fmt.Stringer implementation.
type stringer struct {
	msg string
}

// String implements the fmt.Stringer interface.
func (s *stringer) String() string {
	if s == nil {
		return ""
	}
	return s.msg
}

// Stringer returns a fmt.Stringer for the given value.
func Stringer(v interface{}) fmt.Stringer {
	switch w := v.(type) {
	case fmt.Stringer:
		return w
	case string:
		return &stringer{msg: w}
	case *string:
		return &stringer{msg: *w}
	case error:
		return &stringer{msg: w.Error()}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &stringer{msg: fmt.Sprintf("%d", w)}
	case float32, float64:
		return &stringer{msg: fmt.Sprintf("%f", w)}
	case bool:
		return &stringer{msg: fmt.Sprintf("%t", w)}
	case nil:
		return &stringer{msg: "nil"}
	default:
		return &stringer{msg: fmt.Sprintf("%v", w)}
	}
}

// StructInfo is a map of string keys to fmt.Stringer values.
// It is used to store additional information about an error.
// WARNING: Not thread-safe
type StructInfo struct {
	info map[string]fmt.Stringer
}

// String implements the fmt.Stringer interface.
// It returns a string representation of the struct info.
func (s *StructInfo) String() string {
	var result string
	keys := s.SortedKeys()

	for _, k := range keys {
		v, ok := s.info[k]
		if ok {
			if result == "" {
				result = fmt.Sprintf("%s: %s", k, v)
			} else {
				result = fmt.Sprintf("%s: %s: %s", result, k, v)
			}
		}
	}
	return result
}

// ensureMap ensures that the map is initialized.
func (s *StructInfo) ensureMap() {
	if s.info == nil {
		s.info = make(map[string]fmt.Stringer)
	}
}

// Add adds a key-value pair to the struct info.
func (s *StructInfo) Add(key string, value fmt.Stringer) *StructInfo {
	s.ensureMap()
	s.info[key] = value
	return s
}

// Get returns the value of the given key.
func (s *StructInfo) Get(key string) fmt.Stringer {
	s.ensureMap()
	return s.info[key]
}

// StringBy returns the string value of the given key.
func (s *StructInfo) StringBy(key string) string {
	s.ensureMap()
	return s.info[key].String()
}

// Remove removes the given key from the struct info.
func (s *StructInfo) Remove(key string) *StructInfo {
	s.ensureMap()
	delete(s.info, key)
	return s
}

// Has checks if the given key exists in the struct info.
func (s *StructInfo) Has(key string) bool {
	s.ensureMap()
	_, ok := s.info[key]
	return ok
}

// Keys returns the keys of the struct info.
func (s *StructInfo) Keys() []string {
	s.ensureMap()
	result := make([]string, 0, len(s.info))
	for k := range s.info {
		result = append(result, k)
	}
	return result
}

// SortedKeys returns the sorted keys of the struct info.
func (s *StructInfo) SortedKeys() []string {
	s.ensureMap()
	keys := s.Keys()
	sort.Strings(keys)
	return keys
}

// Update updates the struct info with the given struct info.
func (s *StructInfo) Update(u *StructInfo) *StructInfo {
	s.ensureMap()
	for k, v := range u.info {
		s.info[k] = v
	}
	return s
}

// NewStructInfo creates a new struct info.
func NewStructInfo() *StructInfo {
	return &StructInfo{
		info: make(map[string]fmt.Stringer),
	}
}

type Location string

func (loc *Location) String() string {
	if loc == nil {
		return ""
	}
	return string(*loc)
}

// MaxList is the maximum number of entries allowed in a StackTrace's List.
// When the limit is reached, a sentinel entry is appended and further entries
// are dropped. A value of 0 means unlimited.
var MaxList = 10

// MaxDepth is the maximum Wrapped-chain depth. When the incoming node's depth
// equals or exceeds this limit, the chain is replaced by a sentinel leaf.
// A value of 0 means unlimited.
var MaxDepth = 2000

// TooManyErrorsMsg is the message of the sentinel entry added to List when
// MaxList is exceeded. Callers can test for this to detect truncation.
const TooManyErrorsMsg = "too many errors, some have been omitted"

// ChainTruncatedMsg is the message of the sentinel leaf substituted for the
// Wrapped chain when MaxDepth is exceeded.
const ChainTruncatedMsg = "error chain truncated"

// StackTrace contains information about a parser error.
type StackTrace struct {
	// Severity is the severity of the error.
	Severity *Severity
	// Type is the type of the error.
	Type *Type
	// Location is the location file path of the error.
	Location *Location
	// Position is the position of the error in the file.
	Position *Position

	// Wrapped is the error that wrapped by this error.
	Wrapped *StackTrace
	// Err is the underlying error. It is not used for the error message.
	Err error
	// Message is the error message.
	Message string
	// Info is the additional information about the error.
	Info StructInfo

	// List is the list of stack traces.
	List []*StackTrace

	typeIsSet bool
	// depth is the length of the Wrapped chain rooted at this node (0 = leaf).
	depth int
}

// Header returns the header of the StackTrace.
func (st *StackTrace) Header() string {
	segs := make([]string, 0)
	result := st.Type.String()
	if result != "" {
		segs = append(segs, result)
	}
	loc := st.GetLocWithPos()
	if loc != "" {
		segs = append(segs, loc)
	}
	return strings.Join(segs, ": ")
}

// Option is an option for the StackTrace creation.
type Option interface {
	Apply(*StackTrace)
}

// OrigString returns the original error message without the wrapping messages.
func (st *StackTrace) OrigString() string {
	segs := make([]string, 0)
	header := st.Header()
	if header != "" {
		segs = append(segs, header)
	}
	msg := st.MessageWithInfo()
	if msg != "" {
		segs = append(segs, msg)
	}
	return strings.Join(segs, ": ")
}

// OrigStringW returns the original error message with the wrapped OrigStringW
func (st *StackTrace) OrigStringW() string {
	segs := make([]string, 0)
	orig := st.OrigString()
	if orig != "" {
		segs = append(segs, orig)
	}
	if st.Wrapped != nil {
		segs = append(segs, st.Wrapped.OrigStringW())
	}
	return strings.Join(segs, ": ")
}

func (st *StackTrace) MessageWithInfo() string {
	segments := make([]string, 0)
	if st.Message != "" {
		segments = append(segments, st.Message)
	}
	if len(st.Info.info) > 0 {
		segments = append(segments, st.Info.String())
	}
	return strings.Join(segments, ": ")
}

// String implements the fmt.Stringer interface.
// It returns the string representation of the StackTrace.
func (st *StackTrace) String() string {
	segs := make([]string, 0)
	orig := st.OrigString()
	if orig != "" {
		segs = append(segs, orig)
	}
	if st.Wrapped != nil {
		segs = append(segs, st.Wrapped.String())
	}

	res := strings.Join(segs, ": ")

	if len(st.List) > 0 {
		lists := make([]string, 0)
		lists = append(lists, res)
		for _, elem := range st.List {
			lists = append(lists, elem.String())
		}
		res = strings.Join(lists, "; ")
	}
	return res
}

// StackTrace implements the error interface.
// It returns the string representation of the StackTrace.
func (st *StackTrace) Error() string {
	return st.String()
}

// Unwrap checks if the given error is a StackTrace and returns it.
// It returns false if the error is not a StackTrace.
// if err is a StackTrace, it returns the wrapped StackTrace and true.
// if err is not a StackTrace, it is unwrapped and wrapped with a new StackTrace if it has a wrapped StackTrace.
func Unwrap(err error) (*StackTrace, bool) {
	if err == nil {
		return nil, false
	}
	wrapped, ok := err.(*StackTrace)
	if !ok {
		errWrapped := errors.Unwrap(err)
		if errWrapped == nil {
			return nil, false
		}
		wrapped, ok = Unwrap(errWrapped)
		if ok && wrapped != nil {
			msg, _, _ := strings.Cut(err.Error(), wrapped.Message)
			msg = strings.TrimSuffix(msg, ": ")
			wrapped = New(msg).Wrap(wrapped)
			wrapped.Err = err
		}
	}
	if wrapped == nil {
		return nil, false
	}

	return wrapped, ok
}

// Is checks if the given error is the same as the StackTrace.
func (st *StackTrace) Is(err error) bool {
	if st == nil {
		return false
	}
	if err == nil {
		return false
	}
	wrapped, ok := Unwrap(err)
	if !ok {
		return false
	}
	if st == wrapped {
		return true
	}
	return st.Is(wrapped.Wrapped)
}

// New creates a new StackTrace.
func New(message string, opts ...Option) *StackTrace {
	e := &StackTrace{
		Message: message,
	}
	for _, opt := range opts {
		opt.Apply(e)
	}
	return e
}

type optErrInfo struct {
	Key   string
	Value fmt.Stringer
}

func (o optErrInfo) Apply(e *StackTrace) {
	e.Info.Add(o.Key, o.Value)
}

type optErrPosition struct {
	Pos *Position
}

func (o optErrPosition) Apply(e *StackTrace) {
	e.Position = o.Pos
}

type optErrLocation struct {
	Location Location
}

func (o optErrLocation) Apply(e *StackTrace) {
	e.Location = &o.Location
}

type optErrSeverity struct {
	Severity Severity
}

func (o optErrSeverity) Apply(e *StackTrace) {
	e.Severity = &o.Severity
}

type optErrType struct {
	ErrType Type
}

func (o optErrType) Apply(e *StackTrace) {
	_ = e.SetType(o.ErrType)
}

func WithInfo(key string, value any) Option {
	return optErrInfo{Key: key, Value: Stringer(value)}
}

func WithPosition(pos *Position) Option {
	return optErrPosition{Pos: pos}
}

func WithSeverity(severity Severity) Option {
	return optErrSeverity{Severity: severity}
}

// WithType sets the type of the error with override.
func WithType(errType Type) Option {
	return optErrType{ErrType: errType}
}

// WithLocation sets the location of the error.
func WithLocation(location string) Option {
	return optErrLocation{Location: Location(location)}
}

// NewWrapped creates a new StackTrace from the given go error.
func NewWrapped(message string, err error, opts ...Option) *StackTrace {
	if wrapped, ok := Unwrap(err); ok {
		return New(
			message,
			opts...,
		).Wrap(wrapped).SetErr(wrapped.Err)
	}
	return New(fmt.Sprintf("%s: %s", message, err.Error()), opts...).SetErr(err)
}

// Wrap wraps the given error with the StackTrace if it is not a StackTrace.
// It returns the wrapped StackTrace.
func Wrap(err error, opts ...Option) *StackTrace {
	if st, ok := Unwrap(err); ok {
		for _, opt := range opts {
			opt.Apply(st)
		}
		return st
	}
	return New(err.Error(), opts...).SetErr(err)
}

// SetSeverity sets the severity of the StackTrace and returns it
func (st *StackTrace) SetSeverity(severity Severity) *StackTrace {
	st.Severity = &severity
	return st
}

// SetType sets the type of the StackTrace and returns it, operation can be done only once.
func (st *StackTrace) SetType(t Type) *StackTrace {
	if !st.typeIsSet {
		st.Type = &t
		st.typeIsSet = true
	}
	if st.Wrapped != nil {
		_ = st.Wrapped.SetType(t)
	}
	return st
}

// SetLocation sets the location of the StackTrace and returns it
func (st *StackTrace) SetLocation(location string) *StackTrace {
	loc := Location(location)
	st.Location = &loc
	return st
}

// SetPosition sets the position of the StackTrace and returns it
func (st *StackTrace) SetPosition(pos *Position) *StackTrace {
	st.Position = pos
	return st
}

// SetMessage sets the message of the StackTrace and returns it
func (st *StackTrace) SetMessage(message string, a ...any) *StackTrace {
	st.Message = fmt.Sprintf(message, a...)
	return st
}

// SetErr sets the underlying error of the StackTrace and returns it
func (st *StackTrace) SetErr(err error) *StackTrace {
	st.Err = err
	return st
}

// Wrap wraps the given StackTrace and returns it.
// If MaxDepth > 0 and w's chain depth already equals or exceeds MaxDepth,
// the chain is replaced by a sentinel leaf that carries the message, location,
// and position of the deepest reachable node so the caller still sees the
// root cause rather than only ChainTruncatedMsg.
func (st *StackTrace) Wrap(w *StackTrace) *StackTrace {
	if w == nil {
		return st
	}
	if MaxDepth > 0 && w.depth >= MaxDepth {
		sentinel := &StackTrace{Message: ChainTruncatedMsg}
		if deepest := deepestWrapped(w); deepest != nil && deepest.Message != "" {
			sentinel.Message = deepest.Message
			sentinel.Location = deepest.Location
			sentinel.Position = deepest.Position
		}
		st.Wrapped = sentinel
		st.depth = 1
		return st
	}
	st.Wrapped = w
	st.depth = w.depth + 1
	return st
}

// deepestWrapped returns the last node reachable by following Wrapped links.
func deepestWrapped(st *StackTrace) *StackTrace {
	for st.Wrapped != nil {
		st = st.Wrapped
	}
	return st
}

// sameFingerprint reports whether two nodes share the same (location, position,
// message, info) identity for deduplication purposes.
// Checks are ordered cheapest-first so the common non-matching case short-circuits
// before reaching Info.String(), which is the only allocation.
func (st *StackTrace) sameFingerprint(o *StackTrace) bool {
	if st.Message != o.Message {
		return false
	}
	switch {
	case st.Location == nil && o.Location == nil:
	case st.Location == nil || o.Location == nil:
		return false
	case *st.Location != *o.Location:
		return false
	}
	switch {
	case st.Position == nil && o.Position == nil:
	case st.Position == nil || o.Position == nil:
		return false
	case st.Position.Line != o.Position.Line || st.Position.Column != o.Position.Column:
		return false
	}
	if len(st.Info.info) == 0 && len(o.Info.info) == 0 {
		return true
	}
	return st.Info.String() == o.Info.String()
}

// Append adds the given StackTrace to the list of StackTraces and returns it.
// Duplicate entries — those whose (location, position, message, info) fingerprint
// matches an existing entry — are silently dropped.
// When MaxList > 0 and the list is full, a TooManyErrorsMsg sentinel is appended
// in place of the new entry and all subsequent entries are dropped.
func (st *StackTrace) Append(e *StackTrace) *StackTrace {
	if MaxList > 0 && len(st.List) > MaxList {
		// Sentinel already placed; drop silently.
		return st
	}
	for _, existing := range st.List {
		if existing.sameFingerprint(e) {
			return st
		}
	}
	if st.List == nil {
		st.List = make([]*StackTrace, 0)
	}
	if MaxList > 0 && len(st.List) == MaxList {
		// Capacity reached: replace with sentinel and stop.
		st.List = append(st.List, &StackTrace{Message: TooManyErrorsMsg})
		return st
	}
	st.List = append(st.List, e)
	return st
}

// GetLocWithPos returns the location with position of the StackTrace.
func (st *StackTrace) GetLocWithPos() string {
	res := st.GetLocWithPosPtr()
	if res == nil {
		return ""
	}
	return *res
}

// GetLocWithPosPtr returns the location with position of the StackTrace.
func (st *StackTrace) GetLocWithPosPtr() *string {
	if st.Location == nil {
		return nil
	}
	result := st.Location.String()
	if result != "" {
		result = fmt.Sprintf("%s:%s", result, st.Position)
	}
	return &result
}

// Position contains the line and column where the error occurred.
type Position struct {
	Line      int
	Column    int
	EndLine   int
	EndColumn int // 1-based exclusive (column after the last character); when set, used to construct a multi-character LSP range
}

func (p *Position) String() string {
	if p == nil {
		return "1"
	}
	if p.Line == 0 {
		return "1"
	}
	result := fmt.Sprintf("%d", p.Line)
	if p.Column > 0 {
		result = fmt.Sprintf("%s:%d", result, p.Column)
	}
	return result
}

// NewPosition creates a new position with the given line and column.
func NewPosition(line, column int) *Position {
	return &Position{Line: line, Column: column}
}

// Accumulator collects multiple *StackTrace errors into a single List-bearing
// root node. The zero value is ready to use without any initialization.
//
// It replaces the repetitive nil-check-then-Append pattern:
//
//	var acc stacktrace.Accumulator
//	for _, item := range items {
//		if err := process(item); err != nil {
//			acc.Add(stacktrace.NewWrapped("step", err))
//		}
//	}
//	return acc.Result()
type Accumulator struct {
	head *StackTrace
}

// Add appends st to the accumulator. A nil st is ignored.
func (a *Accumulator) Add(st *StackTrace) {
	if st == nil {
		return
	}
	if a.head == nil {
		a.head = st
	} else {
		a.head = a.head.Append(st)
	}
}

// Result returns the accumulated StackTrace, or nil when no errors were added.
func (a *Accumulator) Result() *StackTrace {
	return a.head
}

// optErrNodePosition is an option to set the position of the error to the position of the given node.
type optErrNodePosition struct {
	pos *Position
}

// Apply sets the position of the error to the given position.
// implements Option
func (o optErrNodePosition) Apply(e *StackTrace) {
	e.Position = o.pos
}
