package meta

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/nestoca/metadata/src/internal/logging"
	"github.com/nestoca/metadata/src/internal/shell"
)

var zeroVersion = semver.New("0.0.0")

func GetCurrentVersionAndTagCommit() (*semver.Version, string, error) {
	var noTagCommit = ""

	// Get prefix
	tagPrefix, ok := os.LookupEnv("GIT_TAG_PREFIX")
	if !ok {
		return nil, "", fmt.Errorf("missing required GIT_TAG_PREFIX env var")
	}
	logging.Log("Using GIT_TAG_PREFIX: %s", tagPrefix)

	// Get latest tag on HEAD
	stdout, stderr, err := shell.Exec("git describe --abbrev=0 --tags --match %s\\*", tagPrefix)
	if err != nil {
		// No tags
		if strings.Contains(stderr, "No names found") ||
			strings.Contains(stderr, "aucun nom trouvé") {
			logging.Log("Current version: %s", zeroVersion)
			return zeroVersion, noTagCommit, nil
		}
		return nil, noTagCommit, fmt.Errorf("getting latest tag on HEAD: %w", err)
	}
	latestTagOnHead := strings.TrimSuffix(stdout, "\n")
	logging.Log("Latest tag on HEAD: %s", latestTagOnHead)

	// Get tag's commit
	stdout, _, err = shell.Exec("git rev-list -n 1 %s", latestTagOnHead)
	if err != nil {
		return nil, noTagCommit, fmt.Errorf("getting tag's commit: %w", err)
	}
	tagCommit := strings.TrimSuffix(stdout, "\n")
	logging.Log("Latest tag's commit: %s", tagCommit)

	// Get all tags containing commit SHA (possibly more recent)
	stdout, _, err = shell.Exec("git tag --contains %s", tagCommit)
	if err != nil {
		return nil, noTagCommit, fmt.Errorf("getting tags containing commit: %w", err)
	}
	tags := strings.Split(stdout, "\n")

	// Find current version
	var currentVersion *semver.Version
	for _, tag := range tags {
		// Consider only tags with proper prefix
		tag = strings.TrimSpace(tag)
		if tag == "" || !strings.HasPrefix(tag, tagPrefix) {
			continue
		}

		// Converge towards current version
		version := semver.New(strings.TrimPrefix(tag, tagPrefix))
		if currentVersion == nil {
			currentVersion = version
		} else if version.Compare(*currentVersion) > 0 {
			currentVersion = version
		}
	}

	logging.Log("Current version: %s", currentVersion)
	return currentVersion, tagCommit, nil
}

// GetNextVersion returns the next version for currently checked out branch
// in current work dir. It finds latest tag with given prefix and increments
// it using the following heuristic: it looks for last merge commit branch containing a Github
// PR For `master` branch, function looks for last
// See: https://semver.org/
func GetNextVersion(now time.Time) (*semver.Version, error) {
	currentBranch, err := getCurrentBranch()
	if err != nil {
		return nil, err
	}

	currentVersion, tagCommit, err := GetCurrentVersionAndTagCommit()
	if err != nil {
		return nil, fmt.Errorf("getting current version: %w", err)
	}

	isMaster := currentBranch == "master" || currentBranch == "main"

	mergedBranch := ""
	if isMaster && !currentVersion.Equal(*zeroVersion) {
		mergedBranch, err = getLastMergedBranchSinceCommit(tagCommit)
		if err != nil {
			return nil, err
		}
	}

	nextVersion := bumpVersion(currentVersion, mergedBranch)

	if !isMaster {
		sha, err := getShortCommitSHA()
		if err != nil {
			log.Fatalln("failed shortCommitSHA:", err)
		}

		sanitizedCurrentBranch := getSanitizedVersionMetaIdentifier(currentBranch)
		dateTime := fmt.Sprintf("%04d-%02d-%02d.%02d-%02d-%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
		nextVersion, err = semver.NewVersion(fmt.Sprintf("%s+%s.%s.%s", nextVersion, sanitizedCurrentBranch, sha, dateTime))
		if err != nil {
			log.Fatalln("failed version:", err)
		}
	}

	logging.Log("Next version: %s", nextVersion)
	return nextVersion, nil
}

var invalidVersionMetaIdentifierCharRegex = regexp.MustCompile(`[^a-zA-Z0-9-]`)

func getSanitizedVersionMetaIdentifier(meta string) string {
	return strings.ToLower(invalidVersionMetaIdentifierCharRegex.ReplaceAllString(meta, "-"))
}

var mergeLogRegex = regexp.MustCompile(`[a-zA-Z0-9]+ (\(tag: .*?\) )?Merge pull request #\d+ from nestoca\/(\S+)`)

func getBranchFromMergeLogLine(line string) string {
	matches := mergeLogRegex.FindStringSubmatch(line)
	if matches == nil {
		return ""
	}
	return matches[2]
}

func getLastMergedBranchSinceCommit(commit string) (string, error) {
	// Log merge commits (optionally since given commit hash)
	command := "git log --merges --pretty=oneline"
	if commit != "" {
		command = fmt.Sprintf("%s %s..HEAD", command, commit)
	}
	stdout, _, err := shell.Exec(command)
	if err != nil {
		return "", fmt.Errorf("listing merge commits: %w", err)
	}
	lines := strings.Split(stdout, "\n")

	// No merges?
	if len(lines) == 0 {
		// It's probably commits made directly to `master` branch
		return "", nil
	}

	// Parse branch from merge log
	branch := getBranchFromMergeLogLine(lines[0])
	if branch != "" {
		logging.Log("Last merged branch: %s", branch)
	} else {
		logging.Log("No merged branch found")
	}
	return branch, nil
}

func getShortCommitSHA() (string, error) {
	stdout, _, err := shell.Exec("git rev-parse --short HEAD")
	if err != nil {
		return "", fmt.Errorf("getting short commit sha: %w", err)
	}
	return strings.TrimSpace(stdout), nil
}

func bumpVersion(version *semver.Version, mergedBranch string) *semver.Version {
	newVersion := semver.New(version.String())
	if strings.HasPrefix(mergedBranch, "fix/") {
		logging.Log("Bumping patch version")
		newVersion.BumpPatch()
	} else { // "feat/" or any other branch name
		logging.Log("Bumping minor version")
		newVersion.BumpMinor()
	}
	return newVersion
}

func getCurrentBranch() (string, error) {
	stdout, _, err := shell.Exec("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return "", fmt.Errorf("getting current branch: %w", err)
	}
	branch := strings.TrimSpace(stdout)
	logging.Log("Current branch: %s", branch)
	return branch, nil
}
