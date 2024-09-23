package stacktrace

import (
	"fmt"
	"log/slog"
	"reflect"
	"testing"
)

func TestErrToSlogAttr(t *testing.T) {
	type args struct {
		err  error
		opts []TracesOpt
	}
	tests := []struct {
		name string
		args args
		want slog.Attr
	}{
		{
			name: "Test simple",
			args: args{
				err:  New("error message", "location.raml"),
				opts: []TracesOpt{},
			},
			want: slog.Group(
				"tracebacks", "traces", []slog.Attr{
					slog.Group(
						"0", "stack", []slog.Attr{
							slog.Group(
								"0",
								slog.String("type", "parsing"),
								slog.String("severity", "error"),
								slog.String("position", "location.raml:1"),
								slog.String("message", "error message"),
							),
						},
					),
				},
			),
		},
		{
			name: "Test is not a stacktrace",
			args: args{
				err:  fmt.Errorf("error message"),
				opts: []TracesOpt{},
			},
			want: slog.String("error", "error message"),
		},
		{
			name: "Test nil err",
			args: args{
				err:  nil,
				opts: []TracesOpt{},
			},
			want: slog.Attr{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrToSlogAttr(tt.args.err, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				slog.Error("got", got)
				slog.Error("want", tt.want)
				t.Errorf("ErrToSlogAttr() = %v, want %v", got, tt.want)
			}
		})
	}
}
