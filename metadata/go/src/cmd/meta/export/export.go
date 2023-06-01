package export

import (
	"github.com/nestoca/cx/src/internal/meta"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "export",
		Short: "Exports all metadata env vars to shell and codefresh",
		Long: "Exports all metadata env vars to shell and codefresh. " +
			"Note: the output of this command must be sourced by using this syntax: $(cx meta export)",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
}

func run() error {
	return meta.Export()
}
