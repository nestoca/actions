package main

import (
	"context"
	"github.com/nestoca/get-build-info/src/internal/meta"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		<-ctx.Done()
		stop()
	}()

	rootCmd := &cobra.Command{
		Use:   "get-build-info",
		Short: "Determines and exports all build info to GitHub outputs",
		Long:  "Determines and exports all build info to GitHub outputs",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return meta.Export()
		},
	}

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(-1)
	}
}
