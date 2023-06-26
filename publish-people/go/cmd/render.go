package cmd

import (
	"github.com/nestoca/actions/publish-people/go/pkg"
	"github.com/spf13/cobra"
)

func NewRenderCmd() *cobra.Command {
	var catalogDir string
	var templateFile string
	var outputFile string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render people pages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pkg.Render(catalogDir, templateFile, outputFile)
		},
	}

	cmd.PersistentFlags().StringVar(&catalogDir, "catalog", "", "Directory of Jac catalog")
	cmd.PersistentFlags().StringVar(&templateFile, "template", "", "Template file to use for rendering")
	cmd.PersistentFlags().StringVar(&outputFile, "output", "", "Output file to write rendered content to")

	return cmd
}
