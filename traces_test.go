package stacktrace

import (
	"reflect"
	"testing"
)

func TestStackTrace_GetTraces(t *testing.T) {
	type fields struct {
		Severity        Severity
		Type            Type
		Location        string
		Position        *Position
		Wrapped         *StackTrace
		Err             error
		Message         string
		WrappingMessage string
		Info            StructInfo
		List            []*StackTrace
		typeIsSet       bool
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
				Severity: SeverityError,
				Type:     TypeValidating,
				Location: "/tmp/location.raml",
				Message:  "error message",
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  "/tmp/location.raml:1",
							Severity: SeverityError,
							Message:  "error message",
							Type:     TypeValidating,
						},
					},
				},
			},
		},
		{
			name: "Test with wrapped",
			fields: fields{
				Severity:        SeverityError,
				Type:            TypeValidating,
				Location:        "/tmp/location.raml",
				Position:        &Position{1, 2},
				Message:         "error message",
				WrappingMessage: "wrapping message",
				Wrapped: &StackTrace{
					Severity: SeverityCritical,
					Type:     TypeParsing,
					Location: "/tmp/location2.raml",
					Position: &Position{3, 4},
					Message:  "error message 2",
				},
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  "/tmp/location.raml:1:2",
							Severity: SeverityError,
							Message:  "wrapping message: error message",
							Type:     TypeValidating,
						},
						{
							LinePos:  "/tmp/location2.raml:3:4",
							Severity: SeverityCritical,
							Message:  "error message 2",
							Type:     TypeParsing,
						},
					},
				},
			},
		},
		{
			name: "Test with wrapped and EnsureDuplicates",
			fields: fields{
				Severity:        SeverityError,
				Type:            TypeValidating,
				Location:        "/tmp/location.raml",
				Position:        &Position{1, 2},
				Message:         "error message",
				WrappingMessage: "wrapping message",
				Wrapped: &StackTrace{
					Severity: SeverityCritical,
					Type:     TypeParsing,
					Location: "/tmp/location2.raml",
					Position: &Position{3, 4},
					Message:  "error message 2",
				},
				List: []*StackTrace{
					{
						Severity:        SeverityCritical,
						Type:            TypeParsing,
						Location:        "/tmp/location3.raml",
						Position:        &Position{5, 6},
						Message:         "error message 3",
						WrappingMessage: "wrapping message 3",
						Wrapped: &StackTrace{
							Severity: SeverityCritical,
							Type:     TypeParsing,
							Location: "/tmp/location2.raml", // duplicate location
							Position: &Position{3, 4},       // duplicate position
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
							LinePos:  "/tmp/location.raml:1:2",
							Severity: SeverityError,
							Message:  "wrapping message: error message",
							Type:     TypeValidating,
						},
						{
							LinePos:  "/tmp/location2.raml:3:4",
							Severity: SeverityCritical,
							Message:  "error message 2",
							Type:     TypeParsing,
						},
					},
				},
			},
		},
		{
			name: "Test with list",
			fields: fields{
				Severity: SeverityError,
				Type:     TypeValidating,
				Location: "/tmp/location.raml",
				Position: &Position{1, 2},
				Message:  "error message",
				List: []*StackTrace{
					{
						Severity: SeverityCritical,
						Type:     TypeParsing,
						Location: "/tmp/location2.raml",
						Position: &Position{3, 4},
						Message:  "error message 2",
					},
					{
						Severity: SeverityCritical,
						Type:     TypeParsing,
						Location: "/tmp/location3.raml",
						Position: &Position{5, 6},
						Message:  "error message 3",
					},
				},
			},
			want: []Trace{
				{
					Stack: []Stack{
						{
							LinePos:  "/tmp/location.raml:1:2",
							Severity: SeverityError,
							Message:  "error message",
							Type:     TypeValidating,
						},
					},
				},
				{
					Stack: []Stack{
						{
							LinePos:  "/tmp/location2.raml:3:4",
							Severity: SeverityCritical,
							Message:  "error message 2",
							Type:     TypeParsing,
						},
					},
				},
				{
					Stack: []Stack{
						{
							LinePos:  "/tmp/location3.raml:5:6",
							Severity: SeverityCritical,
							Message:  "error message 3",
							Type:     TypeParsing,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &StackTrace{
				Severity:        tt.fields.Severity,
				Type:            tt.fields.Type,
				Location:        tt.fields.Location,
				Position:        tt.fields.Position,
				Wrapped:         tt.fields.Wrapped,
				Err:             tt.fields.Err,
				Message:         tt.fields.Message,
				WrappingMessage: tt.fields.WrappingMessage,
				Info:            tt.fields.Info,
				List:            tt.fields.List,
				typeIsSet:       tt.fields.typeIsSet,
			}
			if got := st.GetTraces(tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTraces() = %v, want %v", got, tt.want)
			}
		})
	}
}
