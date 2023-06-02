package meta

import (
	"os"
	"testing"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/nestoca/metadata/src/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestBumpVersion(t *testing.T) {
	cases := []struct {
		name         string
		version      *semver.Version
		mergedBranch string
		expected     *semver.Version
	}{
		{
			name:         "fix merged branch",
			version:      semver.New("1.2.3"),
			mergedBranch: "fix/something",
			expected:     semver.New("1.2.4"),
		},
		{
			name:         "no merged branch",
			version:      semver.New("1.2.3"),
			mergedBranch: "",
			expected:     semver.New("1.3.0"),
		},
		{
			name:         "feature merged branch",
			version:      semver.New("1.2.3"),
			mergedBranch: "feat/sparkling",
			expected:     semver.New("1.3.0"),
		},
		{
			name:         "anything",
			version:      semver.New("1.2.3"),
			mergedBranch: "abcdefghijklmnop",
			expected:     semver.New("1.3.0"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := bumpVersion(c.version, c.mergedBranch)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGetSanitizedVersionMetaIdentifier(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "slash and dashes",
			value:    "fix/campaign-monitor-list",
			expected: "fix-campaign-monitor-list",
		},
		{
			name:     "invalid chars",
			value:    "Élégant.and$sharp feature",
			expected: "-l-gant-and-sharp-feature",
		},
		{
			name:     "uppercase chars",
			value:    "Elegant-And-Sharp-Feature",
			expected: "elegant-and-sharp-feature",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := getSanitizedVersionMetaIdentifier(c.value)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGetBranchFromMergeLogLine(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "with tag",
			value:    "8cfcc55a1e7fb312da03c611f3437e3087e837af (tag: api/v0.97.9) Merge pull request #248 from nestoca/fix/campaign-monitor-list",
			expected: "fix/campaign-monitor-list",
		},
		{
			name:     "without tag",
			value:    "8cfcc55a1e7fb312da03c611f3437e3087e837af Merge pull request #248 from nestoca/fix/campaign-monitor-list",
			expected: "fix/campaign-monitor-list",
		},
		{
			name:     "arbitrary message",
			value:    "hello world",
			expected: "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := getBranchFromMergeLogLine(c.value)
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestGetCurrentVersionAndTagCommit(t *testing.T) {
	cases := []struct {
		name              string
		dir               string
		tagPrefix         string
		expectedVersion   *semver.Version
		expectedTagCommit string
	}{
		{
			name:              "",
			dir:               "aFewCommitsWithoutTag",
			tagPrefix:         "whatever",
			expectedVersion:   semver.New("0.0.0"),
			expectedTagCommit: "",
		},
		{
			name:              "fix merge and specific tag prefix",
			dir:               "aFewTagsAndARecentFixGithubMerge",
			tagPrefix:         "api/v",
			expectedVersion:   semver.New("0.103.0"),
			expectedTagCommit: "e026f23033bc3feff1e3037c75f984054bdec697",
		},
		{
			name:              "feature merge and specific tag prefix",
			dir:               "aFewTagsAndARecentFeatureGithubMerge",
			tagPrefix:         "api/v",
			expectedVersion:   semver.New("0.103.0"),
			expectedTagCommit: "e026f23033bc3feff1e3037c75f984054bdec697",
		},
		{
			name:              "non master branch",
			dir:               "nonMasterBranch",
			tagPrefix:         "api/v",
			expectedVersion:   semver.New("0.103.0"),
			expectedTagCommit: "bd33d8f",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testutil.RunInGitDir(t, c.dir, func(t *testing.T) {
				os.Setenv("GIT_TAG_PREFIX", c.tagPrefix)
				actualVersion, actualTagCommit, err := GetCurrentVersionAndTagCommit()
				os.Unsetenv("GIT_TAG_PREFIX")
				assert.NoError(t, err)
				assert.Equal(t, c.expectedVersion, actualVersion, "version")
				assert.Equal(t, shortenHash(c.expectedTagCommit), shortenHash(actualTagCommit), "tag commit")
			})
		})
	}
}

func shortenHash(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func TestGetNextVersion(t *testing.T) {

	now := time.Date(2021, 12, 31, 23, 58, 59, 0, time.UTC)

	cases := []struct {
		name            string
		dir             string
		tagPrefix       string
		expectedVersion *semver.Version
	}{
		{
			name:            "a few commits without tag",
			dir:             "aFewCommitsWithoutTag",
			tagPrefix:       "whatever",
			expectedVersion: semver.New("0.1.0"),
		},
		{
			name:            "fix merge and specific tag prefix",
			dir:             "aFewTagsAndARecentFixGithubMerge",
			tagPrefix:       "api/v",
			expectedVersion: semver.New("0.103.1"),
		},
		{
			name:            "feature merge and specific tag prefix",
			dir:             "aFewTagsAndARecentFeatureGithubMerge",
			tagPrefix:       "api/v",
			expectedVersion: semver.New("0.104.0"),
		},
		{
			name:            "non master branch",
			dir:             "nonMasterBranch",
			tagPrefix:       "api/v",
			expectedVersion: semver.New("0.104.0+feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testutil.RunInGitDir(t, c.dir, func(t *testing.T) {
				os.Setenv("GIT_TAG_PREFIX", c.tagPrefix)
				actualVersion, err := GetNextVersion(now)
				os.Unsetenv("GIT_TAG_PREFIX")
				assert.NoError(t, err)
				assert.Equal(t, c.expectedVersion, actualVersion, "version")
			})
		})
	}
}

func TestGetShortCommitSHA(t *testing.T) {
	testutil.RunInGitDir(t, "aFewCommitsWithoutTag", func(t *testing.T) {
		actual, err := getShortCommitSHA()
		assert.NoError(t, err)
		assert.Equal(t, "97b26b9", actual)
	})
}

func TestGetCurrentBranch(t *testing.T) {
	testutil.RunInGitDir(t, "someBranchCheckedOut", func(t *testing.T) {
		actual, err := getCurrentBranch()
		assert.NoError(t, err)
		assert.Equal(t, "some-branch", actual)
	})
}
