package promote

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/nestoca/promote-codefresh/src/internal/helpers"
	"github.com/nestoca/promote-codefresh/src/internal/logging"
	"gopkg.in/yaml.v3"
)

type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

type PromoteOptions struct {
	Releases            []string
	ReleaseTemplatePath string
	ValueTemplatePath   string
	DryRun              bool

	IOStreams
}

const releasesFile = "releases.yaml"

var (
	releaseTemplate *template.Template
	valuesTemplate  *template.Template
)

func (o *PromoteOptions) Promote() error {
	logging.Log("Load existing releases.yaml")
	bytes, err := ioutil.ReadFile(releasesFile)
	if err != nil {
		return err
	}
	var doc map[string]map[string]interface{}
	if err = yaml.Unmarshal(bytes, &doc); err != nil {
		return err
	}

	// Working with an empty releases file?
	if doc == nil {
		logging.Log("Starting new releases.yaml")
		doc = map[string]map[string]interface{}{}
	}

	// Load templates
	releaseTemplate = template.Must(template.ParseFiles(o.ReleaseTemplatePath))
	valuesTemplate = template.Must(template.ParseFiles(o.ValueTemplatePath))

	// Patch releases.yaml
	for _, release := range o.Releases {
		entry, exists := doc[release]
		isNew := !exists
		if isNew {
			logging.Log("Adding new release: %s", release)
			entry = make(map[string]interface{})
			doc[release] = entry
		} else {
			logging.Log("Patching existing release: %s", release)
		}
		patch, err := o.patch(release, isNew)
		if err != nil {
			return err
		}
		apply(release, entry, patch)

		o.createValuesFile(release)
	}

	// Write releases.yaml
	var w io.Writer
	if o.DryRun {
		w = o.IOStreams.Out
	} else {
		fi, err := os.Create(releasesFile)
		if err != nil {
			return err
		}
		w = fi
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(doc)
}

// createValues creates the default values file for helm chart for a new service
func (o *PromoteOptions) createValuesFile(release string) error {
	writePath := path.Join("values", fmt.Sprintf("%s.yaml.gotmpl", release))

	if helpers.PathExists(writePath) {
		logging.Log("Leaving existing values file untouched: %s", writePath)
		return nil
	}

	var buf bytes.Buffer
	err := valuesTemplate.Execute(&buf, map[string]interface{}{"Env": env()})

	if err != nil {
		return fmt.Errorf("rendering values template: %w", err)
	}

	logging.Log("Creating new values file: %s", writePath)
	return ioutil.WriteFile(writePath, buf.Bytes(), fs.ModePerm)
}

// patch generates the yaml patch for the current release
func (o *PromoteOptions) patch(release string, isNew bool) (map[string]map[string]interface{}, error) {
	var buf bytes.Buffer
	if err := releaseTemplate.Execute(&buf, map[string]interface{}{
		"Release": release,
		"IsNew":   isNew,
		"Env":     env(),
	}); err != nil {
		return make(map[string]map[string]interface{}), fmt.Errorf("executing release template: %w", err)
	}

	var out map[string]map[string]interface{}
	if err := yaml.Unmarshal(buf.Bytes(), &out); err != nil {
		return make(map[string]map[string]interface{}), fmt.Errorf("unmarshalling template result: %w", err)
	}

	return out, nil
}

// apply applies the patch to the existing yaml representation
func apply(release string, original map[string]interface{}, patch map[string]map[string]interface{}) {
	for k, v := range patch[release] {
		original[k] = v
	}
}

// env returns a map of current process' environment variables suitable for injection into a go template
func env() map[string]string {
	m := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.Split(e, "=")
		m[parts[0]] = parts[1]
	}
	return m
}
