package meta

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nestoca/metadata/src/internal/logging"
	"github.com/nestoca/metadata/src/internal/meta/jen"
)

func GetProjectName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working dir: %w", err)
	}

	name, ok := os.LookupEnv("PROJECT")
	if ok && name != "" {
		logging.Log("Using name from PROJECT env var")
		return name, nil
	}

	project, err := jen.Load(dir)
	if err == nil {
		logging.Log("Found jen.yaml")
		value, ok := project.Vars["PROJECT"]
		if ok {
			name, ok := value.(string)
			if ok {
				logging.Log("Using PROJECT var in jen.yaml")
				return name, nil
			}
		}
	}

	logging.Log("Using name of current directory")
	return filepath.Base(dir), nil
}
