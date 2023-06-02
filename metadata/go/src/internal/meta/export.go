package meta

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// Exports all metadata variables to GitHub Actions outputs.
func Export() error {
	return export(time.Now(), SetGitHubOutputFunc)
}

func export(now time.Time, exportFunc func(key, value string) error) error {
	// Output project
	project, err := GetProjectName()
	if err != nil {
		return fmt.Errorf("determining project: %w", err)
	}
	if err := exportFunc("project", project); err != nil {
		return err
	}

	// Output version
	semver, err := GetNextVersion(now)
	if err != nil {
		return fmt.Errorf("determining version: %w", err)
	}
	version := semver.String()
	if err := exportFunc("version", version); err != nil {
		return err
	}

	// Output git-tag
	gitTagPrefix, ok := os.LookupEnv("GIT_TAG_PREFIX")
	if !ok {
		return fmt.Errorf("missing required GIT_TAG_PREFIX env var")
	}
	gitTag := gitTagPrefix + version
	if err := exportFunc("git-tag", gitTag); err != nil {
		return err
	}

	// Output docker-tag
	dockerTag := strings.ReplaceAll(version, "+", "-")
	if err := exportFunc("docker-tag", dockerTag); err != nil {
		return err
	}

	// Output releases
	releases, err := GetReleases()
	if err != nil {
		return fmt.Errorf("determining releases: %w", err)
	}
	if err := exportFunc("releases", releases); err != nil {
		return err
	}

	return nil
}

func SetGitHubOutputFunc(key, value string) error {
	outputFilePath := os.Getenv("GITHUB_OUTPUT")

	if outputFilePath == "" {
		return fmt.Errorf("GITHUB_OUTPUT environment variable not set")
	}

	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s=%s", key, value)
	if err != nil {
		return err
	}

	return nil
}
