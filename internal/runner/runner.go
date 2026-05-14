package runner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"time"
)

// Check runs commandToRun via shell and matches combined stdout+stderr against pattern.
// timeout of 0 defaults to 30 seconds.
func Check(commandToRun, expectedPattern string, timeout time.Duration) (output string, ok bool, err error) {
	re, err := regexp.Compile(expectedPattern)
	if err != nil {
		return "", false, fmt.Errorf("invalid expected_output regex: %w", err)
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", commandToRun)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", commandToRun)
	}
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	runErr := cmd.Run()
	combined := out.String() + errBuf.String()

	if ctx.Err() == context.DeadlineExceeded {
		return combined, false, fmt.Errorf("command timed out after %s", timeout)
	}
	if runErr != nil {
		if re.MatchString(combined) {
			return combined, true, nil
		}
		return combined, false, nil
	}
	return combined, re.MatchString(combined), nil
}
