package cmd

import (
	"github.com/nestoca/actions/publish-people/go/pkg"
	"github.com/spf13/cobra"
)

func NewRenderCmd() *cobra.Command {
	var catalogDir string
	var templateFile string

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render people pages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return pkg.Render(catalogDir, templateFile)
		},
	}

	cmd.PersistentFlags().StringVar(&catalogDir, "catalog", "", "Directory of Jac catalog")
	cmd.PersistentFlags().StringVar(&templateFile, "template", "", "Template file to use for rendering")

	return cmd
}
