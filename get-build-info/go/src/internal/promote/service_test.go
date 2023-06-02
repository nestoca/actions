package promote

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func NewTestIOStreams() (IOStreams, *bytes.Buffer, *bytes.Buffer, *bytes.Buffer) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}

	return IOStreams{
		In:     in,
		Out:    out,
		ErrOut: errOut,
	}, in, out, errOut
}

func TestPromoteNewService(t *testing.T) {
	stream, _, buf, _ := NewTestIOStreams()

	t.Cleanup(func() {
		os.RemoveAll("testdata/releases/empty/values")
		os.Chdir("../../../")
	})

	t.Setenv("DOCKER_TAG", "0.2.3")
	t.Setenv("REGISTRY", "gcr.io/nesto-ci-1a2b3c")
	t.Setenv("GENERIC_CHART_VERSION", "0.5.4")

	_ = os.Chdir("testdata/releases/empty")

	cmd := PromoteOptions{
		Releases:            []string{"foo"},
		ReleaseTemplatePath: "../../templates/releases.yaml.gotmpl",
		ValueTemplatePath:   "../../templates/values.yaml.gotmpl",
		IOStreams:           stream,
		DryRun:              true,
	}

	err := cmd.Promote()
	assert.NoError(t, err)

	expected := `foo:
  chart: generic
  installed: false
  tag: 0.2.3
  version: 0.5.4
`

	if diff := deep.Equal(expected, buf.String()); diff != nil {
		t.Error(diff)
	}
}

func TestPromoteExistingService(t *testing.T) {
	stream, _, buf, _ := NewTestIOStreams()

	t.Cleanup(func() {
		os.RemoveAll("testdata/releases/existing/values")
		os.Chdir("../../../")
	})

	t.Setenv("DOCKER_TAG", "0.2.4")
	t.Setenv("REGISTRY", "gcr.io/nesto-ci-1a2b3c")
	t.Setenv("GENERIC_CHART_VERSION", "0.5.4")

	_ = os.Chdir("testdata/releases/existing")

	cmd := PromoteOptions{
		Releases:            []string{"foo"},
		ReleaseTemplatePath: "../../templates/releases.yaml.gotmpl",
		ValueTemplatePath:   "../../templates/values.yaml.gotmpl",
		IOStreams:           stream,
		DryRun:              true,
	}

	err := cmd.Promote()
	assert.NoError(t, err)

	expected := `foo:
  chart: generic
  installed: true
  tag: 0.2.4
  version: 0.5.4
`
	if diff := deep.Equal(expected, buf.String()); diff != nil {
		t.Error(diff)
	}
}

func TestPromoteNewServiceSorted(t *testing.T) {
	stream, _, buf, _ := NewTestIOStreams()

	t.Cleanup(func() {
		os.RemoveAll("testdata/releases/existing/values")
		os.Chdir("../../../")
	})

	t.Setenv("DOCKER_TAG", "0.0.4")
	t.Setenv("REGISTRY", "gcr.io/nesto-ci-1a2b3c")
	t.Setenv("GENERIC_CHART_VERSION", "0.5.4")

	_ = os.Chdir("testdata/releases/existing")

	cmd := PromoteOptions{
		Releases:            []string{"bar"},
		ReleaseTemplatePath: "../../templates/releases.yaml.gotmpl",
		ValueTemplatePath:   "../../templates/values.yaml.gotmpl",
		IOStreams:           stream,
		DryRun:              true,
	}

	err := cmd.Promote()
	assert.NoError(t, err)

	expected := `bar:
  chart: generic
  installed: false
  tag: 0.0.4
  version: 0.5.4
foo:
  chart: generic
  installed: true
  tag: 0.2.3
  version: 0.5.4
`

	if diff := deep.Equal(expected, buf.String()); diff != nil {
		t.Error(diff)
	}
}
