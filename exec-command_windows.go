package main

import (
	"context"
	"io"
	"os/exec"
	"syscall"
)

func execCommand(exe string, args []string, stdin io.Reader, stdout io.WriteCloser, ctx context.Context) {

	init := func(cmd *exec.Cmd) {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}

	execCommandCommon(exe, args, stdin, stdout, ctx, init)

}
