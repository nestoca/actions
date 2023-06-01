package project

import (
	"os"
	"testing"

	"github.com/nestoca/cx/src/internal/testutil"
)

func TestProjectFromMetaFile(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	testutil.AssertCmd(t, cmd, "nameFromJenFile", "acme\n", "Found jen.yaml\nUsing PROJECT var in jen.yaml\n")
}

func TestProjectFromDirName(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	testutil.AssertCmd(t, cmd, "nameFromDirName", "nameFromDirName\n", "Using name of current directory\n")
}

func TestProjectFromEnvVar(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	os.Setenv("PROJECT", "someNameFromProjectEnvVar")
	testutil.AssertCmd(t, cmd, "nameFromDirName", "someNameFromProjectEnvVar\n", "Using name from PROJECT env var\n")
	os.Unsetenv("PROJECT")
}
