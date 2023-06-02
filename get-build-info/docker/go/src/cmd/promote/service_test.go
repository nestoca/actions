package promote

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvs(t *testing.T) {
	cmd := New()
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	assert.Error(t, err)
	assert.EqualError(t, err, "missing required env var RELEASES")
}
