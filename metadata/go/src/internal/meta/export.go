package meta

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nestoca/metadata/src/internal/logging"
)

// Exports all metadata env vars to shell and codefresh by executing the `export-var`
// bash script for each one (ie: PROJECT, VERSION, GIT_TAG, DOCKER_TAG).
// Note: Stdout must be sourced by caller in order for those vars to actually be
// exported to current shell (ie: by using the `$(...)` syntax).
func Export() error {
	return export(time.Now(), ExecExportVarFunc)
}

func export(now time.Time, exportFunc func(key, value string) error) error {
	// PROJECT
	project, err := GetProjectName()
	if err != nil {
		return fmt.Errorf("determining PROJECT: %w", err)
	}
	if err := exportFunc("PROJECT", project); err != nil {
		return err
	}

	// VERSION
	semver, err := GetNextVersion(now)
	if err != nil {
		return fmt.Errorf("determining VERSION: %w", err)
	}
	version := semver.String()
	if err := exportFunc("version", version); err != nil {
		return err
	}

	// GIT_TAG
	gitTagPrefix, ok := os.LookupEnv("GIT_TAG_PREFIX")
	if !ok {
		return fmt.Errorf("missing required GIT_TAG_PREFIX env var")
	}
	gitTag := gitTagPrefix + version
	if err := exportFunc("git-tag", gitTag); err != nil {
		return err
	}

	// DOCKER_TAG
	dockerTag := strings.ReplaceAll(version, "+", "-")
	if err := exportFunc("docker-tag", dockerTag); err != nil {
		return err
	}

	// RELEASES
	releases, err := GetReleases()
	if err != nil {
		return fmt.Errorf("determining RELEASES: %w", err)
	}
	if err := exportFunc("releases", releases); err != nil {
		return err
	}

	return nil
}

func ExecExportVarFunc(key, value string) error {
	logging.Log("exporting %s=%s", key, value)
	command := fmt.Sprintf("::set-output name=%s::%s", key, value)
	cmd := exec.Command("sh", "-c", command)
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("executing shell command %q: %w", command, err)
	}
	return nil
}
