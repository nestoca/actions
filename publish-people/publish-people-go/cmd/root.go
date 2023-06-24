package cmd

import (
	"github.com/nestoca/people-renderer/internal"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	var catalogDir string
	var templateFile string

	rootCmd := &cobra.Command{
		Use:   "renderer",
		Short: "CLI tool for rendering team page in confluence from Jac people and group definitions in nestoca/people repo",
		RunE: func(cmd *cobra.Command, args []string) error {
			return internal.Render(catalogDir, templateFile)
		},
	}

	rootCmd.PersistentFlags().StringVar(&catalogDir, "catalog-dir", "", "Directory of Jac catalog")
	rootCmd.PersistentFlags().StringVar(&templateFile, "template", "", "Template file to use for rendering")

	return rootCmd
}
