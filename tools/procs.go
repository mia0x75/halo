package tools

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

// TimeoutedExec executes a timeouted command.
// The program path is defined by the name arguments, args are passed as arguments to the program.
//
// TimeoutedExec returns process output as a string (stdout) , and stderr as an error.
func TimeoutedExec(timeout time.Duration, name string, arg ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, arg...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	if s := string(stderr.Bytes()); len(s) > 0 {
		return "", errors.New(s)
	}

	return string(stdout.Bytes()), nil
}
