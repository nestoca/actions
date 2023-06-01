package version

import (
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Commands for determining latest and next versions",
	}
}
