package project

import (
	"fmt"

	"github.com/nestoca/cx/src/internal/meta"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "project",
		Short: "Determines project's name to use for current build",
		Long: "Determines project's name to use for current build, " +
			"by looking for PROJECT env var in current process then in jen.yaml, and falling back to current dir name.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	name, err := meta.GetProjectName()
	if err != nil {
		return err
	}
	fmt.Println(name)
	return nil
}
