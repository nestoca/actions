package releases

import (
	"fmt"

	"github.com/nestoca/metadata/src/internal/meta"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "releases",
		Short: "Determines list of helm release names to deploy using same container image",
		Long: "Determines list of helm release names to deploy using same container image, " +
			"by reading meta.yaml's 'releases' property and falling back to a single release " +
			"named according to 'meta name' command. " +
			"Defining multiple releases allows the same container image to be deployed with " +
			"different values files.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	releases, err := meta.GetReleases()
	if err != nil {
		return err
	}
	fmt.Println(releases)
	return nil
}
