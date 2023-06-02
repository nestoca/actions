package jen

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Project represents the configuration file in a jen project's root dir
type Project struct {
	Version string                 `yaml:"version"`
	Vars    map[string]interface{} `yaml:"vars"`
}

// Load loads the project file from given project directory
func Load(dir string) (*Project, error) {
	specFilePath := filepath.Join(dir, "jen.yaml")
	buf, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, fmt.Errorf("loading jen project file: %w", err)
	}
	var project Project
	err = yaml.Unmarshal(buf, &project)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling jen project file yaml: %w", err)
	}
	return &project, nil
}
