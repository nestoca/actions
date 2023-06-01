package current

import (
	"fmt"

	"github.com/nestoca/cx/src/internal/meta"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Determines current version of project in current work dir",
		Long:  "Determines current version of project in current work dir",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return cmd
}

func run() error {
	version, _, err := meta.GetCurrentVersionAndTagCommit()
	if err != nil {
		return err
	}
	fmt.Println(version.String())
	return nil
}
