package stacktrace

import (
	"reflect"
	"testing"
)

func TestStackTrace_GetTraces(t *testing.T) {
	var SeverityError Severity = "error"
	var SeverityCritical Severity = "critical"
	var TypeValidating Type = "validating"
	var TypeParsing Type = "parsing"

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
		opts []TracesOpt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Trace
	}{
		{
			name: "Test simple",
			fields: fields{
				Severity: &SeverityError,
				Type:     &TypeValidating,
				Location: func() *Location { s := Location("/tmp/location.raml"); return &s }(),
				Message:  "error message",
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location.raml:1"; return &s }(),
							Severity: &SeverityError,
							Message:  "error message",
							Type:     &TypeValidating,
						},
					},
				},
			},
		},
		{
			name: "Test with wrapped",
			fields: fields{
				Severity: &SeverityError,
				Type:     &TypeValidating,
				Location: func() *Location { s := Location("/tmp/location.raml"); return &s }(),
				Position: &Position{1, 2},
				Message:  "error message",
				Wrapped: &StackTrace{
					Severity: &SeverityCritical,
					Type:     &TypeParsing,
					Location: func() *Location { s := Location("/tmp/location2.raml"); return &s }(),
					Position: &Position{3, 4},
					Message:  "error message 2",
				},
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location.raml:1:2"; return &s }(),
							Severity: &SeverityError,
							Message:  "error message",
							Type:     &TypeValidating,
						},
						{
							LinePos:  func() *string { s := "/tmp/location2.raml:3:4"; return &s }(),
							Severity: &SeverityCritical,
							Message:  "error message 2",
							Type:     &TypeParsing,
						},
					},
				},
			},
		},
		{
			name: "Test with wrapped and EnsureDuplicates",
			fields: fields{
				Severity: &SeverityError,
				Type:     &TypeValidating,
				Location: func() *Location { s := Location("/tmp/location.raml"); return &s }(),
				Position: &Position{1, 2},
				Message:  "error message",
				Wrapped: &StackTrace{
					Severity: &SeverityCritical,
					Type:     &TypeParsing,
					Location: func() *Location { s := Location("/tmp/location2.raml"); return &s }(),
					Position: &Position{3, 4},
					Message:  "error message 2",
				},
				List: []*StackTrace{
					{
						Severity: &SeverityCritical,
						Type:     &TypeParsing,
						Location: func() *Location { s := Location("/tmp/location2.raml"); return &s }(),
						Position: &Position{5, 6},
						Message:  "error message 3",
						Wrapped: &StackTrace{
							Severity: &SeverityCritical,
							Type:     &TypeParsing,
							Location: func() *Location { s := Location("/tmp/location2.raml"); return &s }(), // duplicate location
							Position: &Position{3, 4},                                                        // duplicate position
							Message:  "error message 4",
						},
					},
				},
			},
			args: args{
				opts: []TracesOpt{
					WithEnsureDuplicates(),
				},
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location.raml:1:2"; return &s }(),
							Severity: &SeverityError,
							Message:  "error message",
							Type:     &TypeValidating,
						},
						{
							LinePos:  func() *string { s := "/tmp/location2.raml:3:4"; return &s }(),
							Severity: &SeverityCritical,
							Message:  "error message 2",
							Type:     &TypeParsing,
						},
					},
				},
			},
		},
		{
			name: "Test with list",
			fields: fields{
				Severity: &SeverityError,
				Type:     &TypeValidating,
				Location: func() *Location { s := Location("/tmp/location.raml"); return &s }(),
				Position: &Position{1, 2},
				Message:  "error message",
				List: []*StackTrace{
					{
						Severity: &SeverityCritical,
						Type:     &TypeParsing,
						Location: func() *Location { s := Location("/tmp/location2.raml"); return &s }(),
						Position: &Position{3, 4},
						Message:  "error message 2",
					},
					{
						Severity: &SeverityCritical,
						Type:     &TypeParsing,
						Location: func() *Location { s := Location("/tmp/location3.raml"); return &s }(),
						Position: &Position{5, 6},
						Message:  "error message 3",
					},
				},
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location.raml:1:2"; return &s }(),
							Severity: &SeverityError,
							Message:  "error message",
							Type:     &TypeValidating,
						},
					},
				},
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location2.raml:3:4"; return &s }(),
							Severity: &SeverityCritical,
							Message:  "error message 2",
							Type:     &TypeParsing,
						},
					},
				},
				{
					Stack: []Stack{
						{
							LinePos:  func() *string { s := "/tmp/location3.raml:5:6"; return &s }(),
							Severity: &SeverityCritical,
							Message:  "error message 3",
							Type:     &TypeParsing,
						},
					},
				},
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
			if got := st.GetTraces(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTraces() = %v, want %v", got, tt.want)
			}
		})
	}
}
