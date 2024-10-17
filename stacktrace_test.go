package stacktrace

import (
	"fmt"
	"reflect"
	"testing"
)

func TestStringer(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check stringer with string",
			args: args{
				v: "string",
			},
			want: "string",
		},
		{
			name: "Check stringer with int",
			args: args{
				v: 10,
			},
			want: "10",
		},
		{
			name: "Check stringer with stringer",
			args: args{
				v: Stringer("stringer"),
			},
			want: "stringer",
		},
		{
			name: "Check stringer with nil",
			args: args{
				v: nil,
			},
			want: "nil",
		},
		{
			name: "Check stringer with error",
			args: args{
				v: fmt.Errorf("error"),
			},
			want: "error",
		},
		{
			name: "Check stringer with pointer of string",
			args: args{
				v: func() *string { s := "string"; return &s }(),
			},
			want: "string",
		},
		{
			name: "Check stringer with bool",
			args: args{
				v: true,
			},
			want: "true",
		},
		{
			name: "Check stringer with float",
			args: args{
				v: 10.1,
			},
			want: "10.100000",
		},
		{
			name: "Check stringer with struct",
			args: args{
				v: struct {
					key string
				}{key: "value"},
			},
			want: "{value}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Stringer(tt.args.v); !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("Stringer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_String(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check empty struct info",
			fields: fields{
				info: map[string]fmt.Stringer{},
			},
			want: "",
		},
		{
			name: "Check struct info with one key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			want: "key: value",
		},
		{
			name: "Check struct info with two keys",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key1": Stringer("value1"),
					"key2": Stringer("value2"),
				},
			},
			want: "key1: value1: key2: value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_Add(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		key   string
		value fmt.Stringer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StructInfo
	}{
		{
			name: "Check add key",
			fields: fields{
				info: map[string]fmt.Stringer{},
			},
			args: args{
				key:   "key",
				value: Stringer("value"),
			},
			want: &StructInfo{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.Add(tt.args.key, tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_Get(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fmt.Stringer
	}{
		{
			name: "Check get key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				key: "key",
			},
			want: Stringer("value"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_StringBy(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Check string by key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				key: "key",
			},
			want: "value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.StringBy(tt.args.key); got != tt.want {
				t.Errorf("StringBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_Remove(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StructInfo
	}{
		{
			name: "Check remove key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				key: "key",
			},
			want: &StructInfo{
				info: map[string]fmt.Stringer{},
			},
		},
		{
			name: "Check remove key from empty struct info",
			fields: fields{
				info: map[string]fmt.Stringer{},
			},
			args: args{
				key: "key",
			},
			want: &StructInfo{
				info: map[string]fmt.Stringer{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.Remove(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_Has(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Check has key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				key: "key",
			},
			want: true,
		},
		{
			name: "Check has key in empty struct info",
			fields: fields{
				info: map[string]fmt.Stringer{},
			},
			args: args{
				key: "key",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.Has(tt.args.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStructInfo_Update(t *testing.T) {
	type fields struct {
		info map[string]fmt.Stringer
	}
	type args struct {
		u *StructInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StructInfo
	}{
		{
			name: "Check update struct info",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				u: &StructInfo{
					info: map[string]fmt.Stringer{
						"key2": Stringer("value2"),
					},
				},
			},
			want: &StructInfo{
				info: map[string]fmt.Stringer{
					"key":  Stringer("value"),
					"key2": Stringer("value2"),
				},
			},
		},
		{
			name: "Check update struct info with the same key",
			fields: fields{
				info: map[string]fmt.Stringer{
					"key": Stringer("value"),
				},
			},
			args: args{
				u: &StructInfo{
					info: map[string]fmt.Stringer{
						"key": Stringer("value2"),
					},
				},
			},
			want: &StructInfo{
				info: map[string]fmt.Stringer{
					"key": Stringer("value2"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StructInfo{
				info: tt.fields.info,
			}
			if got := s.Update(tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_SetSeverity(t *testing.T) {
	type fields struct {
		Severity Severity
	}
	type args struct {
		severity Severity
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StackTrace
	}{
		{
			name:   "Check set severity",
			fields: fields{},
			args: args{
				severity: Severity("critical"),
			},
			want: &StackTrace{
				Severity: func() *Severity { s := Severity("critical"); return &s }(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Severity: &tt.fields.Severity,
			}
			if got := v.SetSeverity(tt.args.severity); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_SetType(t *testing.T) {
	type fields struct {
		ErrType Type
	}
	type args struct {
		errType Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StackTrace
	}{
		{
			name:   "Check set type",
			fields: fields{},
			args: args{
				errType: Type("parsing"),
			},
			want: &StackTrace{
				Type: func() *Type { tt := Type("parsing"); return &tt }(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Type: &tt.fields.ErrType,
			}
			if got := v.SetType(tt.args.errType); !reflect.DeepEqual(got.Type, tt.want.Type) {
				t.Errorf("SetType() = %v, want %v", got.Type, tt.want.Type)
			}
		})
	}
}

func TestError_SetLocationAndPosition(t *testing.T) {
	type fields struct {
		Location Location
		Position Position
	}
	type args struct {
		location string
		pos      Position
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StackTrace
	}{
		{
			name: "Check set location and position",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: Position{
					Line:   10,
					Column: 1,
				},
			},
			args: args{
				location: "/usr/local/raml2.raml",
				pos: Position{
					Line:   20,
					Column: 2,
				},
			},
			want: &StackTrace{
				Location: func() *Location { l := Location("/usr/local/raml2.raml"); return &l }(),
				Position: &Position{
					Line:   20,
					Column: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Location: &tt.fields.Location,
				Position: &tt.fields.Position,
			}
			if got := v.SetLocation(tt.args.location).SetPosition(&tt.args.pos); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_SetMessage(t *testing.T) {
	type fields struct {
		Message string
	}
	type args struct {
		message string
		a       []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *StackTrace
	}{
		{
			name: "Check set message",
			fields: fields{
				Message: "message",
			},
			args: args{
				message: "new message",
				a:       []any{},
			},
			want: &StackTrace{
				Message: "new message",
			},
		},
		{
			name: "Check set message with arguments",
			fields: fields{
				Message: "message",
			},
			args: args{
				message: "new message with %s",
				a:       []any{"argument"},
			},
			want: &StackTrace{
				Message: "new message with argument",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Message: tt.fields.Message,
			}
			if got := v.SetMessage(tt.args.message, tt.args.a...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	type fields struct {
		Severity Severity
		ErrType  Type
		Location Location
		Position *Position
		Wrapped  *StackTrace
		Err      error
		Message  string
		Info     StructInfo
		List     []*StackTrace
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check error",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: &Position{Line: 10, Column: 1},
				Message:  "message",
			},
			want: "/usr/local/raml.raml:10:1: message",
		},
		{
			name: "Check error without position",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Message:  "message",
			},
			want: "/usr/local/raml.raml:1: message",
		},
		{
			name: "Check error with info",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: &Position{Line: 10, Column: 1},
				Message:  "message",
				Info:     *NewStructInfo().Add("key", Stringer("value")),
			},
			want: "/usr/local/raml.raml:10:1: message: key: value",
		},
		{
			name: "Check error with empty message",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: &Position{Line: 10, Column: 1},
			},
			want: "/usr/local/raml.raml:10:1",
		},
		{
			name: "Check error with only info",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: &Position{Line: 10, Column: 1},
				Info:     *NewStructInfo().Add("key", Stringer("value")),
			},
			want: "/usr/local/raml.raml:10:1: key: value",
		},
		{
			name: "Check error with only wrapped error",
			fields: fields{
				Location: "/usr/local/raml.raml",
				Position: &Position{Line: 10, Column: 1},
				Wrapped: &StackTrace{
					Message: "message 1",
					Wrapped: &StackTrace{
						Message: "message 2",
					},
				},
				Message: "message",
			},
			want: "/usr/local/raml.raml:10:1: message: message 1: message 2",
		},
		{
			name: "Check error with list",
			fields: fields{
				Message: "message",
				Wrapped: &StackTrace{
					Message: "message 1",
				},
				List: []*StackTrace{
					{
						Message: "message 2",
						Wrapped: &StackTrace{
							Message: "message 3",
						},
					},
					{
						Message: "message 4",
						Wrapped: &StackTrace{
							Message: "message 5",
						},
					},
				},
			},
			want: "message: message 1; message 2: message 3; message 4: message 5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Severity: &tt.fields.Severity,
				Type:     &tt.fields.ErrType,
				Location: &tt.fields.Location,
				Position: tt.fields.Position,
				Wrapped:  tt.fields.Wrapped,
				Err:      tt.fields.Err,
				Message:  tt.fields.Message,
				Info:     tt.fields.Info,
				List:     tt.fields.List,
			}
			if got := v.Error(); got != tt.want {
				t.Errorf("StackTrace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestError_OrigString(t *testing.T) {
	var SeverityError Severity = "error"
	var TypeValidating Type = "validating"
	var loc Location = "/usr/local/raml.raml"
	var loc2 Location = "/usr/local/raml2.raml"

	type fields struct {
		Severity        *Severity
		ErrType         *Type
		Location        *Location
		Position        Position
		Wrapped         *StackTrace
		Err             error
		Message         string
		WrappedMessages string
		Info            StructInfo
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check original string",
			fields: fields{
				Severity: &SeverityError,
				ErrType:  &TypeValidating,
				Location: &loc,
				Position: Position{Line: 10, Column: 1},
				Message:  "message",
				Info:     *NewStructInfo().Add("key", Stringer("value")),
				Wrapped: &StackTrace{
					Severity: &SeverityError,
					Type:     &TypeValidating,
					Location: &loc2,
					Position: &Position{Line: 20, Column: 2},
					Message:  "wrapped",
					Info:     *NewStructInfo().Add("key", Stringer("value")),
				},
			},
			want: "validating: /usr/local/raml.raml:10:1: message: key: value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &StackTrace{
				Severity: tt.fields.Severity,
				Type:     tt.fields.ErrType,
				Location: tt.fields.Location,
				Position: &tt.fields.Position,
				Wrapped:  tt.fields.Wrapped,
				Err:      tt.fields.Err,
				Message:  tt.fields.Message,
				Info:     tt.fields.Info,
			}
			if got := v.OrigString(); got != tt.want {
				t.Errorf("OrigString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWrappedError(t *testing.T) {
	var SeverityCritical Severity = "critical"
	var TypeParsing Type = "parsing"
	var loc Location = "/usr/local/raml.raml"

	type args struct {
		message string
		err     error
		opts    []Option
	}
	tests := []struct {
		name string
		args args
		want *StackTrace
	}{
		{
			name: "Check wrapped error",
			args: args{
				message: "message",
				err:     fmt.Errorf("error"),
				opts: []Option{
					WithLocation("/usr/local/raml.raml"),
					WithSeverity(SeverityCritical),
					WithPosition(NewPosition(10, 1)),
					WithInfo("key", Stringer("value")),
					WithType(TypeParsing),
				},
			},
			want: &StackTrace{
				Severity:  &SeverityCritical,
				Type:      &TypeParsing,
				Message:   "message: error",
				Location:  &loc,
				Position:  &Position{Line: 10, Column: 1},
				Err:       fmt.Errorf("error"),
				Info:      *NewStructInfo().Add("key", Stringer("value")),
				typeIsSet: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWrapped(tt.args.message, tt.args.err, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWrapped() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	var SeverityCritical Severity = "critical"
	var TypeParsing Type = "parsing"
	var loc Location = "/usr/local/raml.raml"

	type args struct {
		message  string
		location string
		opts     []Option
	}
	tests := []struct {
		name string
		args args
		want *StackTrace
	}{
		{
			name: "Check error",
			args: args{
				message: "message",
				opts: []Option{
					WithLocation("/usr/local/raml.raml"),
					WithSeverity(SeverityCritical),
					WithPosition(NewPosition(10, 1)),
					WithInfo("key", Stringer("value")),
					WithType(TypeParsing),
				},
			},
			want: &StackTrace{
				Severity:  &SeverityCritical,
				Type:      &TypeParsing,
				Message:   "message",
				Location:  &loc,
				Position:  &Position{Line: 10, Column: 1},
				Info:      *NewStructInfo().Add("key", Stringer("value")),
				typeIsSet: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.message, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    *Type
		want string
	}{
		{
			name: "Check type no nil",
			t:    func() *Type { tt := Type("type"); return &tt }(),
			want: "type",
		},
		{
			name: "Check type nil",
			t:    nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSeverity_String(t *testing.T) {
	tests := []struct {
		name string
		s    *Severity
		want string
	}{
		{
			name: "Check severity no nil",
			s:    func() *Severity { s := Severity("severity"); return &s }(),
			want: "severity",
		},
		{
			name: "Check severity nil",
			s:    nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stringer_String(t *testing.T) {
	type fields struct {
		stringer *stringer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check stringer not nil",
			fields: fields{
				stringer: &stringer{
					msg: "message",
				},
			},
			want: "message",
		},
		{
			name: "Check stringer nil",
			fields: fields{
				stringer: nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.fields.stringer
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLocation_String(t *testing.T) {
	tests := []struct {
		name string
		loc  *Location
		want string
	}{
		{
			name: "Check location",
			loc:  func() *Location { l := Location("location"); return &l }(),
			want: "location",
		},
		{
			name: "Check location nil",
			loc:  nil,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.loc.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStackTrace_OrigStringW(t *testing.T) {
	type fields struct {
		Severity  *Severity
		Type      *Type
		Location  *Location
		Position  *Position
		Wrapped   *StackTrace
		Err       error
		Message   string
		Info      StructInfo
		List      []*StackTrace
		typeIsSet bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check original string",
			fields: fields{
				Severity: func() *Severity { s := Severity("severity"); return &s }(),
				Type:     func() *Type { t := Type("type"); return &t }(),
				Location: func() *Location { l := Location("location"); return &l }(),
				Position: &Position{Line: 10, Column: 1},
				Message:  "message",
				Info:     *NewStructInfo().Add("key", Stringer("value")),
				Wrapped: &StackTrace{
					Message: "wrapped",
				},
			},
			want: "type: location:10:1: message: key: value: wrapped",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &StackTrace{
				Severity:  tt.fields.Severity,
				Type:      tt.fields.Type,
				Location:  tt.fields.Location,
				Position:  tt.fields.Position,
				Wrapped:   tt.fields.Wrapped,
				Err:       tt.fields.Err,
				Message:   tt.fields.Message,
				Info:      tt.fields.Info,
				List:      tt.fields.List,
				typeIsSet: tt.fields.typeIsSet,
			}
			if got := st.OrigStringW(); got != tt.want {
				t.Errorf("OrigStringW() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want func(t *testing.T, got *StackTrace, got1 bool)
	}{
		{
			name: "Check unwrap: nil",
			args: args{
				err: nil,
			},
			want: func(t *testing.T, got *StackTrace, got1 bool) {
				if got != nil {
					t.Errorf("Unwrap() = %v, want %v", got, nil)
				}
				if got1 {
					t.Errorf("Unwrap() = %v, want %v", got1, false)
				}
			},
		},
		{
			name: "Check unwrap",
			args: args{
				err: fmt.Errorf("error"),
			},
			want: func(t *testing.T, got *StackTrace, got1 bool) {
				if got != nil {
					t.Errorf("Unwrap() = %v, want %v", got, nil)
				}
				if got1 {
					t.Errorf("Unwrap() = %v, want %v", got1, false)
				}
			},
		},
		{
			name: "Check unwrap with stack trace",
			args: args{
				err: &StackTrace{Message: "message"},
			},
			want: func(t *testing.T, got *StackTrace, got1 bool) {
				if got == nil {
					t.Errorf("Unwrap() = %v, want %v", got, nil)
				}
				if !got1 {
					t.Errorf("Unwrap() = %v, want %v", got1, true)
				}
				if got.Message != "message" {
					t.Errorf("Unwrap() = %v, want %v", got.Message, "message")
				}
			},
		},
		{
			name: "Check unwrap with wrapped error",
			args: args{
				err: fmt.Errorf("error: %w", &StackTrace{
					Message: "message",
					Wrapped: &StackTrace{
						Message: "wrapped",
					},
				}),
			},
			want: func(t *testing.T, got *StackTrace, got1 bool) {
				if got == nil {
					t.Errorf("Unwrap() = %v, want %v", got, nil)
				}
				if !got1 {
					t.Errorf("Unwrap() = %v, want %v", got1, true)
				}
				if got.Message != "error" {
					t.Errorf("Unwrap() = %v, want %v", got.Message, "error")
				}
				if got.Wrapped == nil {
					t.Errorf("Unwrap() = %v, want %v", got.Wrapped, nil)
				}
				if got.Wrapped.Message != "message" {
					t.Errorf("Unwrap() = %v, want %v", got.Wrapped.Message, "message")
				}
				if got.Wrapped.Wrapped == nil {
					t.Errorf("Unwrap() = %v, want %v", got.Wrapped.Wrapped, nil)
				}
				if got.Wrapped.Wrapped.Message != "wrapped" {
					t.Errorf("Unwrap() = %v, want %v", got.Wrapped.Wrapped.Message, "wrapped")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Unwrap(tt.args.err)
			if tt.want != nil {
				tt.want(t, got, got1)
			}
		})
	}
}

func TestStackTrace_Is(t *testing.T) {
	checkErr := &StackTrace{Message: "message"}
	diffErr := &StackTrace{Message: "message"}
	type fields struct {
		st *StackTrace
	}
	type args struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Check is",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: checkErr,
			},
			want: true,
		},
		{
			name: "Check is with wrapped error",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: fmt.Errorf("error: %w", checkErr),
			},
			want: true,
		},
		{
			name: "Check is with wrapped stacktrace",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: fmt.Errorf("error: %w", &StackTrace{
					Message: "message",
					Wrapped: checkErr,
				}),
			},
			want: true,
		},
		{
			name: "Negative: Check is with different error",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: diffErr,
			},
			want: false,
		},
		{
			name: "Negative: Check is with wrapped different error",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: fmt.Errorf("error: %w", diffErr),
			},
			want: false,
		},
		{
			name: "Negative: Check is with wrapped stacktrace with different error",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: fmt.Errorf("error: %w", &StackTrace{
					Message: "message",
					Wrapped: diffErr,
				}),
			},
			want: false,
		},
		{
			name: "Negative: Check is with nil",
			fields: fields{
				st: nil,
			},
			args: args{
				err: checkErr,
			},
			want: false,
		},
		{
			name: "Negative: Check is with arg nil",
			fields: fields{
				st: checkErr,
			},
			args: args{
				err: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := tt.fields.st
			if got := st.Is(tt.args.err); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewWrapped(t *testing.T) {
	type args struct {
		message string
		err     error
		opts    []Option
	}
	tests := []struct {
		name string
		args args
		want func(t *testing.T, got *StackTrace)
	}{
		{
			name: "Check with simple error",
			args: args{
				message: "message",
				err:     fmt.Errorf("error"),
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "message: error" {
					t.Errorf("NewWrapped() = %v, want %v", got.Message, "message: error")
				}
				if got.Err == nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Err, nil)
				}
			},
		},
		{
			name: "Check with wrapped error",
			args: args{
				message: "message",
				err:     fmt.Errorf("error: %w", &StackTrace{Message: "wrapped"}),
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "message" {
					t.Errorf("NewWrapped() = %v, want %v", got.Message, "message")
				}
				if got.Err == nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Err, nil)
				}
				if got.Wrapped == nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped, nil)
				}
				if got.Wrapped.Message != "error" {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped.Message, "error")
				}
				if got.Wrapped.Wrapped == nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped.Wrapped, nil)
				}
				if got.Wrapped.Wrapped.Message != "wrapped" {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped.Wrapped.Message, "wrapped")
				}
			},
		},
		{
			name: "Check with wrapped stacktrace",
			args: args{
				message: "message",
				err: &StackTrace{
					Message: "error",
				},
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "message" {
					t.Errorf("NewWrapped() = %v, want %v", got.Message, "message")
				}
				if got.Err != nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Err, nil)
				}
				if got.Wrapped == nil {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped, nil)
				}
				if got.Wrapped.Message != "error" {
					t.Errorf("NewWrapped() = %v, want %v", got.Wrapped.Message, "error")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWrapped(tt.args.message, tt.args.err, tt.args.opts...)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		err  error
		opts []Option
	}
	tests := []struct {
		name string
		args args
		want func(t *testing.T, got *StackTrace)
	}{
		{
			name: "Check wrap",
			args: args{
				err: fmt.Errorf("error"),
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "error" {
					t.Errorf("Wrap() = %v, want %v", got.Message, "error")
				}
				if got.Err == nil {
					t.Errorf("Wrap() = %v, want %v", got.Err, nil)
				}
			},
		},
		{
			name: "Check wrap stacktrace",
			args: args{
				err: &StackTrace{
					Message: "message",
				},
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "message" {
					t.Errorf("Wrap() = %v, want %v", got.Message, "message")
				}
				if got.Err != nil {
					t.Errorf("Wrap() = %v, want %v", got.Err, nil)
				}
			},
		},
		{
			name: "Check wrap stacktrace with options",
			args: args{
				err: &StackTrace{Message: "message"},
				opts: []Option{
					WithLocation("/usr/local/raml.raml"),
				},
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Message != "message" {
					t.Errorf("Wrap() = %v, want %v", got.Message, "message")
				}
				if got.Err != nil {
					t.Errorf("Wrap() = %v, want %v", got.Err, nil)
				}
				if got.Location == nil {
					t.Errorf("Wrap() = %v, want %v", got.Location, nil)
				}
				if *got.Location != "/usr/local/raml.raml" {
					t.Errorf("Wrap() = %v, want %v", *got.Location, "/usr/local/raml.raml")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.args.err, tt.args.opts...)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestStackTrace_SetType(t *testing.T) {
	type fields struct {
		Severity  *Severity
		Type      *Type
		Location  *Location
		Position  *Position
		Wrapped   *StackTrace
		Err       error
		Message   string
		Info      StructInfo
		List      []*StackTrace
		typeIsSet bool
	}
	type args struct {
		t Type
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(t *testing.T, got *StackTrace)
	}{
		{
			name: "Check set type",
			fields: fields{
				Type: func() *Type { t := Type("type"); return &t }(),
			},
			args: args{
				t: Type("new type"),
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Type == nil {
					t.Errorf("SetType() = %v, want %v", got.Type, nil)
				}
				if *got.Type != "new type" {
					t.Errorf("SetType() = %v, want %v", *got.Type, "new type")
				}
			},
		},
		{
			name: "Check set type with wrapped",
			fields: fields{
				Type: func() *Type { t := Type("type"); return &t }(),
				Wrapped: &StackTrace{
					Type: func() *Type { t := Type("type 2"); return &t }(),
				},
			},
			args: args{
				t: Type("new type"),
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Type == nil {
					t.Errorf("SetType() = %v, want %v", got.Type, nil)
				}
				if *got.Type != "new type" {
					t.Errorf("SetType() = %v, want %v", *got.Type, "new type")
				}
				if got.Wrapped.Type == nil {
					t.Errorf("SetType() = %v, want %v", got.Wrapped.Type, nil)
				}
				if *got.Wrapped.Type != "new type" {
					t.Errorf("SetType() = %v, want %v", *got.Wrapped.Type, "new type")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &StackTrace{
				Severity:  tt.fields.Severity,
				Type:      tt.fields.Type,
				Location:  tt.fields.Location,
				Position:  tt.fields.Position,
				Wrapped:   tt.fields.Wrapped,
				Err:       tt.fields.Err,
				Message:   tt.fields.Message,
				Info:      tt.fields.Info,
				List:      tt.fields.List,
				typeIsSet: tt.fields.typeIsSet,
			}
			got := st.SetType(tt.args.t)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestStackTrace_Append(t *testing.T) {
	type fields struct {
		Severity  *Severity
		Type      *Type
		Location  *Location
		Position  *Position
		Wrapped   *StackTrace
		Err       error
		Message   string
		Info      StructInfo
		List      []*StackTrace
		typeIsSet bool
	}
	type args struct {
		e *StackTrace
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(t *testing.T, got *StackTrace)
	}{
		{
			name:   "Check append",
			fields: fields{},
			args: args{
				e: &StackTrace{Message: "message"},
			},
			want: func(t *testing.T, got *StackTrace) {
				if len(got.List) != 1 {
					t.Errorf("Append() = %v, want %v", len(got.List), 1)
				}
				if got.List[0].Message != "message" {
					t.Errorf("Append() = %v, want %v", got.List[0].Message, "message")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &StackTrace{
				Severity:  tt.fields.Severity,
				Type:      tt.fields.Type,
				Location:  tt.fields.Location,
				Position:  tt.fields.Position,
				Wrapped:   tt.fields.Wrapped,
				Err:       tt.fields.Err,
				Message:   tt.fields.Message,
				Info:      tt.fields.Info,
				List:      tt.fields.List,
				typeIsSet: tt.fields.typeIsSet,
			}
			got := st.Append(tt.args.e)
			if tt.want != nil {
				tt.want(t, got)
			}
		})
	}
}

func TestPosition_String(t *testing.T) {
	type fields struct {
		pos *Position
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Check position: line only",
			fields: fields{
				pos: &Position{Line: 10},
			},
			want: "10",
		},
		{
			name: "Check position: line and column",
			fields: fields{
				pos: &Position{Line: 10, Column: 1},
			},
			want: "10:1",
		},
		{
			name: "Check position: column only: ignored, line is 1",
			fields: fields{
				pos: &Position{Column: 3},
			},
			want: "1",
		},
		{
			name:   "Check position: nil: default line is 1",
			fields: fields{},
			want:   "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.fields.pos
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_optErrNodePosition_Apply(t *testing.T) {
	type fields struct {
		pos *Position
	}
	type args struct {
		e *StackTrace
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   func(t *testing.T, got *StackTrace)
	}{
		{
			name: "Check apply",
			fields: fields{
				pos: &Position{Line: 10, Column: 1},
			},
			args: args{
				e: &StackTrace{},
			},
			want: func(t *testing.T, got *StackTrace) {
				if got.Position == nil {
					t.Errorf("Apply() = %v, want %v", got.Position, nil)
				}
				if got.Position.Line != 10 {
					t.Errorf("Apply() = %v, want %v", got.Position.Line, 10)
				}
				if got.Position.Column != 1 {
					t.Errorf("Apply() = %v, want %v", got.Position.Column, 1)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := optErrNodePosition{
				pos: tt.fields.pos,
			}
			o.Apply(tt.args.e)
			if tt.want != nil {
				tt.want(t, tt.args.e)
			}
		})
	}
}
