package meta

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/nestoca/metadata/src/internal/testutil"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {

	now := time.Date(2021, 12, 31, 23, 58, 59, 0, time.UTC)

	cases := []struct {
		name            string
		dir             string
		tagPrefix       string
		expectedExports []string
	}{
		{
			name:      "a few commits without tag",
			dir:       "aFewCommitsWithoutTag",
			tagPrefix: "whatever/v",
			expectedExports: []string{
				"PROJECT=aFewCommitsWithoutTag",
				"VERSION=0.1.0",
				"GIT_TAG=whatever/v0.1.0",
				"DOCKER_TAG=0.1.0",
				"RELEASES=aFewCommitsWithoutTag",
			},
		},
		{
			name:      "fix merge and specific tag prefix",
			dir:       "aFewTagsAndARecentFixGithubMerge",
			tagPrefix: "api/v",
			expectedExports: []string{
				"PROJECT=aFewTagsAndARecentFixGithubMerge",
				"VERSION=0.103.1",
				"GIT_TAG=api/v0.103.1",
				"DOCKER_TAG=0.103.1",
				"RELEASES=aFewTagsAndARecentFixGithubMerge",
			},
		},
		{
			name:      "feature merge and specific tag prefix",
			dir:       "aFewTagsAndARecentFeatureGithubMerge",
			tagPrefix: "api/v",
			expectedExports: []string{
				"PROJECT=aFewTagsAndARecentFeatureGithubMerge",
				"VERSION=0.104.0",
				"GIT_TAG=api/v0.104.0",
				"DOCKER_TAG=0.104.0",
				"RELEASES=aFewTagsAndARecentFeatureGithubMerge",
			},
		},
		{
			name:      "non master branch",
			dir:       "nonMasterBranch",
			tagPrefix: "api/v",
			expectedExports: []string{
				"PROJECT=nonMasterBranch",
				"VERSION=0.104.0+feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"GIT_TAG=api/v0.104.0+feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"DOCKER_TAG=0.104.0-feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"RELEASES=nonMasterBranch",
			},
		},
		{
			name:      "non master branch with jen releases",
			dir:       "nonMasterBranchWithJenReleases",
			tagPrefix: "api/v",
			expectedExports: []string{
				"PROJECT=nonMasterBranchWithJenReleases",
				"VERSION=0.104.0+feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"GIT_TAG=api/v0.104.0+feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"DOCKER_TAG=0.104.0-feat-sparkling-and-new.ad8beab.2021-12-31.23-58-59",
				"RELEASES=release1, release2",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			testutil.RunInGitDir(t, c.dir, func(t *testing.T) {
				var actualExports []string
				mockExportFunc := func(key, value string) error {
					actualExports = append(actualExports, fmt.Sprintf("%s=%s", key, value))
					return nil
				}

				os.Setenv("GIT_TAG_PREFIX", c.tagPrefix)
				err := export(now, mockExportFunc)
				os.Unsetenv("GIT_TAG_PREFIX")
				assert.NoError(t, err)

				if diff := deep.Equal(c.expectedExports, actualExports); diff != nil {
					t.Error(diff)
				}
			})
		})
	}
}
