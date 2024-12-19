// Package tool /**
package tool

import (
	"bytes"
	"log/slog"
	"os/exec"
)

// ExecCommand 执行linux 命令
func ExecCommand(command string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			slog.Error("错误信息---recover---->", "err", err)
		}
	}()
	cmd := exec.Command("cmd", "/c", command)
	var output bytes.Buffer
	//cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return output.Bytes(), err
}
