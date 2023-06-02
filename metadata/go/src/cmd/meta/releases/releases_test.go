package releases

import (
	"testing"

	"github.com/nestoca/metadata/src/internal/testutil"
)

func TestReleasesFromJenReleases(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	testutil.AssertCmd(t, cmd, "releasesFromJenReleases", "acme1 acme2 acme3\n", "Found jen.yaml\nUsing RELEASES var in jen.yaml\n")
}

func TestReleasesFromJenProject(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	testutil.AssertCmd(t, cmd, "releasesFromJenProject", "acme\n", "Found jen.yaml\nFound jen.yaml\nUsing PROJECT var in jen.yaml\nUsing single release from PROJECT var\n")
}

func TestReleasesFromDirName(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})
	testutil.AssertCmd(t, cmd, "releasesFromDirName", "releasesFromDirName\n", "Using name of current directory\nUsing single release from PROJECT var\n")
}
