package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func RunInDir(t *testing.T, dir string, f func(t *testing.T)) {
	t.Helper()
	err := os.Chdir(filepath.Join("testdata", dir))
	assert.NoError(t, err)
	defer os.Chdir("../..")

	f(t)
}
