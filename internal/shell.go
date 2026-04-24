package internal

import (
	"bytes"
	"fmt"
	"os/exec"
)

type ShellResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func Run(name string, args ...string) (ShellResult, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := ShellResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
		return result, fmt.Errorf("命令执行失败 (%s): %s", name, stderr.String())
	}
	if err != nil {
		return result, fmt.Errorf("无法执行命令 %s: %w", name, err)
	}
	return result, nil
}
