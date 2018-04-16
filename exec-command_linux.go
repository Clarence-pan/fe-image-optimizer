package main

import (
	"context"
	"io"
	"os/exec"
)

func execCommand(exe string, args []string, stdin io.Reader, stdout io.WriteCloser, ctx context.Context) {
	init := func(cmd *exec.Cmd) {}

	execCommandCommon(exe, args, stdin, stdout, ctx, init)
}
