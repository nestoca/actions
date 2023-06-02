package next

import (
	"fmt"
	"time"

	"github.com/nestoca/metadata/src/internal/meta"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Determines next version to use for project in current work dir",
		Long:  "Determines next version to use for project in current work dir",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return cmd
}

func run() error {
	version, err := meta.GetNextVersion(time.Now())
	if err != nil {
		return err
	}
	fmt.Println(version.String())
	return nil
}
