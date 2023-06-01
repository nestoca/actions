package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var colorRegex = regexp.MustCompile("\x1b\\[[0-9]+m")

func AssertCmd(t *testing.T, cmd *cobra.Command, dir string, expectedOut, expectedErr string) {
	t.Helper()
	err := os.Chdir(filepath.Join("testdata", dir))
	assert.NoError(t, err)
	defer os.Chdir("../..")

	// Redirect stdout
	oldStdout := os.Stdout
	readStdout, writeStdout, _ := os.Pipe()
	os.Stdout = writeStdout

	// Redirect stderr
	oldStderr := os.Stderr
	readStderr, writeStderr, _ := os.Pipe()
	os.Stderr = writeStderr

	cmd.Execute()

	// Read and restore stdout
	writeStdout.Close()
	actualOut, err := ioutil.ReadAll(readStdout)
	assert.NoError(t, err)
	os.Stdout = oldStdout

	assert.Equal(t, expectedOut, string(actualOut), "stdout")

	// Read and restore stderr
	writeStderr.Close()
	actualErr, err := ioutil.ReadAll(readStderr)
	assert.NoError(t, err)
	os.Stderr = oldStderr

	actualErrStr := colorRegex.ReplaceAllString(string(actualErr), "")
	assert.Equal(t, expectedErr, actualErrStr, "stderr")
}
