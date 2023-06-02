package promote

import (
	"fmt"
	"os"
	"strings"

	"github.com/nestoca/metadata/src/internal/promote"
	"github.com/spf13/cobra"
)

// New creates a cobra command

func New() *cobra.Command {
	o := promote.PromoteOptions{
		IOStreams: promote.IOStreams{
			Out:    os.Stdout,
			In:     os.Stdin,
			ErrOut: os.Stderr,
		},
	}

	var cmd = &cobra.Command{
		Use:   "promote",
		Short: "Promotes a service",
		Long:  "Promoting a service will upgrade the given releases to the specified tag",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validate(); err != nil {
				return err
			}
			o.Releases = strings.Split(os.Getenv("RELEASES"), " ")
			o.ReleaseTemplatePath = os.Getenv("RELEASE_TEMPLATE")
			o.ValueTemplatePath = os.Getenv("VALUES_TEMPLATE")
			return o.Promote()
		},
	}

	cmd.Flags().BoolVar(&o.DryRun, "dry-run", false, "writes output to standard out rather then overwritting release file")

	return cmd
}

// validate validates the runtime args
func validate() error {
	for _, envKey := range []string{"RELEASES", "DOCKER_TAG", "REGISTRY", "GENERIC_CHART_VERSION", "RELEASE_TEMPLATE", "VALUES_TEMPLATE"} {
		if _, present := os.LookupEnv(envKey); !present {
			return fmt.Errorf("missing required env var %s", envKey)
		}
	}

	return nil
}
