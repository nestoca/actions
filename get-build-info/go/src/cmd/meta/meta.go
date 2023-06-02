package meta

import (
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "meta",
		Short: "Commands for working with CI/CD meta data, such as names, releases and versions",
	}
}
