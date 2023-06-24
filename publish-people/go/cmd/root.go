package cmd

import (
	"github.com/nestoca/actions/publish-people/go/internal"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var catalogDir string
	var templateFile string

	rootCmd := &cobra.Command{
		Use:   "publish-people",
		Short: "CLI tool for publishing people pages in Confluence Cloud from Jac catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Render(catalogDir, templateFile)
		},
	}

	rootCmd.PersistentFlags().StringVar(&catalogDir, "catalog", "", "Directory of Jac catalog")
	rootCmd.PersistentFlags().StringVar(&templateFile, "template", "", "Template file to use for rendering")

	return rootCmd
}
