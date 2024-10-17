package slogex

import (
	"fmt"
	"log/slog"
	"reflect"
	"testing"

	"github.com/acronis/go-stacktrace"
)

func TestErrToSlogAttr(t *testing.T) {
	type args struct {
		err  error
		opts []stacktrace.TracesOpt
	}
	tests := []struct {
		name string
		args args
		want slog.Attr
	}{
		{
			name: "Test simple",
			args: args{
				err:  stacktrace.New("error message", "location.raml"),
				opts: []stacktrace.TracesOpt{},
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
				opts: []stacktrace.TracesOpt{},
			},
			want: slog.String("error", "error message"),
		},
		{
			name: "Test nil err",
			args: args{
				err:  nil,
				opts: []stacktrace.TracesOpt{},
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
