package meta

import (
	"fmt"
	"os"

	"github.com/nestoca/metadata/src/internal/logging"
	"github.com/nestoca/metadata/src/internal/meta/jen"
)

func GetReleases() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("getting working dir: %w", err)
	}

	project, err := jen.Load(dir)
	if err == nil {
		logging.Log("Found jen.yaml")
		value, ok := project.Vars["RELEASES"]
		if ok {
			releases, ok := value.(string)
			if ok {
				logging.Log("Using RELEASES var in jen.yaml")
				return releases, nil
			}
		}
	}

	projectName, err := GetProjectName()
	if err != nil {
		return "", err
	}
	logging.Log("Using single release from PROJECT var")
	return projectName, nil
}
