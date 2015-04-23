package util

import (
	"os"
	"os/exec"
	"runtime"
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func ExecScript(script string) (*exec.Cmd, error) {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	if other := os.Getenv("SHELL"); other != "" {
		shell = other
	}
	cmd := exec.Command(shell, flag, script)
	return cmd, nil
}
