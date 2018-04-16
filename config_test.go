package main

import (
	"os"
	"strings"
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

func Test_tOptimizerConfig_isInBlackList(t *testing.T) {
	cfg := &tOptimizerConfig{
		OptimizedSuffix: ".o",
		BlackList: []string{
			"foo",
			"bar",
		},
	}
	type args struct {
		file string
	}

	normalizePath := func(p string) string {
		return strings.Replace(p, "/", string([]rune{os.PathSeparator}), -1)
	}

	tests := []struct {
		name string
		cfg  *tOptimizerConfig
		args args
		want bool
	}{
		{
			cfg:  cfg,
			args: args{file: "test.jpg"},
			want: false,
		},
		{
			cfg:  cfg,
			args: args{file: "a/b/foo"},
			want: true,
		},
		{
			cfg:  cfg,
			args: args{file: normalizePath("foo/test.jpg")},
			want: true,
		},
		{
			cfg:  cfg,
			args: args{file: normalizePath("bar/test.jpg")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.isInBlackList(tt.args.file); got != tt.want {
				t.Errorf("tOptimizerConfig.isInBlackList(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
