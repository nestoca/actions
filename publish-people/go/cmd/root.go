package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "publisher",
		Short: "CLI tool for publishing people pages in Confluence Cloud from Jac catalog",
	}

	rootCmd.AddCommand(NewRenderCmd())
	rootCmd.AddCommand(NewPublishCmd())
	return rootCmd
}
