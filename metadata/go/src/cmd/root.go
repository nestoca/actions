package cmd

import (
	"github.com/spf13/cobra"
)

// NewRoot creates the root cobra command
func NewRoot() *cobra.Command {
	c := &cobra.Command{
		Use:          "cx",
		Short:        "nesto's CI/CD/Continuous Everyting CLI helper",
		Long:         "nesto's CI/CD/Continuous Everyting CLI helper",
		SilenceUsage: true,
	}

	return c
}
