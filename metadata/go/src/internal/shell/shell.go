package shell

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func Exec(format string, args ...interface{}) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command := fmt.Sprintf(format, args...)
	cmd := exec.Command("sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", stderr.String(), fmt.Errorf("executing shell command %q: %w: stderr: %s", command, err, stderr.String())
	}
	return stdout.String(), stderr.String(), nil
}
