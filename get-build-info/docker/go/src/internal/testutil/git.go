package testutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunInGitDir(t *testing.T, dir string, f func(t *testing.T)) {
	t.Helper()

	RunInDir(t, dir, func(t *testing.T) {
		err := os.Rename("_git", ".git")
		assert.NoError(t, err, "renaming _git => .git")
		defer os.Rename(".git", "_git")

		f(t)
	})
}
