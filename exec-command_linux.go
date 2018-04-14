package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

func execCommand(exe string, args []string, stdin io.Reader, stdout io.WriteCloser) {
	log.Printf("[EXEC]: %s %v", exe, args)

	cmd := exec.Command(exe, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	panicIf(err)
}
