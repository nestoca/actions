package cmd

import (
	"github.com/nestoca/actions/publish-people/go/pkg/publish"
	"github.com/spf13/cobra"
)

func NewPublishCmd() *cobra.Command {
	var opts publish.Opts

	cmd := &cobra.Command{
		Use:   "publish",
		Short: "Publish people pages",
		RunE: func(cmd *cobra.Command, args []string) error {
			return publish.Publish(opts)
		},
	}

	cmd.PersistentFlags().StringVar(&opts.BaseUrl, "base-url", "", "Base URL of Confluence Cloud instance (eg: https://your-confluence-instance.atlassian.net)")
	cmd.PersistentFlags().StringVar(&opts.PageID, "page-id", "", "ID of page to update (eg: 12345678)")
	cmd.PersistentFlags().StringVar(&opts.SpaceKey, "space-key", "", "Key of space to update (eg: SPACEKEY)")
	cmd.PersistentFlags().StringVar(&opts.ApiToken, "api-token", "", "API token for authentication")
	cmd.PersistentFlags().StringVar(&opts.PageTitle, "page-title", "", "Updated page title")
	cmd.PersistentFlags().StringVar(&opts.PageFile, "page-file", "", "File containing updated page content")
	cmd.PersistentFlags().StringVar(&opts.Username, "username", "", "Username the update will be made in the name of")

	return cmd
}
