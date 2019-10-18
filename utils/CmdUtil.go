package utils

import (
	"bytes"
	"context"
	"os/exec"
	"time"
)

/*
	指令工具
*/

/* 超时执行指令 */
func RunWithTimeout(cmd *exec.Cmd, timeout time.Duration) (string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := cmd.Start(); err != nil {
		return out.String(), err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return out.String(), ctx.Err()
		case err := <-errCh:
			return out.String(), err
		}
	}
}
