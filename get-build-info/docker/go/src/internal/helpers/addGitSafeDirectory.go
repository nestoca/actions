package helpers

import (
	"fmt"
	"os"

	"github.com/nestoca/metadata/src/internal/logging"
	"github.com/nestoca/metadata/src/internal/shell"
)

// ConfigureGithubWorkspaceAsSafeDirectory marks directory specified by GITHUB_WORKSPACE
// (if defined) in git config as a "safe directory".
//
// Since 2.20.1(?) git has introduced a security measure to validate
// that current user acting upon a git working copy is same user
// as the one who cloned it originally. That introduces the following
// issue with GitHub Actions:
//
// fatal: detected dubious ownership in repository at '/github/workspace'
// To add an exception for this directory, call:
//
//	git config --global --add safe.directory /github/workspace
//
// See:
// - github.com/actions/runner-images/issues/6775
// - https://phabricator.wikimedia.org/T325128
func ConfigureGithubWorkspaceAsSafeDirectory() error {
	githubWorkspace, ok := os.LookupEnv("GITHUB_WORKSPACE")
	if !ok || githubWorkspace == "" {
		return nil
	}

	logging.Log("Configuring GITHUB_WORKSPACE as safe directory: %s", githubWorkspace)
	_, _, err := shell.Exec(fmt.Sprintf("git config --global --add safe.directory %s", githubWorkspace))
	if err != nil {
		return fmt.Errorf("configuring GITHUB_WORKSPACE as safe directory: %w", err)
	}
	return nil
}
