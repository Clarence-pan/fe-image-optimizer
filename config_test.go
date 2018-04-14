package main

import (
	"testing"
)

func Test_tOptimizerConfig_isOptimizedFile(t *testing.T) {
	cfg := &tOptimizerConfig{
		OptimizedSuffix: ".o",
	}

	type args struct {
		file string
	}

	tests := []struct {
		name string
		cfg  *tOptimizerConfig
		args args
		want bool
	}{
		{
			cfg:  cfg,
			args: args{file: "test.o.jpg"},
			want: true,
		},
		{
			cfg:  cfg,
			args: args{file: "test.o"},
			want: true,
		},
		{
			cfg:  cfg,
			args: args{file: "test.jpg"},
			want: false,
		},
		{
			cfg:  cfg,
			args: args{file: "test"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.isOptimizedFile(tt.args.file); got != tt.want {
				t.Errorf("tOptimizerConfig.isOptimizedFile(%#v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
