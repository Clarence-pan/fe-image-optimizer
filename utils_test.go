package main

import (
	"testing"
)

func Test_basenameWithoutExt(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				file: "a/b/c.jpg",
			},
			want: "c",
		},
		{
			name: "test2",
			args: args{
				file: "a/b/c.zip",
			},
			want: "c",
		},
		{
			name: "test3",
			args: args{
				file: "a/b/c.tar.gz",
			},
			want: "c.tar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := basenameWithoutExt(tt.args.file); got != tt.want {
				t.Errorf("basenameWithoutExt() = %v, want %v", got, tt.want)
			}
		})
	}
}
